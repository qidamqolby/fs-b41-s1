package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// interface
var data = map[string]interface{}{
	"Title": "Personal Web",
}

// struct
type Project struct {
	ProjectName          string
	ProjectStartDate     string
	ProjectEndDate       string
	ProjectDescription   string
	ProjectUseNodeJS     string
	ProjectUseReactJS    string
	ProjectUseGolang     string
	ProjectUseJavaScript string
}

// local database
var ProjectList = []Project{
	// dummy data
	{
		ProjectName:          "Test Project Main",
		ProjectStartDate:     "20 October 2022",
		ProjectEndDate:       "31 October 2022",
		ProjectDescription:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
		ProjectUseNodeJS:     "on",
		ProjectUseReactJS:    "on",
		ProjectUseGolang:     "on",
		ProjectUseJavaScript: "on",
	},
}

func main() {
	route := mux.NewRouter()

	// route path folder for public folder
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// routing
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project", projectForm).Methods("GET")
	route.HandleFunc("/create-project", createProject).Methods("POST")

	fmt.Println(("Server running on port 5000"))
	http.ListenAndServe("localhost:5000", route)
}

// function route home page
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	// create render home project

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

// function route contact page
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

// function route project page
func projectForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

// function create project and adding to local database
func createProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var projectName = r.PostForm.Get("project-name")
	var projectStartDate = r.PostForm.Get("date-start")
	var projectEndDate = r.PostForm.Get("date-end")
	var projectDescription = r.PostForm.Get("project-description")
	var projectUseNodeJS = r.PostForm.Get("node-js")
	var projectUseReactJS = r.PostForm.Get("react-js")
	var projectUseGolang = r.PostForm.Get("golang")
	var projectUseJavaScript = r.PostForm.Get("javascript")

	var newProject = Project{
		ProjectName:          projectName,
		ProjectStartDate:     projectStartDate,
		ProjectEndDate:       projectEndDate,
		ProjectDescription:   projectDescription,
		ProjectUseNodeJS:     projectUseNodeJS,
		ProjectUseReactJS:    projectUseReactJS,
		ProjectUseGolang:     projectUseGolang,
		ProjectUseJavaScript: projectUseJavaScript,
	}
	ProjectList = append(ProjectList, newProject)

	fmt.Println(ProjectList)

	http.Redirect(w, r, "/project", http.StatusMovedPermanently)

}

// function delete project in local database

// function update project in local database
