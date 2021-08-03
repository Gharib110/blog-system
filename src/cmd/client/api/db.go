package api

import (
	zerolog "github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
)

type MongoTools struct {
	MSession    *mgo.Session
	Mdb         *mgo.Database
	MCollection map[string]*mgo.Collection
}

func NewMSession(dsn string) (*mgo.Session, error) {
	session, err := mgo.Dial(dsn)
	if err != nil {
		zerolog.Error().Msg(err.Error() + "; In NewMSession occurred.")
		return nil, err
	}

	err = session.Ping()
	if err != nil {
		zerolog.Error().Msg(err.Error() + "; Problem with Pinging.")
		return nil, err
	}

	return session, err
}

// NewMDatabase use for creating or adding a database into the db
func (rcf *MongoTools) NewMDatabase(dbname string) {
	rcf.Mdb = rcf.MSession.DB(dbname)

}

// NewMCollection use for adding or creating a collection into the db
func (rcf *MongoTools) NewMCollection(cname string) {
	rcf.MCollection[cname] = rcf.Mdb.C(cname)
}
