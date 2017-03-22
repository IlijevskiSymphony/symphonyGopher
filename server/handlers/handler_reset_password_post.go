package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/basiqio/developer-dashboard/server/configuration"
	"github.com/basiqio/developer-dashboard/server/partners"
	"github.com/basiqio/developer-dashboard/server/passwordreset"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"
)

type ResetPasswordPost struct {
	Password  string `json:"password"`
	Reference string `json:"id"`
}

type HandlerResetPasswordPost struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
	Store         *sessions.CookieStore
}

func (h HandlerResetPasswordPost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var posted ResetPasswordPost
	if err := decoder.Decode(&posted); err != nil {
		logrus.Infof("Cannot decode registration. Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !partners.ValidatePasswordFormat(posted.Password) {
		logrus.Infof("Password format invalid.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	partnersRepo := partners.NewRepository(h.SessionFn())
	defer partnersRepo.Close()
	passwordReset := passwordreset.NewRepository(h.SessionFn())
	defer passwordReset.Close()

	rp, err := passwordReset.FindByReference(posted.Reference)

	if err != nil {
		logrus.Infof("Cannot find password reset for partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rp == nil {
		logrus.Info("Password reset with given id doesn't exists")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	p, err := partnersRepo.FindByID(rp.PartnerID)

	if err != nil {
		logrus.Infof("Cannot find partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if p == nil {
		logrus.Info("Partner with given email doesn't exists")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	p.PasswordHash = partners.Hash(posted.Password + h.Configuration.HashSalt)

	err = partnersRepo.Update(p)

	if err != nil {
		logrus.Infof("Cannot update partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
