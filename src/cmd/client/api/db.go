package api

import (
	zerolog "github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
)

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

// NewMDatabase use for creating or adding the database into rest api configuration
func (rcf *RestConf) NewMDatabase(dbname string) *mgo.Database {
	database := rcf.MSession.DB(dbname)

	return database
}

func (rcf *RestConf) NewMCollection(cname string) *mgo.Collection {
	collection := rcf.Mdb.C(cname)

	return collection
}
