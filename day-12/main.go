package main

import (
	"context"
	"day-12/connection"
	"day-12/middleware"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// ==================================================
// STRUCT TEMPLATE
// ==================================================

// PROJECT STRUCT
type Project struct {
	ID                     int
	ProjectName            string
	ProjectStartDate       time.Time
	ProjectEndDate         time.Time
	ProjectStartDateString string
	ProjectEndDateString   string
	ProjectDuration        string
	ProjectDescription     string
	ProjectTechnologies    []string
	ProjectImage           string
	UserID                 int
}

// ACCOUNT STRUCT
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// SESSION STRUCT
type Session struct {
	IsLogin   bool
	UserID    int
	UserName  string
	FlashData string
}

// ADDITIONAL
var Data = Session{}

// ==================================================
// MAIN FUNCTION
// ==================================================

// MAIN
func main() {
	route := mux.NewRouter()

	// CONNECTION TO DATABASE
	connection.DatabaseConnect()

	// ROUTE STATIC FOLDER
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	// ROUTE RENDER HTML
	route.HandleFunc("/", HomePage).Methods("GET")
	route.HandleFunc("/contact", ContactPage).Methods("GET")
	route.HandleFunc("/create-project", ProjectPage).Methods("GET")
	route.HandleFunc("/detail-project/{id}", ProjectDetail).Methods("GET")
	// ROUTE RENDER HTML FOR UPDATE
	route.HandleFunc("/update-project/{id}", EditProject).Methods("GET")
	// ROUTE RENDER HTML FOR REGISTER ACCOUNT
	route.HandleFunc("/register", Register).Methods("GET")
	// ROUTE RENDER HTML FOR LOGIN ACCOUNT
	route.HandleFunc("/login", Login).Methods("GET")

	// CREATE PROJECT
	route.HandleFunc("/created-project", middleware.UploadFile(CreateProject)).Methods("POST")
	// UPDATE PROJECT
	route.HandleFunc("/updated-project/{id}", middleware.UploadFile(UpdateProject)).Methods("POST")
	// DELETE PROJECT
	route.HandleFunc("/deleted-project/{id}", DeleteProject).Methods("GET")

	// REGISTER ACCOUNT
	route.HandleFunc("/registered", Registered).Methods("POST")
	// ACCOUNT LOGIN
	route.HandleFunc("/logged-in", LoggedIn).Methods("POST")
	// ACCOUNT LOGOUT
	route.HandleFunc("/logout", Logout).Methods("GET")

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
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
			// PROCESS FLASH MESSAGE
			fm := session.Flashes("message")

			var flashes []string
			// CHECKING IF FLASH MESSAGE IS AVAILABLE
			if len(fm) > 0 {
				session.Save(r, w)
				for _, f1 := range fm {
					flashes = append(flashes, f1.(string))
				}
			}
			// ADDING FLASH MESSAGE TO HTML
			Data.FlashData = strings.Join(flashes, "")
			Data.UserID = 0
		} else {
			// PROCESS FLASH MESSAGE
			fm := session.Flashes("message")

			var flashes []string
			// CHECKING IF FLASH MESSAGE IS AVAILABLE
			if len(fm) > 0 {
				session.Save(r, w)
				for _, f1 := range fm {
					flashes = append(flashes, f1.(string))
				}
			}
			// ADDING FLASH MESSAGE TO HTML
			Data.FlashData = strings.Join(flashes, "")
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)
		}

		var renderData []Project
		item := Project{}

		if !Data.IsLogin {
			// GET ALL PROJECTS DATA FROM POSTGRESQL
			rows, _ := connection.Conn.Query(context.Background(), `SELECT * FROM public.tb_project`)
			// PARSE PROJECT
			for rows.Next() {
				// CONNECT EACH ITEM WITH STRUCT
				err := rows.Scan(&item.ID, &item.ProjectName, &item.ProjectStartDate, &item.ProjectEndDate, &item.ProjectDescription, &item.ProjectTechnologies, &item.ProjectImage, &item.UserID)
				//ERROR HANDLING GET PROJECTS FROM POSTGRESQL
				if err != nil {
					fmt.Println(err.Error())
					return
				} else {
					// PARSING DATE
					item := Project{
						ID:                  item.ID,
						ProjectName:         item.ProjectName,
						ProjectDuration:     GetDuration(item.ProjectStartDate, item.ProjectEndDate),
						ProjectDescription:  item.ProjectDescription,
						ProjectTechnologies: item.ProjectTechnologies,
						ProjectImage:        item.ProjectImage,
						UserID:              item.UserID,
					}
					renderData = append(renderData, item)
				}
			}
			response := map[string]interface{}{
				"renderData": renderData,
				"Data":       Data,
			}
			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)
		} else {
			// GET USER PROJECTS FROM POSTGRESQL
			rows, _ := connection.Conn.Query(context.Background(), `SELECT * FROM public.tb_project WHERE "UserID" = $1`, Data.UserID)
			// PARSE PROJECT
			for rows.Next() {
				// CONNECT EACH ITEM WITH STRUCT
				err := rows.Scan(&item.ID, &item.ProjectName, &item.ProjectStartDate, &item.ProjectEndDate, &item.ProjectDescription, &item.ProjectTechnologies, &item.ProjectImage, &item.UserID)
				//ERROR HANDLING GET PROJECTS FROM POSTGRESQL
				if err != nil {
					fmt.Println(err.Error())
					return
				} else {
					// PARSING DATE
					item := Project{
						ID:                  item.ID,
						ProjectName:         item.ProjectName,
						ProjectDuration:     GetDuration(item.ProjectStartDate, item.ProjectEndDate),
						ProjectDescription:  item.ProjectDescription,
						ProjectTechnologies: item.ProjectTechnologies,
						ProjectImage:        item.ProjectImage,
						UserID:              item.UserID,
					}
					renderData = append(renderData, item)
				}
			}
			response := map[string]interface{}{
				"renderData": renderData,
				"Data":       Data,
			}
			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)
		}

	}
}

// RENDER CONTACT PAGE
func ContactPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)
		}
		response := map[string]interface{}{
			"Data": Data,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}

}

// RENDER CREATE PROJECT PAGE
func ProjectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/create-project.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)
		}
		response := map[string]interface{}{
			"Data": Data,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

// RENDER UPDATE PROJECT PAGE
func EditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/update-project.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		ID, _ := strconv.Atoi(mux.Vars(r)["id"])
		updateData := Project{}
		// GET PROJECT BY ID FROM POSTGRESQL
		err = connection.Conn.QueryRow(context.Background(), `SELECT *  FROM public.tb_project WHERE "ID" = $1`, ID).Scan(&updateData.ID, &updateData.ProjectName, &updateData.ProjectStartDate, &updateData.ProjectEndDate, &updateData.ProjectDescription, &updateData.ProjectTechnologies, &updateData.ProjectImage, &updateData.UserID)
		// ERROR HANDLING GET PROJECT FROM POSTGRESQL
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)
		}

		updateData = Project{
			ID:                     updateData.ID,
			ProjectName:            updateData.ProjectName,
			ProjectStartDateString: ReturnDate(updateData.ProjectStartDate),
			ProjectEndDateString:   ReturnDate(updateData.ProjectEndDate),
			ProjectDescription:     updateData.ProjectDescription,
			ProjectTechnologies:    updateData.ProjectTechnologies,
			ProjectImage:           updateData.ProjectImage,
			UserID:                 updateData.UserID,
		}
		response := map[string]interface{}{
			"updateData": updateData,
			"Data":       Data,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

// RENDER DETAIL PROJECT PAGE
func ProjectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/detail-project.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		ID, _ := strconv.Atoi(mux.Vars(r)["id"])
		renderDetail := Project{}
		// GET PROJECT BY ID FROM POSTGRESQL
		err = connection.Conn.QueryRow(context.Background(), `SELECT *
		FROM public.tb_project WHERE "ID" = $1`, ID).Scan(&renderDetail.ID, &renderDetail.ProjectName, &renderDetail.ProjectStartDate, &renderDetail.ProjectEndDate, &renderDetail.ProjectDescription, &renderDetail.ProjectTechnologies, &renderDetail.ProjectImage, &renderDetail.UserID)
		// ERROR HANDLING GET PROJECT DATA BY ID
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
		} else {
			// GET COOKIES DATA
			store := sessions.NewCookieStore([]byte("SESSION_KEY"))
			session, _ := store.Get(r, "SESSION_KEY")
			// CHECK LOGIN STATUS
			if session.Values["IsLogin"] != true {
				Data.IsLogin = false
			} else {
				Data.IsLogin = session.Values["IsLogin"].(bool)
				Data.UserName = session.Values["UserName"].(string)
				Data.UserID = session.Values["UserID"].(int)
			}
			//PARSING DATE
			renderDetail := Project{
				ID:                     renderDetail.ID,
				ProjectName:            renderDetail.ProjectName,
				ProjectStartDateString: FormatDate(renderDetail.ProjectStartDate),
				ProjectEndDateString:   FormatDate(renderDetail.ProjectEndDate),
				ProjectDuration:        GetDuration(renderDetail.ProjectStartDate, renderDetail.ProjectEndDate),
				ProjectDescription:     renderDetail.ProjectDescription,
				ProjectTechnologies:    renderDetail.ProjectTechnologies,
				ProjectImage:           renderDetail.ProjectImage,
				UserID:                 renderDetail.UserID,
			}
			response := map[string]interface{}{
				"renderDetail": renderDetail,
				"Data":         Data,
			}
			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)
		}
	}
}

// RENDER REGISTER PAGE
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/register.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
			// PROCESS FLASH MESSAGE
			fm := session.Flashes("message")

			var flashes []string
			// CHECKING IF FLASH MESSAGE IS AVAILABLE
			if len(fm) > 0 {
				session.Save(r, w)
				for _, f1 := range fm {
					flashes = append(flashes, f1.(string))
				}
			}
			// ADDING FLASH MESSAGE TO HTML
			Data.FlashData = strings.Join(flashes, "")

			response := map[string]interface{}{
				"Data": Data,
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)

		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)

			response := map[string]interface{}{
				"Data": Data,
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)
		}
	}
}

// RENDER LOGIN PAGE
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/login.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
			// PROCESS FLASH MESSAGE
			fm := session.Flashes("message")

			var flashes []string
			// CHECKING IF FLASH MESSAGE IS AVAILABLE
			if len(fm) > 0 {
				session.Save(r, w)
				for _, f1 := range fm {
					flashes = append(flashes, f1.(string))
				}
			}
			// ADDING FLASH MESSAGE TO HTML
			Data.FlashData = strings.Join(flashes, "")

			response := map[string]interface{}{
				"Data": Data,
			}
			tmpl.Execute(w, response)

		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)

			response := map[string]interface{}{
				"Data": Data,
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)
		}
	}
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
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)
		}

		ProjectName := r.PostForm.Get("project-name")
		ProjectStartDate := r.PostForm.Get("date-start")
		ProjectEndDate := r.PostForm.Get("date-end")
		ProjectDescription := r.PostForm.Get("project-description")
		ProjectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}
		UserID := Data.UserID

		dataContext := r.Context().Value("dataFile")
		ProjectImage := dataContext.(string)

		// INSERT PROJECT TO POSTGRESQL
		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.tb_project( "ProjectName", "ProjectStartDate", "ProjectEndDate", "ProjectDescription", "ProjectTechnologies", "ProjectImage", "UserID")
			VALUES ( $1, $2, $3, $4, $5, $6, $7)`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies, ProjectImage, UserID)
		// ERROR HANDLING INSERT PROJECT TO POSTGRESQL
		if err != nil {
			session.AddFlash("Cannot add project", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			session.AddFlash("Project added", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
	}
}

// UPDATE PROJECT
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		// GET COOKIES DATA
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// CHECK LOGIN STATUS
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserID = session.Values["UserID"].(int)
		}

		ID, _ := strconv.Atoi(mux.Vars(r)["id"])
		ProjectName := r.PostForm.Get("project-name")
		ProjectStartDate := r.PostForm.Get("date-start")
		ProjectEndDate := r.PostForm.Get("date-end")
		ProjectDescription := r.PostForm.Get("project-description")
		ProjectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("golang"), r.PostForm.Get("typescript")}

		dataContext := r.Context().Value("dataFile")
		ProjectImage := dataContext.(string)

		if ProjectImage == "empty" {
			// UPDATE PROJECT TO POSTGRESQL
			_, err = connection.Conn.Exec(context.Background(), `UPDATE public.tb_project
		SET "ProjectName"=$1, "ProjectStartDate"=$2, "ProjectEndDate"=$3, "ProjectDescription"=$4, "ProjectTechnologies"=$5 WHERE "ID"=$6`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies, ID)
		} else {
			// UPDATE PROJECT TO POSTGRESQL
			_, err = connection.Conn.Exec(context.Background(), `UPDATE public.tb_project
		SET "ProjectName"=$1, "ProjectStartDate"=$2, "ProjectEndDate"=$3, "ProjectDescription"=$4, "ProjectTechnologies"=$5, "ProjectImage" =$6 WHERE "ID"=$7`, ProjectName, ProjectStartDate, ProjectEndDate, ProjectDescription, ProjectTechnologies, ProjectImage, ID)
		}
		// ERROR HANDLING INSERT PROJECT TO POSTGRESQL
		if err != nil {
			fmt.Println(err)
			session.AddFlash("Cannot update project", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			session.AddFlash("Project updated", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
	}
}

// DELETE PROJECT
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(mux.Vars(r)["id"])

	// GET COOKIES DATA
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	// CHECK LOGIN STATUS
	if session.Values["IsLogin"] != true {
		session.AddFlash("Please login to continue", "message")
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	} else {
		// DELETE PROJECT BY ID AT POSTGRESQL
		_, err := connection.Conn.Exec(context.Background(), `DELETE FROM public.tb_project WHERE "ID" = $1`, ID)
		// ERROR HANDLING DELETE PROJECT AT POSTGRESQL
		if err != nil {
			session.AddFlash("Cannot delete project", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			session.AddFlash("Project deleted", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}

	}
}

// ==================================================
// ACCOUNT MANAGEMENT
// ==================================================

// REGISTER ACCOUNT
func Registered(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		Name := r.PostForm.Get("name-account")
		Email := r.PostForm.Get("email-account")
		Password := r.PostForm.Get("password-account")
		// ENCRYPT PASSWORD WITH BCRYPT
		PasswordHash, _ := bcrypt.GenerateFromPassword([]byte(Password), 10)

		// INSERT ACCOUNT TO POSTGRESQL
		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.tb_user("Name", "Email", "Password") VALUES($1, $2, $3)`, Name, Email, PasswordHash)
		//ERROR HANDLING INSERT ACCOUNT TO POSTGRESQL
		if err != nil {
			var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
			session, _ := store.Get(r, "SESSION_KEY")

			session.Values["IsLogin"] = false
			session.AddFlash("Email is registered! Please login", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		} else {
			var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
			session, _ := store.Get(r, "SESSION_KEY")

			session.Values["IsLogin"] = false
			session.AddFlash("Account created", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		}

	}
}

// ACCOUNT LOGIN
func LoggedIn(w http.ResponseWriter, r *http.Request) {
	// SETUP COOKIE STORE
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	} else {
		Email := r.PostForm.Get("email-account")
		Password := r.PostForm.Get("password-account")

		LoginUser := User{}
		// GET USER BY EMAIL FROM POSTGRESQL
		err = connection.Conn.QueryRow(context.Background(), `SELECT * FROM public.tb_user WHERE "Email" = $1`,
			Email).Scan(&LoginUser.ID, &LoginUser.Name, &LoginUser.Email, &LoginUser.Password)
		// ERROR HANDLING GET USER FROM POSTGRESQL (NO MATCHING EMAIL ADDRESS)
		if err != nil {
			var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
			session, _ := store.Get(r, "SESSION_KEY")

			session.Values["IsLogin"] = false
			session.AddFlash("Email not registered!", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
			return
		} else {
			// CHECK PASSWORD
			err = bcrypt.CompareHashAndPassword([]byte(LoginUser.Password), []byte(Password))
			// ERROR HANDLING
			if err != nil {
				var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
				session, _ := store.Get(r, "SESSION_KEY")

				session.Values["IsLogin"] = false
				session.AddFlash("Wrong password!", "message")
				session.Save(r, w)

				http.Redirect(w, r, "/login", http.StatusMovedPermanently)
				return
			} else {
				// CREATE SESSION CACHE
				session.Values["IsLogin"] = true
				session.Values["UserName"] = LoginUser.Name
				session.Values["UserID"] = LoginUser.ID

				// SETUP DURATION LOGGED IN
				session.Options.MaxAge = 10800 // 10800 SECONDS = 3 HOURS

				//ADDING FLASH MESSAGE
				session.AddFlash("Login successfully! Welcome "+LoginUser.Name, "message")
				session.Save(r, w)

				http.Redirect(w, r, "/", http.StatusMovedPermanently)
			}
		}
	}
}

// ACCOUNT LOGOUT
func Logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")
	// REMOVE SESSION
	session.Options.MaxAge = -1
	// ADDING LOGIN STATUS
	session.Values["IsLogin"] = false
	session.Values["UserID"] = 0
	// ADDING FLASH MESSAGE
	session.AddFlash("Successfully Logout!", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ==================================================
// ADDITIONAL FUNCTION
// ==================================================

// GET DURATION
func GetDuration(startDate time.Time, endDate time.Time) string {

	margin := endDate.Sub(startDate).Hours() / 24
	var duration string

	if margin >= 30 {
		if (margin / 30) == 1 {
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
func FormatDate(InputDate time.Time) string {

	Formated := InputDate.Format("02 January 2006")

	return Formated
}

// RETURN DATE FORMAT
func ReturnDate(InputDate time.Time) string {

	Formated := InputDate.Format("2006-01-02")

	return Formated
}
