package common

import (
	"net/http"
	"os"
	"pi-software/repos"
	"strings"

	b64 "encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func SubmitAttendant(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	mosqueid := sanitize(query.Get("mosque"))
	firstName := sanitize(query.Get("fname"))
	lastName := sanitize(query.Get("lname"))
	phone := sanitize(query.Get("phone"))
	address := sanitize(query.Get("address"))
	location := sanitize(query.Get("location"))
	if isValid(mosqueid, firstName, lastName, phone, address, location) && repos.DoesDBExist(mosqueid) {
		user := repos.StringToUser(firstName, lastName, phone, address, location)
		repos.PushToDB(mosqueid, repos.GetEncryptedUser(user, mosqueid))
	}
}

func checkPassword(password string) bool {
	decodedPassword, _ := b64.StdEncoding.DecodeString(os.Getenv("addmosquepw"))
	err := bcrypt.CompareHashAndPassword(decodedPassword, []byte(password))
	return err == nil
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
