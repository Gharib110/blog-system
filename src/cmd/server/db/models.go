package db

import (
	"gopkg.in/mgo.v2/bson"
)

type BlogItem struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AuthorID string        `json:"author_id" bson:"author_id"`
	Content  string        `json:"content" bson:"content"`
	Title    string        `json:"title" bson:"title"`
}

type Author struct {
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name" bson:"name"`
	Career string        `json:"career" bson:"career"`
}
