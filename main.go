package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/mordredp/wololo/auth"
)

// Global variables
var appConfig Config
var appData Data

func main() {

	setWorkingDir()

	loadConfig()

	loadData()

	router := chi.NewRouter()
	//loginRouter := chi.NewRouter()

	authenticator := auth.New(
		120,
		auth.Static(appConfig.StaticPass),
		auth.LDAP(
			appConfig.LDAPAddr,
			appConfig.LDAPBaseDN,
			appConfig.LDAPBindUser,
			appConfig.LDAPBindPass),
	)

	router.Use(middleware.Logger)

	router.Use(authenticator.Identify)
	router.Use(authenticator.Clear)

	router.Get("/", renderHomePage)

	router.Post("/login", authenticator.Login)
	router.Get("/logout", authenticator.Logout)
	router.Get("/refresh", authenticator.Refresh)

	router.Route("/wake/{deviceName}", func(r chi.Router) { r.Get("/", wakeUpWithDeviceName) })

	router.Post("/data/save", saveData)
	router.Get("/data/get", getData)

	router.Get("/health", checkHealth)

	//router.PathPrefix(basePath + "/static/").Handler(http.StripPrefix(basePath+"/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	httpListen := appConfig.IP + ":" + strconv.Itoa(appConfig.Port)

	log.Printf("starting webserver on %q", httpListen)
	log.Fatal(http.ListenAndServe(httpListen, router))
}

func setWorkingDir() {

	thisApp, err := os.Executable()
	if err != nil {
		log.Fatalf("error determining the executable directory: %s", err)
	}
	appPath := filepath.Dir(thisApp)
	os.Chdir(appPath)
	log.Printf("set working directory: %q", appPath)

}

func loadConfig() {

	err := cleanenv.ReadConfig("config.json", &appConfig)
	if err != nil {
		log.Fatalf("error loading configuration file: %s", err)
	}
	log.Printf("configuration loaded")
}
