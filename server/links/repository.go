package links

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const linksColl = "links"

type repository struct {
	m *mgo.Session
}

func NewRepository(session *mgo.Session) *repository {
	return &repository{m: session}
}

func (r *repository) Update(link *Link) error {
	c := r.m.DB("").C(linksColl)

	if _, err := c.Upsert(bson.M{"id": link.ID}, link); err != nil {
		return errors.Wrapf(err, "Cannot update links.")
	}

	return nil
}

func (r *repository) FindByID(id string) (*Link, error) {
	c := r.m.DB("").C(linksColl)

	var links []Link
	err := c.Find(bson.M{"id": id}).All(&links)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read links.")
	}
	if links == nil || len(links) == 0 {
		return nil, nil
	}
	return &links[0], nil
}

func (r *repository) FindByPartnerID(partnerId string) (*[]Link, error) {
	c := r.m.DB("").C(linksColl)

	var links []Link
	err := c.Find(bson.M{"partnerid": partnerId}).All(&links)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read applications.")
	}
	if links == nil || len(links) == 0 {
		return nil, nil
	}
	return &links, nil
}

func (r *repository) Close() {
	r.m.Close()
}
