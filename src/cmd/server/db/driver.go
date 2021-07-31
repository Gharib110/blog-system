package db

import (
	zerolog "github.com/rs/zerolog/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// MDatabase use it when you just have one Database and multiple collections
type MDatabase struct {
	MSession     *mgo.Session
	Mdb          *mgo.Database
	MCollections map[string]*mgo.Collection
}

func NewSession(dsn string) (*mgo.Session, error) {
	session, err := mgo.Dial(dsn)
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return nil, err
	}
	err = session.Ping()
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return nil, err
	}

	return session, err
}

// AddDatabase use for adding database for our app
func (m *MDatabase) AddDatabase(dbname string) {
	m.Mdb = m.MSession.DB(dbname)
}

func (m *MDatabase) AddCollection(cname string) {
	m.MCollections[cname] = m.Mdb.C(cname)
}

// tempAuthorCreator just hard coding for adding two authors
func (m *MDatabase) tempAuthorCreator() {
	err := m.MCollections["authors"].Insert(&Author{
		ID:     bson.NewObjectIdWithTime(time.Now()),
		Name:   "Dapper",
		Career: "Developer",
	})
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}

	err = m.MCollections["authors"].Insert(&Author{
		ID:     bson.NewObjectIdWithTime(time.Now()),
		Name:   "V",
		Career: "Mercenary Outlaw",
	})
	if err != nil {
		zerolog.Fatal().Msg(err.Error())
		return
	}
}
