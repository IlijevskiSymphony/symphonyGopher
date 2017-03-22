package server

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/Sirupsen/logrus"
	"github.com/IlijevskiSymphony/symphonyGopher/server/configuration"
	"github.com/IlijevskiSymphony/symphonyGopher/server/handlers"
	"github.com/IlijevskiSymphony/symphonyGopher/server/persistance"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//Router function represents router for Server
func Router(conf configuration.Configuration, sessionFn func() *mgo.Session, store *sessions.CookieStore) *mux.Router {
	router := mux.NewRouter()

	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(conf.StaticContentDir))))

	router.
		PathPrefix("/signup/accept/").
		Handler(handlers.HandlerRegisterVerificationGet{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		PathPrefix("/resetPasswordCheck/").
		Handler(handlers.HandlerResetPasswordCheckGet{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("GET").
		Path("/login").
		Handler(handlers.HandlerLoginGet{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("GET").
		Path("/logout").
		Handler(handlers.HandlerLogoutGet{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("POST").
		Path("/login").
		Handler(handlers.HandlerLoginPost{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("POST").
		Path("/verification/resend").
		Handler(handlers.HandlerRegisterVerificationResendPost{Configuration: conf, SessionFn: sessionFn})

	router.
		Methods("POST").
		PathPrefix("/register").
		Handler(handlers.HandlerRegisterPost{Configuration: conf, SessionFn: sessionFn})

	router.
		Methods("POST").
		Path("/settings").
		Handler(handlers.HandlerSettingsPost{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("POST").
		Path("/resetPasswordEmail").
		Handler(handlers.HandlerResetPasswordEmailPost{Configuration: conf, SessionFn: sessionFn})

	router.
		Methods("POST").
		Path("/resetPassword").
		Handler(handlers.HandlerResetPasswordPost{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.
		Methods("POST").
		Path("/addApplication").
		Handler(handlers.HandlerApplicationAddPost{Configuration: conf, SessionFn: sessionFn, Store: store})

	router.PathPrefix("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, conf.StaticContentDir+"/index.html")
	})

	return router
}

//Start function starts the web server
func Start() {
	// var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(64)))
	var store = sessions.NewCookieStore([]byte("something-very-secret-xx"))

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

	dashboardDS, err := persistance.NewDataStore(conf.DeveloperDashboardDB)
	if err != nil {
		logrus.Fatalf("Cannot create developer dashboard data store. Error: %s.", err)
	}
	defer dashboardDS.Close()

	var sessionFn = func() *mgo.Session { return dashboardDS.Session() }

	logrus.Fatal(http.ListenAndServe(":"+conf.Port, Router(*conf, sessionFn, store)))
}
