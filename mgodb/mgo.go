package mgodb

import (
	"errors"
	"github.com/hzxiao/goutil/log"
	"gopkg.in/bson.v2"
	"gopkg.in/mgo.v2"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Uri string
}

var DB *Database

type Database struct {
	uri     string
	Session *mgo.Session
	DB      *mgo.Database
	C       func(name string) *mgo.Collection
}

func Init(cfg *Config) error {
	var err error
	DB, err = openDB(cfg.Uri)
	if err != nil {
		return err
	}
	return err
}

func openDB(uri string) (*Database, error) {
	db := &Database{uri: uri}
	var err error
	db.Session, err = mgo.Dial(uri)
	if err != nil {
		return db, err
	}

	dbName, err := db.getDBName(uri)
	if err != nil {
		return db, err
	}
	db.DB = db.Session.DB(dbName)
	db.C = db.DB.C

	go db.pingLoop()
	return db, nil
}

func (db *Database) pingLoop() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		<-ticker.C
		err := db.ping()
		if err == nil {
			continue
		}
		//handle err
		for {
			DB, err = openDB(db.uri)
			if err != nil {
				log.Error("[mgodb] try to dial mongo by url(%v) fail: %v", db.uri, err)
				time.Sleep(5 * time.Second)
				continue
			}
			log.Info("[mgodb] reconnect to mongo success")
			return
		}
	}
}

func (db *Database) ping() (err error) {
	errClosed := errors.New("Closed explicitly")
	defer func() {
		if pe := recover(); pe != nil {
			if db.Session != nil {
				db.Session.Clone()
				err = errClosed
			}
		}
	}()

	err = db.Session.Ping()
	if err == nil {
		return nil
	}
	if err.Error() == "Closed explicitly" || err.Error() == "EOF" {
		db.Session.Clone()
		return errClosed
	}
	return
}

func (db *Database) EnsureAllIndex(indexMap map[string][]mgo.Index) (err error) {
	for coll, indexs := range indexMap {
		for _, index := range indexs {
			err = db.C(coll).EnsureIndex(index)
			if err != nil {
				return
			}
		}
	}
	return
}

func (db *Database) getDBName(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		return "", errors.New("empty dbName")
	}
	return dbName, nil
}

func (db *Database) Insert(coll string, docs ...interface{}) error {
	return DB.C(coll).Insert(docs...)
}

func (db *Database) One(coll string, find, selector bson.M, v interface{}) error {
	return DB.C(coll).Find(find).Select(selector).One(v)
}

func (db *Database) Update(coll string, finder, updater bson.M) error {
	return DB.C(coll).Update(finder, updater)
}

func (db *Database) UpdateAll(coll string, finder, updater bson.M) (*mgo.ChangeInfo, error) {
	return DB.C(coll).UpdateAll(finder, updater)
}

func (db *Database) All(collectionName string, cond, selector bson.M, sort []string, skip, limit int, needCount bool, v interface{}) (int, error) {
	query := DB.C(collectionName).Find(cond).Sort(sort...).Select(selector)
	var count int
	var err error
	if needCount {
		count, err = query.Count()
		if err != nil {
			return 0, err
		}
	}

	if skip > 0 {
		query = query.Skip(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	return count, query.All(v)
}

func (db *Database) FindAndModify(coll string, finder, updater bson.M, upsert, returnNew, remove bool, result interface{}) (*mgo.ChangeInfo, error) {
	info, err := DB.C(coll).Find(finder).Apply(mgo.Change{
		Update:    updater,
		Upsert:    upsert,
		ReturnNew: returnNew,
		Remove:    remove,
	}, result)
	return info, err
}
