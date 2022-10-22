package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// STRUCT TEMPLATE
type Project struct {
	ProjectName          string
	ProjectStartDate     string
	ProjectEndDate       string
	ProjectDuration      string
	ProjectDescription   string
	ProjectUseNodeJS     string
	ProjectUseReactJS    string
	ProjectUseGolang     string
	ProjectUseTypeScript string
}

// LOCAL DATABASE
var ProjectList = []Project{
	{
		ProjectName:          "Test Project Main",
		ProjectStartDate:     "01 October 2022",
		ProjectEndDate:       "01 December 2022",
		ProjectDuration:      "2 Months",
		ProjectDescription:   "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.",
		ProjectUseNodeJS:     "checked",
		ProjectUseReactJS:    "checked",
		ProjectUseGolang:     "checked",
		ProjectUseTypeScript: "checked",
	},
	{
		ProjectName:          "Test Project Additional",
		ProjectStartDate:     "20 October 2022",
		ProjectEndDate:       "21 November 2022",
		ProjectDuration:      "1 Months",
		ProjectDescription:   "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.",
		ProjectUseNodeJS:     "checked",
		ProjectUseReactJS:    "",
		ProjectUseGolang:     "checked",
		ProjectUseTypeScript: "",
	},
}

// MAIN
func main() {
	route := mux.NewRouter()

	// ROUTE PUBLIC FOLDER
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// ROUTE RENDER HTML
	route.HandleFunc("/", HomePage).Methods("GET")
	route.HandleFunc("/contact", ContactPage).Methods("GET")
	route.HandleFunc("/project", CreateProjectPage).Methods("GET")
	route.HandleFunc("/detail-project/{index}", ProjectDetail).Methods("GET")

	// CREATE PROJECT
	route.HandleFunc("/project/create", CreateProject).Methods("POST")
	// UPDATE PROJECT
	route.HandleFunc("/update-project/{index}", UpdateProject).Methods("GET")
	// DELETE PROJECT
	route.HandleFunc("/delete-project/{index}", DeleteProject).Methods("GET")

	// PORT HANDLING
	fmt.Println(("Server running on port 5000"))
	http.ListenAndServe("localhost:5000", route)
}

// RENDER HOME PAGE
func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	response := map[string]interface{}{
		"ProjectList": ProjectList,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)
}

// RENDER CONTACT PAGE
func ContactPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

// RENDER ADD PROJECT PAGE
func CreateProjectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

// RENDER PROJECT DETAIL
func ProjectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/project-detail.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		var renderDetail = Project{}
		index, _ := strconv.Atoi(mux.Vars(r)["index"])

		for i, data := range ProjectList {
			if index == i {
				renderDetail = Project{
					ProjectName:          data.ProjectName,
					ProjectStartDate:     data.ProjectStartDate,
					ProjectEndDate:       data.ProjectEndDate,
					ProjectDuration:      data.ProjectDuration,
					ProjectDescription:   data.ProjectDescription,
					ProjectUseNodeJS:     data.ProjectUseNodeJS,
					ProjectUseReactJS:    data.ProjectUseReactJS,
					ProjectUseGolang:     data.ProjectUseGolang,
					ProjectUseTypeScript: data.ProjectUseTypeScript,
				}
			}
		}
		data := map[string]interface{}{
			"renderDetail": renderDetail,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}

// CREATE PROJECT
func CreateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		projectName := r.PostForm.Get("project-name")
		projectStartDate := r.PostForm.Get("date-start")
		projectEndDate := r.PostForm.Get("date-end")
		projectDescription := r.PostForm.Get("project-description")
		projectUseNodeJS := r.PostForm.Get("nodejs")
		projectUseReactJS := r.PostForm.Get("reactjs")
		projectUseGolang := r.PostForm.Get("golang")
		projectUseTypeScript := r.PostForm.Get("typescript")

		var newProject = Project{
			ProjectName:          projectName,
			ProjectStartDate:     FormatDate(projectStartDate),
			ProjectEndDate:       FormatDate(projectEndDate),
			ProjectDuration:      GetDuration(projectStartDate, projectEndDate),
			ProjectDescription:   projectDescription,
			ProjectUseNodeJS:     projectUseNodeJS,
			ProjectUseReactJS:    projectUseReactJS,
			ProjectUseGolang:     projectUseGolang,
			ProjectUseTypeScript: projectUseTypeScript,
		}

		fmt.Println(newProject)

		ProjectList = append(ProjectList, newProject)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

// DELETE PROJECT
func DeleteProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	ProjectList = append(ProjectList[:index], ProjectList[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

// GET PROJECT TO EDIT
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/update-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		var updateData = Project{}
		index, _ := strconv.Atoi(mux.Vars(r)["index"])

		for i, data := range ProjectList {
			if index == i {
				updateData = Project{
					ProjectName:          data.ProjectName,
					ProjectStartDate:     ReturnDate(data.ProjectStartDate),
					ProjectEndDate:       ReturnDate(data.ProjectEndDate),
					ProjectDescription:   data.ProjectDescription,
					ProjectUseNodeJS:     data.ProjectUseNodeJS,
					ProjectUseReactJS:    data.ProjectUseReactJS,
					ProjectUseGolang:     data.ProjectUseGolang,
					ProjectUseTypeScript: data.ProjectUseTypeScript,
				}
				ProjectList = append(ProjectList[:index], ProjectList[index+1:]...)
			}
		}
		data := map[string]interface{}{
			"updateData": updateData,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, data)
	}
}

// ADDITIONAL FUNCTION

// GET DURATION
func GetDuration(startDate string, endDate string) string {

	layout := "2006-01-02"

	date1, _ := time.Parse(layout, startDate)
	date2, _ := time.Parse(layout, endDate)

	margin := date2.Sub(date1).Hours() / 24
	var duration string

	if margin > 30 {
		if (margin / 30) <= 1 {
			duration = "1 Month"
		} else {
			duration = strconv.Itoa(int(margin)/30) + " Months"
		}
	} else {
		if margin <= 1 {
			duration = "1 Day"
		} else {
			duration = strconv.Itoa(int(margin)) + " Days"
		}
	}

	return duration
}

// FORMAT DATE
func FormatDate(InputDate string) string {

	layout := "2006-01-02"
	t, _ := time.Parse(layout, InputDate)

	Formated := t.Format("02 January 2006")

	return Formated
}

// RETURN DATE FORMAT BEFORE EDIT
func ReturnDate(InputDate string) string {

	layout := "02 January 2006"
	t, _ := time.Parse(layout, InputDate)

	Formated := t.Format("2006-01-02")

	return Formated
}
