package common

import (
	"net/http"
	"pi-software/repos"
	"regexp"
	"strings"
	"text/template"
)

type TemplateStruct struct {
	Users []repos.User
	Date  string
}

// Handler for Login Page used with POST by submitting Loginform
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	mosque := request.URL.Query().Get("mosqueid")
	date := request.URL.Query().Get("date")
	if mosque != "" && date != "" {
		if !repos.DoesDBExist(mosque) {
			http.Redirect(response, request, "/?wrong", 302)
		} else {
			if ok, _ := regexp.MatchString("\\d\\d\\d\\d-\\d\\d-\\d\\d", date); !ok {
				return
			}
			date = formatDate(date)
			users, err := repos.GetEntriesForDate(mosque, date)
			if err != nil {
				http.Redirect(response, request, "/?wrong", 302)
			} else {
				templateStruct := TemplateStruct{
					Users: users,
					Date:  repos.GetCurrentDate(),
				}
				t, _ := template.ParseFiles("templates/getRegistrations.gohtml", "templates/base.tmpl", "templates/footer.tmpl")
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
