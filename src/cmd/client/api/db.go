package api

import (
	"github.com/DapperBlondie/blog-system/src/cmd/client/models"
	zerolog "github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoTools struct {
	MSession    *mgo.Session
	Mdb         *mgo.Database
	MCollection *mgo.Collection
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

// NewMDatabase use for creating or adding the database into rest api configuration
func (rcf *MongoTools) NewMDatabase(dbname string) *mgo.Database {
	database := rcf.MSession.DB(dbname)

	return database
}

func (rcf *MongoTools) NewMCollection(cname string) *mgo.Collection {
	collection := rcf.Mdb.C(cname)

	return collection
}

func (rcf *MongoTools) GetAuthorIDByItsUsername(uname string) (string, error) {
	var author *models.AuthorPayload
	err := rcf.MCollection.Find(bson.M{"name": uname}).One(&author)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		return "", err
	}

	return author.ID.Hex(), nil
}
