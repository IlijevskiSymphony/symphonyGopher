package handlers

import (
	"io"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"

	"github.com/Sirupsen/logrus"
	"github.com/basiqio/developer-dashboard/server/configuration"
	"github.com/basiqio/developer-dashboard/server/partners"
	"github.com/basiqio/developer-dashboard/server/registrations"
	"github.com/gorilla/sessions"
	"github.com/pborman/uuid"
)

type HandlerRegisterVerificationGet struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
	Store         *sessions.CookieStore
}

func (h HandlerRegisterVerificationGet) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/signup/accept/")

	if uuid.Parse(id) == nil {
		logrus.Infof("Registration verification link is not valid.")
		http.Redirect(w, req, "/404", http.StatusFound)
		return
	}

	registrationsRepo := registrations.NewRepository(h.SessionFn())
	defer registrationsRepo.Close()
	partnersRepo := partners.NewRepository(h.SessionFn())
	defer partnersRepo.Close()

	r, err := registrationsRepo.FindByReference(id)

	if err != nil {
		logrus.Errorf("Cannot find registration. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r == nil {
		logrus.Info("Registration doesn't exists.")
		http.Redirect(w, req, "/404", http.StatusFound)
		return
	}

	p, err := partnersRepo.FindByID(r.PartnerID)

	if err != nil {
		logrus.Errorf("Cannot find partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p.Status == partners.StatusVerified {
		http.Redirect(w, req, "/dashboard", http.StatusFound)
		io.WriteString(w, id)
		return
	}

	if r.Status != registrations.StatusVerificationSent {
		logrus.Infof("Registration verification link is not valid.")
		http.Redirect(w, req, "/verification-invalid", http.StatusFound)
		return
	}

	p.Status = partners.StatusActive
	if err := partnersRepo.Update(p); err != nil {
		logrus.Infof("Cannot update partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := registrationsRepo.Delete(r); err != nil {
		logrus.Infof("Cannot delete registration. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session, err := h.Store.Get(req, "login-session")
	session.Values["id"] = p.ID
	session.Save(req, w)

	http.Redirect(w, req, "/dashboard/settings", http.StatusFound)
	io.WriteString(w, id)
}
