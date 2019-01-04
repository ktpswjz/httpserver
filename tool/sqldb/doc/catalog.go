package doc

import "github.com/ktpswjz/database/sqldb"

type catalog struct {
	name  string
	level int

	children catalogs
	tables   []*sqldb.SqlTable
}

type catalogs []*catalog

func (s catalogs) getCatalog(name string) *catalog {
	count := len(s)
	for i := 0; i < count; i++ {
		item := s[i]
		if item.name == name {
			return item
		}
	}
	return nil
}
