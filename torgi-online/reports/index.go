package reports

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Constants for departments
const (
	B24ManagerDepartment = 5
	B24ITDepartment      = 10
)

// Struct for user data
type User struct {
	ID       string
	Name     string
	LastName string
}

// Struct for task data
type Task struct {
	Count int
	Tasks []TaskDetail
}

// Struct for task detail
type TaskDetail struct {
	Link  string
	Title string
}

// PageData holds data for rendering HTML
type PageData struct {
	Department int
	Users      []User
	Start      string
	Finish     string
	Tasks      map[string]Task
}

func main() {
	http.HandleFunc("/", indexHandler)
	log.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form inputs
	r.ParseForm()
	department, _ := strconv.Atoi(r.FormValue("department"))
	start := r.FormValue("start")
	finish := r.FormValue("finish")

	if start == "" || finish == "" {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		if start == "" {
			start = yesterday
		}
		if finish == "" {
			finish = yesterday
		}
	}

	// Example users - replace with real data from Bitrix
	users := []User{
		{ID: "1", Name: "John", LastName: "Doe"},
		{ID: "2", Name: "Jane", LastName: "Smith"},
	}

	// Example tasks - replace with real logic
	tasks := map[string]Task{
		"1": {
			Count: 3,
			Tasks: []TaskDetail{
				{Link: "https://example.com/task/1", Title: "Task 1"},
				{Link: "https://example.com/task/2", Title: "Task 2"},
			},
		},
		"2": {
			Count: 2,
			Tasks: []TaskDetail{
				{Link: "https://example.com/task/3", Title: "Task 3"},
			},
		},
	}

	// Prepare data for rendering
	data := PageData{
		Department: department,
		Users:      users,
		Start:      start,
		Finish:     finish,
		Tasks:      tasks,
	}

	// Render template
	tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl.Execute(w, data)
}
