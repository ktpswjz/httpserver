package memory

import (
	"github.com/hashicorp/go-memdb"
	"github.com/ktpswjz/httpserver/example/webserver/model"
	"github.com/ktpswjz/httpserver/types"
	"time"
)

type Token interface {
	Set(entity *model.Token) error
	Del(entity interface{}) error
	Get(id string) (*model.Token, error)
	List(userId uint64) ([]model.Token, error)
}

func NewToken(expirationMinutes int64, log types.Log) (Token, error) {
	instance := &innerToken{}
	tableName := instance.tableName()

	// Create the DB schema
	dbSchema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableName: &memdb.TableSchema{
				Name: tableName,
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"userId": &memdb.IndexSchema{
						Name:    "userId",
						Unique:  false,
						Indexer: &memdb.UintFieldIndex{Field: "UserID"},
					},
				},
			},
		},
	}

	// Create a new database
	db, err := memdb.NewMemDB(dbSchema)
	if err != nil {
		return nil, err
	}
	instance.SetLog(log)
	instance.db = db
	instance.expiration = time.Duration(expirationMinutes) * time.Minute

	if expirationMinutes > 0 {
		go func(interval time.Duration) {
			instance.checkExpiration(interval)
		}(5 * time.Minute)
	}

	return instance, nil
}

type innerToken struct {
	types.Base

	db         *memdb.MemDB
	expiration time.Duration
}

func (s *innerToken) tableName() string {
	return "token"
}

func (s *innerToken) checkExpiration(interval time.Duration) {
	for {
		time.Sleep(interval)
		s.LogDebug("begin checking token expiration...")
		now := time.Now()
		s.deleteExpiration()
		s.LogDebug("end checking token expiration, time elapse: ", time.Now().Sub(now))
	}
}

func (s *innerToken) deleteExpiration() {
	txn := s.db.Txn(false)
	defer txn.Abort()

	tableName := s.tableName()
	raw, err := txn.Get(tableName, "userId")
	if err != nil || raw == nil {
		return
	}

	expTime := time.Now().Add(-s.expiration)
	entities := make([]*model.Token, 0)
	row := raw.Next()
	for row != nil {
		entity := row.(*model.Token)
		row = raw.Next()
		if entity == nil {
			continue
		}
		if entity.ActiveTime.After(expTime) {
			continue
		}

		entities = append(entities, entity)

		s.LogDebug("token expired: id=", entity.ID, ", userId=", entity.UserID, ", userAccount=", entity.UserAccount)
	}

	count := len(entities)
	if count < 1 {
		return
	}

	txn.Abort()
	txn = s.db.Txn(true)
	defer txn.Abort()

	for i := 0; i < count; i++ {
		id := entities[i].ID
		err = txn.Delete(tableName, entities[i])
		if err == nil {
			continue
		}

		s.LogWarning("delete expired token(id=", id, ") fail")
	}
	txn.Commit()
}

func (s *innerToken) Set(entity *model.Token) error {
	txn := s.db.Txn(true)
	err := txn.Insert(s.tableName(), entity)
	if err != nil {
		return err
	}
	txn.Commit()

	return nil
}

func (s *innerToken) Get(id string) (*model.Token, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First(s.tableName(), "id", id)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, err
	}

	return raw.(*model.Token), nil
}

func (s *innerToken) Del(entity interface{}) error {
	txn := s.db.Txn(true)
	err := txn.Delete(s.tableName(), entity)
	if err != nil {
		return err
	}
	txn.Commit()

	return nil
}

func (s *innerToken) List(userId uint64) ([]model.Token, error) {
	entities := make([]model.Token, 0)
	txn := s.db.Txn(false)
	defer txn.Abort()

	raw, err := txn.Get(s.tableName(), "userId", userId)
	if err != nil || raw == nil {
		return entities, err
	}

	row := raw.Next()
	for row != nil {
		entity := row.(*model.Token)
		row = raw.Next()
		if entity == nil {
			continue
		}

		info := model.Token{}
		entity.CopyTo(&info)
		entities = append(entities, info)
	}

	return entities, nil
}
