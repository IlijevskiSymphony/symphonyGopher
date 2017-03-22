package applications

import (
  "github.com/pkg/errors"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const applicationsColl = "applications"

type repository struct {
  m *mgo.Session
}

func NewRepository(session *mgo.Session) *repository {
  return &repository{m: session}
}

func (r* repository) Update(application *Application) error {
  c := r.m.DB("").C(applicationsColl)

  if _, err := c.Upsert(bson.M{"id": application.ID}, application); err != nil {
    return errors.Wrapf(err, "Cannot update application.")
  }

  return nil
}

func (r *repository) FindByID(id string) (*Application, error) {
	c := r.m.DB("").C(applicationsColl)

	var applications []Application
	err := c.Find(bson.M{"id": id}).All(&applications)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read applications.")
	}
	if applications == nil || len(applications) == 0 {
		return nil, nil
	}
	return &applications[0], nil
}

func (r *repository) FindByPartnerID(partnerId string) (*[]Application, error) {
	c := r.m.DB("").C(applicationsColl)

	var applications []Application
	err := c.Find(bson.M{"partnerid": partnerId}).All(&applications)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read applications.")
	}
	if applications == nil || len(applications) == 0 {
		return nil, nil
	}
	return &applications, nil
}


func (r *repository) Close() {
	r.m.Close()
}
