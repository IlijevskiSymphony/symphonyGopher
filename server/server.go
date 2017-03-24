package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/IlijevskiSymphony/symphonyGopher/server/configuration"
	"github.com/IlijevskiSymphony/symphonyGopher/server/handlers"
	"github.com/IlijevskiSymphony/symphonyGopher/server/persistance"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

//Router function represents router for Server
func Router(conf configuration.Configuration, sessionFn func() *mgo.Session, store *sessions.CookieStore) *mux.Router {
	router := mux.NewRouter()

	router.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello and welcome!")
	})

	router.
		Methods("POST").
		PathPrefix("/register").
		Handler(handlers.HandlerRegisterPost{Configuration: conf, SessionFn: sessionFn})

	router.
		PathPrefix("/signup/accept/").
		Handler(handlers.HandlerRegisterVerificationGet{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("POST").
		Path("/verification/resend").
		Handler(handlers.HandlerRegisterVerificationResendPost{Configuration: conf, SessionFn: sessionFn})

	router.
		Methods("POST").
		Path("/resetPassword").
		Handler(handlers.HandlerResetPasswordPost{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("POST").
		Path("/resetPasswordEmail").
		Handler(handlers.HandlerResetPasswordEmailPost{Configuration: conf, SessionFn: sessionFn})

	router.
		PathPrefix("/resetPasswordCheck/").
		Handler(handlers.HandlerResetPasswordCheckGet{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("GET").
		Path("/login").
		Handler(handlers.HandlerLoginGet{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("POST").
		Path("/login").
		Handler(handlers.HandlerLoginPost{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("GET").
		Path("/logout").
		Handler(handlers.HandlerLogoutGet{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("POST").
		Path("/addLink").
		Handler(handlers.HandlerLinkAddPost{Configuration: conf, SessionFn: sessionFn, Store: store})

	return router
}

//Start function starts the web server
func Start() {
	var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(64)))

	store.Options = &sessions.Options{
		Domain:   "localhost",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		Path:     "/",
	}

	conf, errs := configuration.Read()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Fatalf("Invalid configuration: %s", err)
		}
		return
	}
	if print, err := json.MarshalIndent(conf, "", "    "); err == nil {
		logrus.Infof("%s\n", print)
	} else {
		logrus.Error("Cannot print environment variables.")
	}

	dashboardDS, err := persistance.NewDataStore(conf.SymphonyGopherDB)
	if err != nil {
		logrus.Fatalf("Cannot create developer dashboard data store. Error: %s.", err)
	}
	defer dashboardDS.Close()

	var sessionFn = func() *mgo.Session { return dashboardDS.Session() }

	logrus.Fatal(http.ListenAndServe(":"+conf.SymphonyGopherPort, Router(*conf, sessionFn, store)))
}
