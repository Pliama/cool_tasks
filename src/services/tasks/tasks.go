//Package tasks implements task handlers
package tasks

import (
	"net/http"
	"time"

	"github.com/Nastya-Kruglikova/cool_tasks/src/models"
	"github.com/Nastya-Kruglikova/cool_tasks/src/services/common"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"

	"github.com/Nastya-Kruglikova/cool_tasks/src/services/auth"
)

type successCreate struct {
	Status string      `json:"message"`
	Result models.Task `json:"result"`
}

type successDelete struct {
	Status string `json:"message"`
}

//GetTasksHandler gets Tasks from DB
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	if auth.CheckPermission(r, "admin", "") == false {
		common.SendError(w, r, http.StatusForbidden, "Wrong user role", nil)
		return
	}
	tasks, err := models.GetTasks()

	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't get tasks", err)
		return
	}

	common.RenderJSON(w, r, tasks)
}

//GetTaskHandler gets Task from DB by taskID
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	taskID, err := uuid.FromString(params["id"])

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong task ID (can't convert string to uuid)", err)
		return
	}

	task, err := models.GetTask(taskID)
	itemOwner, err := models.GetUserByID(task.UserID)
	if auth.CheckPermission(r, "owner", itemOwner.Login) == false {
		common.SendError(w, r, http.StatusForbidden, "Wrong user role", nil)
		return
	}
	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't get task by ID", err)
		return
	}

	common.RenderJSON(w, r, task)
}

//AddTaskHandler creates and saves Task in DB
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	var newTask models.Task
	var resultTask models.Task

	err := r.ParseForm()

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Can't parse POST Body", err)
		return
	}

	timeNow := time.Now()
	userID, err := uuid.FromString(r.Form.Get("user_id"))

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong User ID", err)
		return
	}

	newTask.UserID = userID
	newTask.Name = r.Form.Get("name")
	newTime := r.Form.Get("time")
	newTask.CreatedAt = timeNow
	newTask.UpdatedAt = timeNow
	newTask.Desc = r.Form.Get("desc")

	parsedTime, err := time.Parse(time.UnixDate, newTime)

	if err != nil {
		common.SendUnsupportedMediaType(w, r, "ERROR: Wrong date(can't convert string to int)", err)
		return
	}

	newTask.Time = parsedTime

	resultTask, err = models.AddTask(newTask)

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Can't add new task", err)
		return
	}

	common.RenderJSON(w, r, successCreate{Status: "201 Created", Result: resultTask})
}

//DeleteTaskHandler deletes Task from DB
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	taskID, err := uuid.FromString(params["id"])

	if err != nil {
		common.SendBadRequest(w, r, "ERROR: Wrong task ID (can't convert string to uuid)", err)
		return
	}

	err = models.DeleteTask(taskID)

	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't delete this task", err)
		return
	}

	common.RenderJSON(w, r, successDelete{Status: "204 No Content"})
}

//GetUserTasksHandler gets Tasks related to current User
func GetUserTasksHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	idUser, err := uuid.FromString(params["id"])
	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't get user", err)
		return
	}

	tasks, err := models.GetUserTasks(idUser)
	if err != nil {
		common.SendNotFound(w, r, "ERROR: Can't get user", err)
		return
	}

	itemOwner, err := models.GetUserByID(idUser)
	if auth.CheckPermission(r, "owner", itemOwner.Login) == false {
		common.SendError(w, r, http.StatusForbidden, "Wrong user role", nil)
		return
	}
	common.RenderJSON(w, r, tasks)
}
