package handlers

import (
	"cv-landing/pkg/files"
	"net/http"
	"path/filepath"
)

type SkillsHandler struct {
	Repo files.FileHandler
}

func (h *SkillsHandler) Get(w http.ResponseWriter, r *http.Request) {
	skillsType := r.PathValue("type")
	switch skillsType {
	case "hard":
		h.getGeneric(w, append(h.Repo.BasePath, "hardskills.json")...)
	case "soft":
		h.getGeneric(w, append(h.Repo.BasePath, "softskills.json")...)
	}
}

func (h *SkillsHandler) getGeneric(w http.ResponseWriter, path ...string) {
	file, err := h.Repo.Get(filepath.Join(path...))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(file.Content))
}
