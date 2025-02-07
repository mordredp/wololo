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

	"github.com/mordredp/auth"
	"github.com/mordredp/auth/provider"
	"github.com/mordredp/auth/provider/ldap"
)

var (
	appConfig Config
	appData   Data
)

func main() {

	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("cannot determine the executable directory: %s", err)
	}
	workingDir := filepath.Dir(executablePath)
	os.Chdir(workingDir)
	log.Printf("set working directory: %q", workingDir)

	err = cleanenv.ReadConfig("config.json", &appConfig)
	if err != nil {
		log.Fatalf("error loading configuration file: %s", err)
	}
	log.Printf("configuration loaded")

	loadData()

	router := chi.NewRouter()
	authRouter := chi.NewRouter()

	authenticator := auth.NewAuthenticator(
		auth.MaxSessionLength(appConfig.MaxSessionSeconds),
		auth.Static(appConfig.StaticPass),
		auth.LDAP(
			appConfig.LDAPAddr,
			appConfig.LDAPBaseDN,
			provider.Credentials{
				Username: appConfig.LDAPBindUser,
				Password: appConfig.LDAPBindPass,
			},
			ldap.IdKey(appConfig.LDAPIdKey),
			ldap.QueryParams(appConfig.LDAPQueryParams),
		),
	)

	router.Use(middleware.Logger)

	router.Use(authenticator.Identify)

	authRouter.Use(authenticator.Authorize)

	router.Get("/", renderHomePage)
	router.Post("/login", authenticator.Login)
	router.Get("/health", checkHealth)

	authRouter.Get("/logout", authenticator.Logout)

	authRouter.Route("/wake/{deviceName}", func(r chi.Router) { r.Get("/", wakeUpWithDeviceName) })

	authRouter.Post("/data/save", saveData)
	authRouter.Get("/data/get", getData)

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	})
	authRouter.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.Mount("/", authRouter)

	address := appConfig.IP + ":" + strconv.Itoa(appConfig.Port)

	log.Printf("starting webserver on %q", address)
	log.Fatal(http.ListenAndServe(address, router))
}
