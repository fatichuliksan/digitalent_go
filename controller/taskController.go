package controller

import (
	"dts-task/model"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type NewTaskController struct {
	Db *gorm.DB
}

func (t *NewTaskController) Index(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
		"dateFormat": func(datetime time.Time) string {
			return datetime.Format("2006-01-02")
		},
	}

	var mainLayout = path.Join("views", "layouts", "main.html")
	var filepath = path.Join("views", "task", "index.html")
	var tmpl, err = template.New("").Funcs(funcMap).ParseFiles(filepath, mainLayout)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tasks []model.Task
	t.Db.Find(&tasks)

	var data = map[string]interface{}{
		"title": "Task",
		"data":  tasks,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *NewTaskController) Create(w http.ResponseWriter, r *http.Request) {
	var mainLayout = path.Join("views", "layouts", "main.html")
	var filepath = path.Join("views", "task", "create.html")
	var tmpl, err = template.New("").ParseFiles(filepath, mainLayout)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = map[string]interface{}{
		"title": "Create Task",
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *NewTaskController) Edit(w http.ResponseWriter, r *http.Request) {

	funcMap := template.FuncMap{
		"selectedStatusPending": func(val string) string {
			if val == "Pending" {
				return "selected"
			}
			return ""
		},
		"selectedStatusDone": func(val string) string {
			if val == "Done" {
				return "selected"
			}
			return ""
		},
	}

	var mainLayout = path.Join("views", "layouts", "main.html")
	var filepath = path.Join("views", "task", "edit.html")
	var tmpl, err = template.New("").Funcs(funcMap).ParseFiles(filepath, mainLayout)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/edit/")
	var tasks model.Task
	t.Db.Find(&tasks, id)

	var data = map[string]interface{}{
		"title":       "Edit Task",
		"assignee":    tasks.Assignee,
		"description": tasks.Description,
		"status":      tasks.Status,
		"deadlineAt":  tasks.DeadlineAt.Format("2006-01-02"),
		"id":          id,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *NewTaskController) Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		description := r.FormValue("description")
		assignee := r.FormValue("assignee")
		deadline := r.FormValue("deadline")
		status := r.FormValue("status")

		date, err := time.Parse("2006-01-02", deadline)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data := model.Task{}
		data.Description = description
		data.Assignee = assignee
		data.DeadlineAt = date
		data.Status = status

		t.Db.Create(&data)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (t *NewTaskController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		description := r.FormValue("description")
		assignee := r.FormValue("assignee")
		deadline := r.FormValue("deadline")
		status := r.FormValue("status")
		id := r.FormValue("id")

		date, err := time.Parse("2006-01-02", deadline)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tempId, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		data := model.Task{}
		data.ID = uint(tempId)
		data.Description = description
		data.Assignee = assignee
		data.DeadlineAt = date
		data.Status = status
		t.Db.Updates(&data)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (t *NewTaskController) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	t.Db.Delete(&model.Task{}, id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
