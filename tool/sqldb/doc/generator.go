package doc

import (
	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
	"fmt"
	"github.com/ktpswjz/database/sqldb"
	"io"
	"strings"
)

type Generator struct {
	Database sqldb.SqlDatabase
}

func (s *Generator) CreateWord(file io.Writer, title string) error {
	if s.Database == nil {
		return fmt.Errorf("sql database is nil")
	}
	tables, err := s.Database.Tables()
	if err != nil {
		return err
	}
	rootCatalogs := make(catalogs, 0)
	otherCatalogName := "其它"
	otherTables := make([]*sqldb.SqlTable, 0)
	for _, table := range tables {
		// 0: 标题1
		// 1：    标题2
		// 2：        标题3
		// 3+：注释
		comments := s.getLines(table.Description)
		lineCount := len(comments)
		if lineCount < 1 {
			otherTables = append(otherTables, table)
			continue
		}
		rootCatalogName := comments[0]
		if rootCatalogName == "" {
			otherTables = append(otherTables, table)
			continue
		}
		rootCatalog := rootCatalogs.getCatalog(rootCatalogName)
		if rootCatalog == nil {
			rootCatalog = &catalog{
				name:     rootCatalogName,
				level:    1,
				children: make(catalogs, 0),
				tables:   make([]*sqldb.SqlTable, 0),
			}
			rootCatalogs = append(rootCatalogs, rootCatalog)
		}
		if lineCount < 2 {
			rootCatalog.tables = append(rootCatalog.tables, table)
			continue
		}
		childCatalogName := comments[1]
		if childCatalogName == "" {
			childCatalogName = otherCatalogName
		}
		childCatalog := rootCatalog.children.getCatalog(childCatalogName)
		if childCatalog == nil {
			childCatalog = &catalog{
				name:     childCatalogName,
				level:    2,
				children: make(catalogs, 0),
				tables:   make([]*sqldb.SqlTable, 0),
			}
			rootCatalog.children = append(rootCatalog.children, childCatalog)
		}
		childCatalog.tables = append(childCatalog.tables, table)
	}

	if len(otherTables) > 0 {
		rootCatalog := rootCatalogs.getCatalog(otherCatalogName)
		if rootCatalog == nil {
			rootCatalog = &catalog{
				name:     otherCatalogName,
				level:    1,
				children: make(catalogs, 0),
				tables:   make([]*sqldb.SqlTable, 0),
			}
			rootCatalogs = append(rootCatalogs, rootCatalog)
		}
		rootCatalog.tables = append(rootCatalog.tables, otherTables...)
	}

	doc := document.New()
	pfs := doc.Styles.AddStyle("ParagraphFontStyle", wml.ST_StyleTypeParagraph, true)
	pfs.RunProperties().SetFontFamily("宋体")

	if len(title) > 0 {
		paragraph := doc.AddParagraph()
		paragraph.Properties().SetAlignment(wml.ST_JcCenter)

		run := paragraph.AddRun()
		run.Properties().SetBold(true)
		run.Properties().SetSize(18)
		run.AddText(title)

		doc.AddParagraph().AddRun().AddPageBreak()
	}

	err = s.createCatalogs(doc, rootCatalogs, "")
	if err != nil {
		return err
	}

	ftr := doc.AddFooter()
	para := ftr.AddParagraph()
	para.Properties().AddTabStop(3*measurement.Inch, wml.ST_TabJcCenter, wml.ST_TabTlcNone)
	run := para.AddRun()
	run.AddTab()
	run.AddFieldWithFormatting(document.FieldCurrentPage, "", false)
	run.AddText("/")
	run.AddFieldWithFormatting(document.FieldNumberOfPages, "", false)
	doc.BodySection().SetFooter(ftr, wml.ST_HdrFtrDefault)

	return doc.Save(file)
}

func (s *Generator) createCatalogs(doc *document.Document, values catalogs, prefix string) error {
	count := len(values)
	if count < 1 {
		return nil
	}

	for index := 0; index < count; index++ {
		value := values[index]

		paragraph := doc.AddParagraph()
		paragraph.Properties().SetHeadingLevel(value.level)
		paragraph.AddRun().AddText(fmt.Sprint(prefix, index+1, ". ", value.name))

		if len(value.children) > 0 {
			s.createCatalogs(doc, value.children, fmt.Sprint(index+1, "."))
		}

		tableCount := len(value.tables)
		for tableIndex := 0; tableIndex < tableCount; tableIndex++ {
			table := value.tables[tableIndex]
			columns, err := s.Database.Columns(table.Name)
			if err != nil {
				return err
			}

			// 0: 标题1
			// 1：    标题2
			// 2：        标题3
			// 3+：注释
			comments := s.getLines(table.Description)
			lineCount := len(comments)

			paragraph := doc.AddParagraph()
			paragraph.Properties().SetHeadingLevel(value.level + 1)
			run := paragraph.AddRun()
			run.AddText(fmt.Sprint(prefix, index+1, ".", tableIndex+1, ". "))
			if lineCount > 2 {
				run.AddText(comments[2])
			}
			run.AddText(table.Name)

			s.createWordTable(doc.AddTable(), columns)

			if lineCount > 3 {
				for lineIndex := 3; lineIndex < lineCount; lineIndex++ {
					doc.AddParagraph().AddRun().AddText(comments[lineIndex])
				}
			}
		}
	}

	return nil
}

func (s *Generator) createWordTable(table document.Table, columns []*sqldb.SqlColumn) error {
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderThick, color.Gray, measurement.Point)

	// title
	row := table.AddRow()
	s.addWordTableRowTitle(row, "序号", "字段中文名", "字段名", "类型", "可空", "主键", "说明", "默认值")

	// contents
	if len(columns) < 1 {
		return nil
	}
	for index, column := range columns {
		row = table.AddRow()
		comments := s.getLines(column.Comment)

		//序号
		c1 := fmt.Sprint(index + 1)

		// 字段中文名
		c2 := ""
		if len(comments) > 0 {
			c2 = comments[0]
		}

		// 字段名
		c3 := column.Name

		// 类型
		c4 := column.Type

		// 可空
		c5 := ""
		if !column.Nullable {
			c5 = "N"
		}

		// 主键
		c6 := ""
		if column.PrimaryKey {
			c6 = "Y"
		}

		// 说明
		c7 := ""
		if len(comments) > 3 {
			c7 = comments[3]
		}

		// 默认值
		c8 := column.DataDisplay

		s.addWordTableRowContent(row, c1, c2, c3, c4, c5, c6, c7, c8)
	}

	return nil
}

func (s *Generator) addWordTableRowTitle(row document.Row, titles ...string) {
	if len(titles) < 1 {
		return
	}

	for index, title := range titles {
		cell := row.AddCell()
		cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.LightGray)
		paragraph := cell.AddParagraph()
		run := paragraph.AddRun()
		run.AddText(title)
		run.Properties().SetSize(8)

		if index == 0 || index == 4 || index == 5 {
			cell.Properties().SetWidth(25)
		} else if index == 7 {
			cell.Properties().SetWidth(30)
		}
	}
}

func (s *Generator) addWordTableRowContent(row document.Row, titles ...string) {
	if len(titles) < 1 {
		return
	}

	for index, title := range titles {
		cell := row.AddCell()
		paragraph := cell.AddParagraph()
		run := paragraph.AddRun()
		run.AddText(title)
		run.Properties().SetSize(8)

		if index == 0 {
			paragraph.Properties().SetAlignment(wml.ST_JcCenter)
		} else if index == 4 || index == 5 {
			paragraph.Properties().SetAlignment(wml.ST_JcCenter)
		}
	}
}

func (s *Generator) getLines(value string) []string {
	lines := strings.Replace(value, "\r", "", -1)
	return strings.Split(lines, "\n")
}
