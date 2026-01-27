package handlers

import (
	"cv-landing/pkg/activity"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ActivityHandler struct {
	Repo activity.ActivityHandler
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
