package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"github.com/Suy56/GradeUpNow/internal/models"
	"github.com/gorilla/mux"
	//"strconv"
)
func (app *application) login(w http.ResponseWriter, r *http.Request) {
    // Check if it's a POST request
    if r.Method == http.MethodPost {
        // Retrieve the username and password from the form
        username := r.FormValue("username")
        password := r.FormValue("password")
        if username == "" || password == "" {
            // Invalid input, set error message
            errorMessage := "Username and password cannot be empty."
            http.Redirect(w, r, "/login?error="+errorMessage, http.StatusSeeOther)
            return
        }

        // Check if the username and password exist in the database
        user, err := app.user.Get(username)
        if err != nil {
            if errors.Is(err, models.ErrNoRecord) {
                errorMessage := "User not found."
                http.Redirect(w, r, "/login?error="+errorMessage, http.StatusSeeOther)
                return
            } else {
                app.serverError(w, err)
                return
            }
        }

        // Check if the entered password matches the user's password
        if user.Password != password {
            errorMessage := "Invalid password."
            http.Redirect(w, r, "/login?error="+errorMessage, http.StatusSeeOther)
            return
        }

        // Authentication successful, redirect to home page
        http.Redirect(w, r, "/sample", http.StatusSeeOther)
        return
    }

    // Render the login form
    files := []string{
        "./ui/html/login.html",
        "./ui/html/index.html",
    }
    tmpl, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    // Get the error message from the query parameter (if any)
    errorMessage := r.URL.Query().Get("error")
    data := struct {
        ErrorMessage string
    }{
        ErrorMessage: errorMessage,
    }

    err = tmpl.ExecuteTemplate(w, "login.html", data)
    if err != nil {
        app.serverError(w, err)
        return
    }
}
func (app *application)sample_home(w http.ResponseWriter, r *http.Request){
	tmpl,err:=template.ParseFiles("./ui/html/sample.html")
	if err!=nil{
		app.serverError(w,err)
		return
	}
	err=tmpl.ExecuteTemplate(w,"sample.html",nil)
}

func (app *application)home(w http.ResponseWriter, r *http.Request){
	files := []string{
		"./ui/html/index.html",
		"./ui/html/login.html",
	}
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	tmpl,err:=template.ParseFiles(files...)
	if err!=nil{
		app.notFound(w)
	}
	
	err=tmpl.ExecuteTemplate(w,"index.html",nil)
	if err!=nil{
		app.serverError(w,err)
	}
	
}
func (app *application) sign_up(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        // Retrieve form values
        email := r.FormValue("email")
        username := r.FormValue("username")
        pass := r.FormValue("password")

        if email == "" || username == "" || pass == "" {
            errorMessage := "Email, display name, and password cannot be empty."
            http.Redirect(w, r, "/signup?error="+errorMessage, http.StatusSeeOther)
            return
        }

        // Check if username already exists
        usernameExists, err := app.user.Check_if_exist(username, "")
        if err != nil {
            app.serverError(w, err)
            return
        }
        if usernameExists {
            errorMessage := "Username already exists."
            http.Redirect(w, r, "/signup?error="+errorMessage, http.StatusSeeOther)
            return
        }

        // Check if email already exists
        emailExists, err := app.user.Check_if_exist("", email)
        if err != nil {
            app.serverError(w, err)
            return
        }
        if emailExists {
            errorMessage := "Email already exists."
            http.Redirect(w, r, "/signup?error="+errorMessage, http.StatusSeeOther)
            return
        }

        // Proceed with user registration
        _, err = app.user.SignUp(username, email, pass)
        if err != nil {
            app.serverError(w, err)
            return
        }

        // Redirect to the login page after successful signup
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Render the signup form for GET requests
    files := []string{
        "./ui/html/signup.html",
        "./ui/html/index.html",
    }
    tmpl, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    // Get the error message from the query parameter (if any)
    errorMessage := r.URL.Query().Get("error")
    data := struct {
        ErrorMessage string
    }{
        ErrorMessage: errorMessage,
    }

    err = tmpl.ExecuteTemplate(w, "signup.html", data)
    if err != nil {
        app.serverError(w, err)
        return
    }
}

func (app *application)get_usr_stats(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)
	username:=vars["username"]
	fmt.Print(username)
	fmt.Println(username)
	user,err:=app.user.Get(username)
	if err!=nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(w)
		}else{
			app.serverError(w,err)
		}
		return
	}
	fmt.Fprintln(w, "username:", user.Username,"Theory:",user.Theory_score,"Mcq:",user.Mcq_score,"Total:",user.Total_score)
}

func (app *application)leader_board(w http.ResponseWriter, r *http.Request){
	leader_board,err:=app.user.Leader_board()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _,usr:=range leader_board{
		fmt.Fprintf(w,"%v\n",usr)
	}
}

func (app *application)q_type_handler(w http.ResponseWriter,r *http.Request){
	vars:=mux.Vars(r)
	subject:=vars["subject"]
	q_type:=vars["type"]
	if subject!="java" && subject!="fse" && subject!="dsa" &&subject!="interview"{
		fmt.Fprint(w,"Subject not available")
		return
	}
	switch q_type{
	case "mcq":
		mcq,err:=app.user.Get_Mcq(subject)
		if err!=nil{
			app.serverError(w,err)
			return
		}
		for _,question :=range mcq{
			fmt.Fprint(w,question.MQ_num,".",question.MQ_question,"\n1.",
		question.Option1,"\n2.",question.Option2,"\n3.",question.Option3,"\n4.",question.Option4,"\n\n")
		}
	case "theory":
		theory,err:=app.user.Get_Theory(subject)
		if err!=nil{
			app.serverError(w,err)
			return
		}
		for _,question:=range theory{
			fmt.Fprintf(w,"%v",question)
		}
	}
	
}

