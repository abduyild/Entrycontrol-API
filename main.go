package main

import (
	"log"
	"net/http"
	"pi-software/common"
	"pi-software/repos"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

// Functions for handling pagecalls like localhost:8080/login
func main() {
	if err := repos.InitDB(); err != nil {
		log.Fatal("Error initializing the Database, error:" + err.Error())
		return
	}
	router.HandleFunc("/", common.LoginHandler)
	router.HandleFunc("/addUser", common.SubmitAttendant)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("icons"))))
	http.Handle("/", router)
	log.Println("All handlers set and ready to listen")
	log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/camii.online/fullchain.pem", "/etc/letsencrypt/live/camii.online/privkey.pem", nil))
}
