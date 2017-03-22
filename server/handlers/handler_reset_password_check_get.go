package handlers

import (
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/basiqio/developer-dashboard/server/configuration"
	"github.com/basiqio/developer-dashboard/server/partners"
	"github.com/basiqio/developer-dashboard/server/passwordreset"
	"github.com/gorilla/sessions"
	"github.com/pborman/uuid"
	mgo "gopkg.in/mgo.v2"
)

type HandlerResetPasswordCheckGet struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
	Store         *sessions.CookieStore
}

func (h HandlerResetPasswordCheckGet) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/resetPasswordCheck/")

	if uuid.Parse(id) == nil {
		logrus.Infof("Registration verification link is not valid.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	partnersRepo := partners.NewRepository(h.SessionFn())
	defer partnersRepo.Close()
	passwordReset := passwordreset.NewRepository(h.SessionFn())
	defer passwordReset.Close()

	rp, err := passwordReset.FindByReference(id)

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

	return
}
