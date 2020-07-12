package main

import (
	"columnist_go/models"
	"columnist_go/routes"
	"columnist_go/utils"
	"net/http"
)

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
