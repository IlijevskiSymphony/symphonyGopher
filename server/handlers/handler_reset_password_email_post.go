package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/Sirupsen/logrus"
	"github.com/IlijevskiSymphony/symphonyGopher/server/configuration"
	"github.com/IlijevskiSymphony/symphonyGopher/server/mail"
	"github.com/IlijevskiSymphony/symphonyGopher/server/partners"
	"github.com/IlijevskiSymphony/symphonyGopher/server/passwordreset"
	"github.com/pborman/uuid"
)

type RecivedPostData struct {
	Email string `json:"email"`
}

type HandlerResetPasswordEmailPost struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
}

func (h HandlerResetPasswordEmailPost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var posted RecivedPostData
	if err := decoder.Decode(&posted); err != nil {
		logrus.Infof("Cannot decode email. Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !partners.ValidateEmailFormat(posted.Email) {
		logrus.Infof("Email format invalid.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	partnersRepo := partners.NewRepository(h.SessionFn())
	defer partnersRepo.Close()
	passwordReset := passwordreset.NewRepository(h.SessionFn())
	defer passwordReset.Close()

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

	rp, err := passwordReset.FindByPartner(p.ID)

	if err != nil {
		logrus.Infof("Cannot find password reset for partner. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rp == nil {
		rp = passwordreset.New(uuid.New(), uuid.New(), p.ID)
	} else {
		rp = passwordreset.New(rp.ID, uuid.New(), p.ID)
	}

	err = passwordReset.Update(rp)

	if err != nil {
		logrus.Infof("Cannot update password reset. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ms := mail.Service{Host: h.Configuration.MandrillHost}
	m := mail.Mail{
		Key: h.Configuration.MandrillApiKey,
		Message: mail.Message{
			Html:      fmt.Sprintf(mail.PasswordResetTemplate.Message, p.Name, rp.Reference, rp.Reference),
			Text:      "",
			Subject:   mail.PasswordResetTemplate.Subject,
			FromEmail: mail.PasswordResetTemplate.Sender.Email,
			FromName:  mail.PasswordResetTemplate.Sender.Name,
			To:        mail.Receivers{mail.Receiver{Email: p.Email, Name: "", Type: "to"}},
		},
		Async: false,
	}

	if err := ms.Send(m); err != nil {
		logrus.Infof("Cannot send registration mail. Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}
