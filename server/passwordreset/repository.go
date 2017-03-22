package passwordreset

import (
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const passwordResetColl = "password_reset"

func NewRepository(session *mgo.Session) *repository {
	return &repository{m: session}
}

type repository struct {
	m *mgo.Session
}

func (r *repository) Update(passwordreset *PasswordReset) error {
	c := r.m.DB("").C(passwordResetColl)

	if _, err := c.Upsert(bson.M{"id": passwordreset.ID}, passwordreset); err != nil {
		return errors.Wrapf(err, "Cannot update password_reset.")
	}
	return nil
}

func (r *repository) Delete(passwordreset *PasswordReset) error {
	c := r.m.DB("").C(passwordResetColl)

	if err := c.Remove(bson.M{"id": passwordreset.ID}); err != nil {
		return errors.Wrapf(err, "Cannot delete registration.")
	}
	return nil
}

func (r *repository) FindByPartner(partner string) (*PasswordReset, error) {
	c := r.m.DB("").C(passwordResetColl)

	var rs []PasswordReset
	err := c.Find(bson.M{"partnerid": partner}).All(&rs)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read password_reset.")
	}
	if rs == nil || len(rs) == 0 {
		return nil, nil
	}
	return &rs[0], nil
}

func (r *repository) FindByReference(reference string) (*PasswordReset, error) {
	c := r.m.DB("").C(passwordResetColl)

	var rs []PasswordReset
	err := c.Find(bson.M{"reference": reference}).All(&rs)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot read password_reset.")
	}
	if rs == nil || len(rs) == 0 {
		return nil, nil
	}
	return &rs[0], nil
}

func (r *repository) Close() {
	r.m.Close()
}
