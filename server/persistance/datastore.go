package persistance

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

type DataStore struct {
	master *mgo.Session
}

func NewDataStore(connection string) (*DataStore, error) {
	session, err := mgo.Dial(connection)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to connect to Mongo DB on '%s'.", connection)
	}
	return &DataStore{session}, nil
}

func (ds *DataStore) Close() {
	ds.master.Close()
}

func (ds *DataStore) Session() *mgo.Session {
	s := ds.master.Copy()
	s.SetMode(mgo.Monotonic, true)
	return s
}
