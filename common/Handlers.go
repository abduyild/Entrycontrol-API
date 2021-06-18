package common

import (
	"log"
	"net/http"
	"pi-software/repos"
	"regexp"
	"strings"
	"text/template"
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
			if ok, _ := regexp.MatchString("\\d\\d:\\d\\d", date); !ok {
				http.Redirect(response, request, "/?wrong", http.StatusFound)
				return
			}
			date = formatDate(date)
			users, err := repos.GetEntriesForDate(mosque, date)
			if err != nil {
				// TODO this is not an error, this only shows that there is no entry for that date
				http.Redirect(response, request, "/?wrong", http.StatusFound)
				return
			} else {
				locations := getLocations(users)
				templateStruct := TemplateStruct{
					Users:     users,
					Date:      repos.GetCurrentDate(),
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
