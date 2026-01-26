package handlers

import (
	"cv-landing/pkg/activity"
	"encoding/json"
	"fmt"
	"net/http"
)

type ActivityHandler struct {
	Repo activity.RepositoryHandler
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
	default:
		fmt.Fprintln(w, "ah! not found(")
	}
}

func (h *ActivityHandler) getGeneric(w http.ResponseWriter, activityType string) {
	activities, err := h.Repo.GetAllOfType(activityType)
	if err != nil {
		w.WriteHeader(501)
		w.Write([]byte(err.Error()))
		return
	}
	result, err := json.Marshal(activities)
	if err != nil {
		w.WriteHeader(502)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(result)
}
