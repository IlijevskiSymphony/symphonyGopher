package registrations

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const registrationsColl = "registrations"

func NewRepository(session *mgo.Session) *repository {
	return &repository{m: session}
}

type repository struct {
	m *mgo.Session
}

func (r *repository) Create(registration *registration) error {
	c := r.m.DB("").C(registrationsColl)

	if err := c.Insert(registration); err != nil {
		return errors.Wrapf(err, "Cannot create registration.")
	}
	return nil
}

func (r *repository) FindByID(id string) (*registration, error) {
	c := r.m.DB("").C(registrationsColl)

	var rs []registration
	err := c.Find(bson.M{"id": id}).All(&rs)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read registrations.")
	}
	if rs == nil || len(rs) == 0 {
		return nil, nil
	}
	return &rs[0], nil
}

func (r *repository) FindByPartner(partner string) (*registration, error) {
	c := r.m.DB("").C(registrationsColl)

	var rs []registration
	err := c.Find(bson.M{"partnerid": partner}).All(&rs)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read registrations.")
	}
	if rs == nil || len(rs) == 0 {
		return nil, nil
	}
	return &rs[0], nil
}

func (r *repository) FindByReference(reference string) (*registration, error) {
	c := r.m.DB("").C(registrationsColl)

	var rs []registration
	err := c.Find(bson.M{"reference": reference}).All(&rs)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read registrations.")
	}
	if rs == nil || len(rs) == 0 {
		return nil, nil
	}
	return &rs[0], nil
}

func (r *repository) Update(registration *registration) error {
	c := r.m.DB("").C(registrationsColl)

	if _, err := c.Upsert(bson.M{"id": registration.ID}, registration); err != nil {
		return errors.Wrapf(err, "Cannot update registration.")
	}
	return nil
}

func (r *repository) Delete(registration *registration) error {
	c := r.m.DB("").C(registrationsColl)

	if err := c.Remove(bson.M{"id": registration.ID}); err != nil {
		return errors.Wrapf(err, "Cannot delete registration.")
	}
	return nil
}

func (r *repository) Close() {
	r.m.Close()
}
