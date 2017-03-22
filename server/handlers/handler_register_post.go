package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/Sirupsen/logrus"
	"github.com/basiqio/developer-dashboard/server/configuration"
	"github.com/basiqio/developer-dashboard/server/mail"
	"github.com/basiqio/developer-dashboard/server/partners"
	"github.com/basiqio/developer-dashboard/server/registrations"
	"github.com/pborman/uuid"
)

type RegistrationPost struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HandlerRegisterPost struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
}

func (h HandlerRegisterPost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var posted RegistrationPost
	if err := decoder.Decode(&posted); err != nil {
		logrus.Infof("Cannot decode registration. Error: %s", err)
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

	registrationsRepo := registrations.NewRepository(h.SessionFn())
	defer registrationsRepo.Close()
	partnersRepo := partners.NewRepository(h.SessionFn())
	defer partnersRepo.Close()

	p, err := partnersRepo.FindByEmail(posted.Email)

	if err != nil {
		logrus.Infof("Cannot find partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p != nil && p.Status == partners.StatusActive {
		logrus.Info("Partner with given email already exists.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if p != nil {
		p = partners.New(p.ID, posted.Email, posted.Password, h.Configuration.HashSalt)
	} else {
		p = partners.New(uuid.New(), posted.Email, posted.Password, h.Configuration.HashSalt)
	}

	err = partnersRepo.Update(p)

	if err != nil {
		logrus.Infof("Cannot update partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r := registrations.New(uuid.New(), uuid.New(), p.ID)

	if err := registrationsRepo.Create(r); err != nil {
		logrus.Infof("Cannot create registration. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ms := mail.Service{Host: h.Configuration.MandrillHost}
	m := mail.Mail{
		Key: h.Configuration.MandrillApiKey,
		Message: mail.Message{
			Html:      fmt.Sprintf(mail.Template.Message, r.Reference, r.Reference),
			Text:      "",
			Subject:   mail.Template.Subject,
			FromEmail: mail.Template.Sender.Email,
			FromName:  mail.Template.Sender.Name,
			To:        mail.Receivers{mail.Receiver{Email: p.Email, Name: "", Type: "to"}},
		},
		Async: false,
	}

	if err := ms.Send(m); err != nil {
		logrus.Infof("Cannot send registration mail. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Status = registrations.StatusVerificationSent
	if err := registrationsRepo.Update(r); err != nil {
		logrus.Infof("Cannot update registration. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Please check your inbox!")
}
