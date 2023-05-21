package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/mordredp/auth"
)

func renderHomePage(w http.ResponseWriter, r *http.Request) {

	pageData := struct {
		Devices []Device
		BCastIP string
	}{
		Devices: appData.Devices,
		BCastIP: appConfig.BCastIP,
	}

	var user auth.User
	if v, ok := r.Context().Value(auth.UserKey).(auth.User); ok {
		user = v
	}

	if !user.Authenticated {
		tpl := template.Must(template.ParseGlob("templates/*.gohtml"))
		tpl.ExecuteTemplate(w, "index.gohtml", user)
	} else {
		tmpl, _ := template.ParseFiles("index.html")
		tmpl.Execute(w, pageData)
	}
}

func checkHealth(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "alive")

}
