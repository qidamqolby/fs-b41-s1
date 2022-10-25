package main

import (
	"context"
	"day-10/connection"
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
	ID                  int
	ProjectName         string
	ProjectStartDate    string
	ProjectEndDate      string
	ProjectDuration     string
	ProjectDescription  string
	ProjectTechnologies []string
	ProjectImage        string
}

// LOCAL DATABASE
var ProjectList = []Project{}

// MAIN
func main() {
	route := mux.NewRouter()

	// CONNECTION TO DATABASE
	connection.DatabaseConnect()

	// ROUTE PUBLIC FOLDER
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// ROUTE RENDER HTML
	route.HandleFunc("/", HomePage).Methods("GET")
	route.HandleFunc("/contact", ContactPage).Methods("GET")
	route.HandleFunc("/project", ProjectPage).Methods("GET")
	route.HandleFunc("/detail-project/{id}", ProjectDetail).Methods("GET")
	// ROUTE RENDER HTML FOR UPDATE
	route.HandleFunc("/update-project/{id}", EditProject).Methods("GET")

	// CREATE PROJECT
	route.HandleFunc("/created-project", CreateProject).Methods("POST")
	// UPDATE PROJECT
	route.HandleFunc("/updated-project/{id}", UpdateProject).Methods("POST")
	// DELETE PROJECT
	route.HandleFunc("/deleted-project/{id}", DeleteProject).Methods("GET")

	// PORT HANDLING
	fmt.Println(("Server running on port 5000"))
	http.ListenAndServe("localhost:5000", route)
}

// ==================================================
// RENDER HTML TEMPLATE
// ==================================================

// RENDER HOME PAGE
func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// GET ALL PROJECTS FROM POSTGRESQL
		rows, _ := connection.Conn.Query(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDuration", "ProjectDescription", "ProjectTechnologies", "ProjectImage" FROM public.tb_project`)
		var renderData []Project
		// PARSE PROJECT
		for rows.Next() {

			item := Project{}

			err := rows.Scan(&item.ID, &item.ProjectName, &item.ProjectStartDate, &item.ProjectEndDate, &item.ProjectDuration, &item.ProjectDescription, &item.ProjectTechnologies, &item.ProjectImage)
			//ERROR HANDLING GET PROJECTS FROM POSTGRESQL
			if err != nil {
				fmt.Println(err.Error())
				return
			} else {
				renderData = append(renderData, item)
			}
		}
		response := map[string]interface{}{
			"renderData": renderData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
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

// RENDER PROJECT PAGE
func ProjectPage(w http.ResponseWriter, r *http.Request) {
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
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		ID, _ := strconv.Atoi(mux.Vars(r)["id"])
		renderDetail := Project{}
		// GET PROJECT BY ID FROM POSTGRESQL
		err = connection.Conn.QueryRow(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDuration", "ProjectDescription", "ProjectTechnologies", "ProjectImage"
		FROM public.tb_project WHERE "ID" = $1`, ID).Scan(&renderDetail.ID, &renderDetail.ProjectName, &renderDetail.ProjectStartDate, &renderDetail.ProjectEndDate, &renderDetail.ProjectDuration, &renderDetail.ProjectDescription, &renderDetail.ProjectTechnologies, &renderDetail.ProjectImage)
		// ERROR HANDLING GET PROJECT DATA BY ID
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
		} else {
			response := map[string]interface{}{
				"renderDetail": renderDetail,
			}
			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)
		}
	}
}

// RENDER PROJECT PAGE TO EDIT
func EditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/update-project.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	ID, _ := strconv.Atoi(mux.Vars(r)["id"])
	renderData := Project{}
	// GET PROJECT BY ID FROM POSTGRESQL
	err = connection.Conn.QueryRow(context.Background(), `SELECT "ID", "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies", "ProjectImage" FROM public.tb_project WHERE "ID" = $1`, ID).Scan(&renderData.ID, &renderData.ProjectName, &renderData.ProjectStartDate, &renderData.ProjectEndDate, &renderData.ProjectDescription, &renderData.ProjectTechnologies, &renderData.ProjectImage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	updateData := Project{
		ID:                  renderData.ID,
		ProjectName:         renderData.ProjectName,
		ProjectStartDate:    ReturnDate(renderData.ProjectStartDate),
		ProjectEndDate:      ReturnDate(renderData.ProjectEndDate),
		ProjectDescription:  renderData.ProjectDescription,
		ProjectTechnologies: renderData.ProjectTechnologies,
		ProjectImage:        renderData.ProjectImage,
	}
	response := map[string]interface{}{
		"updateData": updateData,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, response)

}

// ==================================================
// CRUD API
// ==================================================

// CREATE PROJECT
func CreateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		StartDate := r.PostForm.Get("date-start")
		EndDate := r.PostForm.Get("date-end")

		ProjectName := r.PostForm.Get("project-name")
		ProjectStartDate := FormatDate(StartDate)
		ProjectEndDate := FormatDate(EndDate)
		ProjectDuration := GetDuration(StartDate, EndDate)
		ProjectDescription := r.PostForm.Get("project-description")
		ProjectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}
		ProjectImage := r.PostForm.Get("upload-image")

		// INSERT PROJECT TO POSTGRESQL
		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.tb_project( "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDuration", "ProjectDescription", "ProjectTechnologies", "ProjectImage")
			VALUES ( $1, $2, $3, $4, $5, $6, $7)`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDuration, ProjectDescription, ProjectTechnologies, ProjectImage)
		// ERROR HANDLING INSERT PROJECT TO POSTGRESQL
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

// UPDATE PROJECT
func UpdateProject(w http.ResponseWriter, r *http.Request) {
}

// DELETE PROJECT
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(mux.Vars(r)["id"])

	// DELETE PROJECT BY ID AT POSTGRESQL
	_, err := connection.Conn.Exec(context.Background(), `DELETE FROM public.tb_project WHERE "ID" = $1`, ID)
	// ERROR HANDLING DELETE PROJECT AT POSTGRESQL
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// ADDITIONAL FUNCTION

// GET DURATION
func GetDuration(startDate string, endDate string) string {

	layout := "2006-01-02"

	date1, _ := time.Parse(layout, startDate)
	date2, _ := time.Parse(layout, endDate)

	margin := date2.Sub(date1).Hours() / 24
	var duration string

	if margin >= 30 {
		if (margin / 30) <= 1 {
			duration = "1 Month"
		} else {
			duration = strconv.Itoa(int(margin/30)) + " Months"
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

// CHANGE DATE FORMAT
func FormatDate(InputDate string) string {

	layout := "2006-01-02"
	t, _ := time.Parse(layout, InputDate)

	Formated := t.Format("02 January 2006")

	return Formated
}

// RETURN DATE FORMAT
func ReturnDate(InputDate string) string {

	layout := "02 January 2006"
	t, _ := time.Parse(layout, InputDate)

	Formated := t.Format("2006-01-02")

	return Formated
}
