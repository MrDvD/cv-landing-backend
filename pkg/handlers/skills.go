package handlers

import (
	"cv-landing-backend/pkg/files"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

type SkillsHandler struct {
	Repo files.FileRepository
}

func (h *SkillsHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	skillsType, has := vars["type"]
	if !has {
		w.WriteHeader(400)
		return
	}
	switch skillsType {
	case "hard":
		h.getGeneric(w, append(h.Repo.BasePath(), "hardskills.json")...)
	case "soft":
		h.getGeneric(w, append(h.Repo.BasePath(), "softskills.json")...)
	default:
		w.WriteHeader(400)
	}
}

func (h *SkillsHandler) getGeneric(w http.ResponseWriter, path ...string) {
	file, err := h.Repo.Get(filepath.Join(path...))
	if hasError(w, err) {
		return
	}
	w.Write([]byte(file.Content))
}
