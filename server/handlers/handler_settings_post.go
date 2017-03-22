package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/Sirupsen/logrus"
	"github.com/IlijevskiSymphony/symphonyGopher/server/configuration"
	"github.com/IlijevskiSymphony/symphonyGopher/server/partners"
	"github.com/gorilla/sessions"
)

type SettingsPost struct {
	Name        string `json:"name"`
	Company     string `json:"company"`
	CompanySize string `json:"companySize"`
	Language    string `json:"language"`
}

type HandlerSettingsPost struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
	Store         *sessions.CookieStore
}

func (h HandlerSettingsPost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var posted SettingsPost
	if err := decoder.Decode(&posted); err != nil {
		logrus.Infof("Cannot decode settings. Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cs, ok := partners.CompanySize(posted.CompanySize)
	if !ok {
		logrus.Infof("Cannot decode company size.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l, ok := partners.Language(posted.Language)
	if !ok {
		logrus.Infof("Cannot decode language.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	partnersRepo := partners.NewRepository(h.SessionFn())
	defer partnersRepo.Close()

	session, err := h.Store.Get(req, "login-session")
	if err != nil {
		w.Write([]byte("err"))
	}

	id, ok := session.Values["id"].(string)
	if !ok {
		logrus.Errorf("Wrong ID type. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	partner, err := partnersRepo.FindByID(id)
	if err != nil {
		logrus.Errorf("Cannot find partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if partner == nil {
		logrus.Errorf("Partner with id %s doesn't exists. Error: %s", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	partner.Name = posted.Name
	partner.Company = posted.Company
	partner.CompanySize = cs
	partner.Language = l
	partner.Status = partners.StatusActive
	err = partnersRepo.Update(partner)
	if err != nil {
		logrus.Errorf("Cannot find partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
