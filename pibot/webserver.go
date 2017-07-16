package pibot

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

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

// StartWebServer serves the main web endpoints
func StartWebServer() {
	r := mux.NewRouter()
	registerDashboard(r)
	RegisterAPI(r)

	s := GetSettings()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.HTTPPort), r))
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
		Version:          Version,
		Favicon:          "/content/img/favicon.png",
		ControllerMethod: controllerMethod,
	}
}

func overviewHandler(w http.ResponseWriter, r *http.Request) {
	p := getDefaultPageData("overview")

	type data struct {
		Metrics  []Metric
		CPUCount int
	}

	startTime := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
	metrics := GetHostMetricsByTime(startTime, nil)
	info := GetHostInfo()

	p.Data = data{
		Metrics:  metrics,
		CPUCount: len(info.Processors),
	}

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
		sf1, _ := strconv.Atoi(r.FormValue("sensor_front_1"))
		sf2, _ := strconv.Atoi(r.FormValue("sensor_front_2"))
		sb1, _ := strconv.Atoi(r.FormValue("sensor_back_1"))
		sb2, _ := strconv.Atoi(r.FormValue("sensor_back_2"))

		settings := Settings{
			HTTPPort:    port,
			MotorLeft:   [2]int{ml1, ml2},
			MotorRight:  [2]int{mr1, mr2},
			SensorFront: [2]int{sf1, sf2},
			SensorBack:  [2]int{sb1, sb2},
		}

		settings.Save()
	}

	settings := GetSettings()
	p.Data = settings

	templates := template.Must(template.ParseFiles(templatePath+"layout.html", templatePath+"settings.html"))
	templates.ExecuteTemplate(w, "layout", p)
}