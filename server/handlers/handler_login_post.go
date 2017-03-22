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

type LoginPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email  string          `json:"email"`
	Status partners.Status `json:"status"`
}

type HandlerLoginPost struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
	Store         *sessions.CookieStore
}

func (h HandlerLoginPost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var posted LoginPost
	if err := decoder.Decode(&posted); err != nil {
		logrus.Infof("Cannot decode settings. Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !partners.ValidateEmailFormat(posted.Email) {
		logrus.Infof("Email format invalid.")
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

	p, err := partnersRepo.FindByEmail(posted.Email)

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

	passwordHash := partners.Hash(posted.Password + h.Configuration.HashSalt)
	if p.PasswordHash == passwordHash {
		if p != nil && p.Status == partners.StatusNotVerified {
			response := LoginResponse{p.Email, partners.StatusNotVerified}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			return
		}

		session, err := h.Store.Get(req, "login-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["id"] = p.ID
		session.Save(req, w)

		if p != nil && p.Status == partners.StatusVerified {

			response := LoginResponse{p.Email, partners.StatusVerified}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)
			return
		}

		if p != nil && p.Status == partners.StatusActive {
			response := LoginResponse{p.Email, partners.StatusActive}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)

			return
		}
	} else {
		logrus.Infof("Invalid password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}
