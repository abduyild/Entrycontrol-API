package common

import (
	"net/http"
	"pi-software/repos"
	"strings"
)

func SubmitAttendant(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	mosque := sanitize(query.Get("mosque"))
	firstName := sanitize(query.Get("fname"))
	lastName := sanitize(query.Get("lname"))
	phone := sanitize(query.Get("phone"))
	address := sanitize(query.Get("address"))
	time := sanitize(query.Get("time"))
	location := sanitize(query.Get("location"))
	if isValid(mosque, firstName, lastName, phone, address, time, location) {
		user := repos.StringToUser(firstName, lastName, phone, address, time, location)
		repos.PushToDB(mosque, user)
	}
}

func isValid(values ...string) bool {
	for _, value := range values {
		if value == "" {
			return false
		}
	}
	return true
}

func sanitize(input string) string {
	return strings.ReplaceAll(input, "+", " ")
}
