package handlers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/IlijevskiSymphony/symphonyGopher/server/configuration"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"
)

type HandlerLogoutGet struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
	Store         *sessions.CookieStore
}

func (h HandlerLogoutGet) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	session, err := h.Store.Get(req, "login-session")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.Values["id"] == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	logrus.Infof("Cookie invalidation, partener uuid: %s", session.Values["id"])

	session.Options.MaxAge = -1
	session.Save(req, w)

	w.WriteHeader(http.StatusOK)

}
