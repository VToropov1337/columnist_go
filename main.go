package main

import (
	"columnist_go/models"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)


var templates *template.Template
var store = sessions.NewCookieStore([]byte("SESSION_KEY_s3cr3t"))

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["username"]
		fmt.Println("authrequired ===>>>>>",session.Values, session.Name())
		if !ok {
			http.Redirect(w,r,"/login",302)
			return
		}
		handler.ServeHTTP(w,r)
	}
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := models.GetComments()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	templates.ExecuteTemplate(w, "index.html",comments)
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("textcomment")
	err := models.PostComment(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return

	}
	http.Redirect(w,r,"/", 302)

}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html",nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.AuthenticateUser(username,password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			templates.ExecuteTemplate(w, "login.html","unknown user")
		case models.ErrInvalidLogin:
			templates.ExecuteTemplate(w, "login.html", "invalid login")
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}
		return
	}
	session, _ := store.Get(r, "session")
	fmt.Println("loginposthandler_session----->>",session.Name(),session.Values)
	session.Values["username"] = username
	session.Options.MaxAge = 5
	session.Save(r,w)
	http.Redirect(w,r,"/",302)
}

func testGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r,"session")
	fmt.Println("testhandler_===>",session.Values)
	untyped, ok := session.Values["username"]
	if !ok {
		return
	}
	username, ok := untyped.(string)
	if !ok {
		return
	}
	w.Write([]byte(username))
}


func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register.html",nil)

}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username,password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	http.Redirect(w,r,"/login",302)

}


func main() {
	models.Init()
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r:=mux.NewRouter()
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/",fs))
	r.HandleFunc("/",AuthRequired(indexGetHandler)).Methods("GET")
	r.HandleFunc("/",AuthRequired(indexPostHandler)).Methods("POST")
	r.HandleFunc("/login",loginGetHandler).Methods("GET")
	r.HandleFunc("/login",loginPostHandler).Methods("POST")
	r.HandleFunc("/register",registerGetHandler).Methods("GET")
	r.HandleFunc("/register", registerPostHandler).Methods("POST")
	r.HandleFunc("/test",testGetHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8080",nil)
}
