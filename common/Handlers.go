package common

import (
	"fmt"
	"log"
	"net/http"
	"pi-software/repos"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

type TemplateStruct struct {
	Users     []repos.User
	Date      string
	Locations []string
}

// Function for handling the input data. Checks the data and if everything is valid prepares the data
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	mosque := request.URL.Query().Get("mosqueid")
	date := request.URL.Query().Get("date")
	if mosque != "" && date != "" {
		if !repos.DoesDBExist(mosque) {
			http.Redirect(response, request, "/?wrong", http.StatusFound)
			return
		} else {
			if ok, _ := regexp.MatchString("\\d\\d\\d\\d-\\d\\d-\\d\\d", date); !ok {
				http.Redirect(response, request, "/?wrong", http.StatusFound)
				return
			}
			date = formatDate(date)
			users, err := repos.GetEntriesForDate(mosque, date)
			if err != nil {
				// TODO this is not an error, this only shows that there is no entry for that date
				log.Println("mosque-date")
				http.Redirect(response, request, "/?wrong", http.StatusFound)
				return
			} else {
				locations := getLocations(users)
				templateStruct := TemplateStruct{
					Users:     users,
					Date:      date,
					Locations: locations,
				}
				t, err := template.ParseFiles("templates/getRegistrations.gohtml", "templates/base.tmpl", "templates/footer.tmpl")
				if err != nil {
					log.Println(err)
					return
				}
				t.Execute(response, templateStruct)
			}
		}
	} else {
		t, _ := template.ParseFiles("templates/userlogin.gohtml", "templates/base.tmpl", "templates/footer.tmpl")
		t.Execute(response, nil)
	}
}

// Function for handling the input data. Checks the data and if everything is valid prepares the data
func MosqueHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	mosquename := request.URL.Query().Get("mosquename")
	password := request.URL.Query().Get("password")
	location := request.URL.Query().Get("location")
	mosqueid := request.URL.Query().Get("mosqueid")

	if isValid(mosquename, location, password) {
		if checkPassword(password) {
			mosqueid = strings.Replace(uuid.NewString(), "-", "", -1)
			if !repos.DoesDBExist(mosqueid) {
				mosque := repos.Mosque{Name: mosquename, Location: location}
				repos.AddMosque(mosqueid, mosque)
				http.Redirect(response, request, fmt.Sprintf("/addMosque?mosqueid=%v&success", mosqueid), http.StatusFound)
				return
			} else {
				http.Redirect(response, request, "/addMosque?exists", http.StatusFound)
				return
			}
		} else {
			http.Redirect(response, request, "/addMosque?wrong", http.StatusFound)
			return
		}
	} else {
		t, _ := template.ParseFiles("templates/addmosque.gohtml", "templates/base.tmpl", "templates/footer.tmpl")
		t.Execute(response, mosqueid)
	}
}

func formatDate(date string) string {
	dates := strings.Split(date, "-")
	return dates[2] + "-" + dates[1] + "-" + dates[0]
}

func getLocations(users []repos.User) []string {
	var locations []string
	for _, user := range users {
		location := user.Location
		if !contains(location, locations) {
			locations = append(locations, location)
		}
	}
	return locations
}

func contains(location string, locations []string) bool {
	for _, loc := range locations {
		if location == loc {
			return true
		}
	}
	return false
}
