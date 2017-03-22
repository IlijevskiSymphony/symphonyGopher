package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/IlijevskiSymphony/symphonyGopher/server/configuration"
	"github.com/IlijevskiSymphony/symphonyGopher/server/links"
	"github.com/IlijevskiSymphony/symphonyGopher/server/partners"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/pborman/uuid"
)

type LinkPost struct {
	Link string   `json:"link"`
	Tags []string `json:"tags"`
}

type HandlerLinkAddPost struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
	Store         *sessions.CookieStore
}

func (h HandlerLinkAddPost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var posted LinkPost
	if err := decoder.Decode(&posted); err != nil {
		logrus.Infof("Cannot decode add application data. Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := h.Store.Get(req, "login-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["id"] == nil {
		logrus.Infof("Session id is nil")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	partnerID := session.Values["id"].(string)

	partnersRepo := partners.NewRepository(h.SessionFn())
	defer partnersRepo.Close()

	partner, err := partnersRepo.FindByID(partnerID)

	if err != nil {
		logrus.Errorf("Cannot find partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if partner.Status != partners.StatusActive {
		logrus.Infof("Partner status is not active.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	link := links.New(uuid.New(), uuid.New(), posted.Link, posted.Tags)

	linksRepo := links.NewRepository(h.SessionFn())
	defer linksRepo.Close()

	err = linksRepo.Update(link)

	if err != nil {
		logrus.Infof("Cannot update link. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
