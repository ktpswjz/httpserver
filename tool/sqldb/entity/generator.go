package entity

import (
	"fmt"
	"github.com/ktpswjz/database/sqldb"
)

type Generator struct {
	Database sqldb.SqlDatabase
}

func (s *Generator) CreateEntity(database, model *Package) error {
	if s.Database == nil {
		return fmt.Errorf("sql database is nil")
	}

	if database == nil && model == nil {
		return fmt.Errorf("parameter invalid")
	}

	tables, err := s.Database.Tables()
	if err != nil {
		return err
	}
	views, err := s.Database.Views()
	if err == nil {
		if len(views) > 0 {
			tables = append(tables, views...)
		}
	}

	modelEntity := &entityModel{}
	modelEntity.pkg = model
	dbEntity := &entityDatabase{}
	dbEntity.pkg = database

	for _, table := range tables {
		columns, err := s.Database.Columns(table.Name)
		if err != nil {
			return err
		}

		if modelEntity.pkg != nil {
			err = modelEntity.create(table, columns)
			if err != nil {
				return err
			}
		}

		if dbEntity.pkg != nil {
			err = dbEntity.create(table, columns, modelEntity)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
