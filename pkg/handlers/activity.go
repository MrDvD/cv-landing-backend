package handlers

import (
	"cv-landing-backend/pkg/activity"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ActivityHandler struct {
	Repo activity.ActivityRepository
}

func (h *ActivityHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityType, has := vars["type"]
	if !has {
		w.WriteHeader(400)
		return
	}
	switch activityType {
	case "projects":
		h.getGeneric(w, "project")
	case "education":
		h.getGeneric(w, "education")
	case "events":
		h.getGeneric(w, "event")
	default:
		w.WriteHeader(400)
	}
}

func (h *ActivityHandler) getGeneric(w http.ResponseWriter, activityType string) {
	activities, err := h.Repo.GetAllOfType(activityType)
	if hasError(w, err) {
		return
	}
	result, err := json.Marshal(activities)
	if hasError(w, err) {
		return
	}
	w.Write(result)
}

func (h *ActivityHandler) Add(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var rawActivity activity.Activity
	err := decoder.Decode(&rawActivity)
	if hasError(w, err) {
		return
	}
	newActivity, err := h.Repo.Add(rawActivity)
	if hasError(w, err) {
		return
	}
	result, err := json.Marshal(newActivity)
	if hasError(w, err) {
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (h *ActivityHandler) Remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawActivityId, has := vars["id"]
	if !has {
		w.WriteHeader(400)
		return
	}
	activityId, err := strconv.Atoi(rawActivityId)
	if hasError(w, err) {
		return
	}
	err = h.Repo.Remove(activityId)
	if hasError(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
