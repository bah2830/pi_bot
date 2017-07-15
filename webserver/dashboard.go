package webserver

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/bah2830/pi_bot/pibot"
	"github.com/gorilla/mux"
)

var templatePath = "resources/web_templates/"

type page struct {
	Title            string
	Version          string
	Favicon          string
	ControllerMethod string
	Data             interface{}
}

func registerDashboard(r *mux.Router) {
	r.HandleFunc("/", overviewHandler)
	r.HandleFunc("/control", controlHandler)
	r.HandleFunc("/settings", settingsHandler)

	// Setup file server for html resources
	s := http.StripPrefix("/content/", http.FileServer(http.Dir("./resources/web_content/")))
	r.PathPrefix("/content/").Handler(s)
	http.Handle("/", r)
}

func getDefaultPageData(controllerMethod string) page {
	return page{
		Title:            "PiBot",
		Version:          "0.1-pre-alpha",
		Favicon:          "/content/img/favicon.png",
		ControllerMethod: controllerMethod,
	}
}

func overviewHandler(w http.ResponseWriter, r *http.Request) {
	p := getDefaultPageData("overview")

	templates := template.Must(template.ParseFiles(templatePath+"layout.html", templatePath+"overview.html"))
	templates.ExecuteTemplate(w, "layout", p)
}

func controlHandler(w http.ResponseWriter, r *http.Request) {
	p := getDefaultPageData("control")

	templates := template.Must(template.ParseFiles(templatePath+"layout.html", templatePath+"control.html"))
	templates.ExecuteTemplate(w, "layout", p)
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	p := getDefaultPageData("settings")

	if r.Method == "POST" {
		r.ParseForm()
		port, _ := strconv.Atoi(r.FormValue("http_port"))
		ml1, _ := strconv.Atoi(r.FormValue("motor_left_1"))
		ml2, _ := strconv.Atoi(r.FormValue("motor_left_2"))
		mr1, _ := strconv.Atoi(r.FormValue("motor_right_1"))
		mr2, _ := strconv.Atoi(r.FormValue("motor_right_2"))

		settings := pibot.Settings{
			HTTPPort:    port,
			MotorLeft1:  ml1,
			MotorLeft2:  ml2,
			MotorRight1: mr1,
			MotorRight2: mr2,
		}

		settings.Save()
	}

	settings := pibot.GetSettings()
	p.Data = settings

	templates := template.Must(template.ParseFiles(templatePath+"layout.html", templatePath+"settings.html"))
	templates.ExecuteTemplate(w, "layout", p)
}
