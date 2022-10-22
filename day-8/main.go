package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// struct
type Project struct {
	ProjectID           int
	ProjectName         string
	ProjectStartDate    string
	ProjectEndDate      string
	ProjectDescription  string
	ProjectTechnologies []string
}

// local database
var ProjectList = []Project{
	// dummy data
	{
		ProjectID:           0,
		ProjectName:         "Test Project Main",
		ProjectStartDate:    "2022-10-20",
		ProjectEndDate:      "2022-10-31",
		ProjectDescription:  "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.",
		ProjectTechnologies: []string{"nodejs", "reactjs"},
	},
	{
		ProjectID:           1,
		ProjectName:         "Test Project Additional",
		ProjectStartDate:    "2022-10-20",
		ProjectEndDate:      "2022-10-31",
		ProjectDescription:  "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.",
		ProjectTechnologies: []string{"nodejs", "reactjs"},
	},
}

func main() {
	route := mux.NewRouter()

	// route path folder for public folder
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// routing
	route.HandleFunc("/", homePage).Methods("GET")
	route.HandleFunc("/contact", contactPage).Methods("GET")
	route.HandleFunc("/project", projectPage).Methods("GET")
	route.HandleFunc("/create-project", createProject).Methods("POST")
	route.HandleFunc("/delete-project/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/update-project/{id}", updateProject).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")

	fmt.Println(("Server running on port 5000"))
	http.ListenAndServe("localhost:5000", route)
}

// function route home page
func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// create render home project
	responseData := map[string]interface{}{
		"ProjectList": ProjectList,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, responseData)
}

// function route contact page
func contactPage(w http.ResponseWriter, r *http.Request) {
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
func projectPage(w http.ResponseWriter, r *http.Request) {
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
	var projectUseNodeJS = r.PostForm.Get("nodejs")
	var projectUseReactJS = r.PostForm.Get("reactjs")
	var projectUseGolang = r.PostForm.Get("golang")
	var projectUseTypeScript = r.PostForm.Get("typescript")

	var newProject = Project{
		ProjectName:         projectName,
		ProjectStartDate:    projectStartDate,
		ProjectEndDate:      projectEndDate,
		ProjectDescription:  projectDescription,
		ProjectTechnologies: []string{projectUseNodeJS, projectUseReactJS, projectUseGolang, projectUseTypeScript},
	}
	ProjectList = append(ProjectList, newProject)

	fmt.Println(ProjectList)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func GenerateProjectID() {

}

// function delete project in local database
func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectList = append(ProjectList[:index], ProjectList[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

// function update project in local database
func updateProject(w http.ResponseWriter, r *http.Request) {
	// work in progress
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	// work in progress
}
