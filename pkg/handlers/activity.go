package handlers

import (
	"cv-landing/pkg/activity"
	"encoding/json"
	"net/http"
)

type ActivityHandler struct {
	Repo activity.ActivityHandler
}

func (h *ActivityHandler) Get(w http.ResponseWriter, r *http.Request) {
	activityType := r.PathValue("type")
	switch activityType {
	case "projects":
		h.getGeneric(w, "project")
	case "education":
		h.getGeneric(w, "education")
	case "events":
		h.getGeneric(w, "event")
	}
}

func (h *ActivityHandler) getGeneric(w http.ResponseWriter, activityType string) {
	activities, err := h.Repo.GetAllOfType(activityType)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	result, err := json.Marshal(activities)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(result)
}
