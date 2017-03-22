package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/IlijevskiSymphony/symphonyGopher/server/configuration"
	"github.com/IlijevskiSymphony/symphonyGopher/server/partners"
	"github.com/IlijevskiSymphony/symphonyGopher/server/registrations"
	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
)

type ResendPost struct {
	Email string `json:"email"`
}

type HandlerRegisterVerificationResendPost struct {
	Configuration configuration.Configuration
	SessionFn     func() *mgo.Session
}

func (h HandlerRegisterVerificationResendPost) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var posted ResendPost
	if err := decoder.Decode(&posted); err != nil {
		logrus.Infof("Cannot decode email. Error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if posted.Email != "" {
		registrationsRepo := registrations.NewRepository(h.SessionFn())
		defer registrationsRepo.Close()

		partnersRepo := partners.NewRepository(h.SessionFn())
		defer partnersRepo.Close()

		partner, err := partnersRepo.FindByEmail(posted.Email)
		if err != nil {
			logrus.Infof("Cannot find partner. Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if partner == nil {
			logrus.Info("Partner doesn't exists")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		registration, err := registrationsRepo.FindByPartner(partner.ID)

		if err != nil {
			logrus.Infof("Cannot find registration. Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if registration == nil {
			logrus.Info("Registration doesn't exists")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		registration.Reference = uuid.New()

		if err := registrationsRepo.Update(registration); err != nil {
			logrus.Infof("Cannot update registration. Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//Email sent here
		//
		// ms := mail.Service{Host: h.Configuration.MandrillHost}
		// m := mail.Mail{
		// 	Key: h.Configuration.MandrillApiKey,
		// 	Message: mail.Message{
		// 		Html:      fmt.Sprintf(mail.Template.Message, registration.Reference, registration.Reference),
		// 		Text:      "",
		// 		Subject:   mail.Template.Subject,
		// 		FromEmail: mail.Template.Sender.Email,
		// 		FromName:  mail.Template.Sender.Name,
		// 		To:        mail.Receivers{mail.Receiver{Email: partner.Email, Name: "", Type: "to"}},
		// 	},
		// 	Async: false,
		// }
		//
		// if err := ms.Send(m); err != nil {
		// 	logrus.Infof("Cannot send registration mail. Error: %s", err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		registration.Status = registrations.StatusVerificationSent
		if err := registrationsRepo.Update(registration); err != nil {
			logrus.Infof("Cannot update registration. Error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Please check your inbox!")
	} else {
		logrus.Info("Empty email.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
