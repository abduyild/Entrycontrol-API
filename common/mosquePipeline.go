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
	location := sanitize(query.Get("location"))
	if isValid(mosque, firstName, lastName, phone, address, location) && repos.DoesDBExist(mosque) {
		user := repos.StringToUser(firstName, lastName, phone, address, location)
		repos.PushToDB(mosque, repos.GetEncryptedUser(user, mosque))
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
