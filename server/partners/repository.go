package partners

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const partnersColl = "partners"

func NewRepository(session *mgo.Session) *repository {
	return &repository{m: session}
}

type repository struct {
	m *mgo.Session
}

func (r *repository) Update(partner *Partner) error {
	c := r.m.DB("").C(partnersColl)

	if _, err := c.Upsert(bson.M{"email": partner.Email}, partner); err != nil {
		return errors.Wrapf(err, "Cannot update partner.")
	}
	return nil
}

func (r *repository) FindByEmail(email string) (*Partner, error) {
	c := r.m.DB("").C(partnersColl)

	var ps []Partner
	err := c.Find(bson.M{"email": email}).All(&ps)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read partners.")
	}
	if ps == nil || len(ps) == 0 {
		return nil, nil
	}
	return &ps[0], nil
}

func (r *repository) FindByID(id string) (*Partner, error) {
	c := r.m.DB("").C(partnersColl)

	var ps []Partner
	err := c.Find(bson.M{"id": id}).All(&ps)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read partners.")
	}
	if ps == nil || len(ps) == 0 {
		return nil, nil
	}
	return &ps[0], nil
}

func (r *repository) Close() {
	r.m.Close()
}
