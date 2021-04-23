package common

import (
	"net/http"
	"pi-software/repos"
)

func SubmitAttendant(response http.ResponseWriter, request *http.Request) {
	mosque := request.URL.Query().Get("mosque")
	firstName := request.URL.Query().Get("fname")
	lastName := request.URL.Query().Get("lname")
	phone := request.URL.Query().Get("phone")
	address := request.URL.Query().Get("address")
	time := request.URL.Query().Get("time")
	user := repos.StringToUser(firstName, lastName, phone, address, time)
	repos.PushToDB(mosque, user)
}
