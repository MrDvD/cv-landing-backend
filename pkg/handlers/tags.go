package handlers

import (
	"cv-landing-backend/pkg/activity"
	"cv-landing-backend/pkg/tags"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TagsHandler struct {
	Repo tags.TagsRepository
}

func (t *TagsHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagsFilter := tags.TagFilter{}
	rawActivityId, has := vars["id"]
	if has {
		activityId, err := strconv.Atoi(rawActivityId)
		if hasError(w, err) {
			return
		}
		tagsFilter.ActivityID = &activityId
	}
	tagType, has := vars["type"]
	if has {
		tagsFilter.TagType = &tagType
	}
	foundTags, err := t.Repo.Get(tagsFilter)
	if hasError(w, err) {
		return
	}
	result, err := json.Marshal(foundTags)
	if hasError(w, err) {
		return
	}
	w.Write(result)
}

func (h *TagsHandler) Add(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var rawTag tags.Tag
	err := decoder.Decode(&rawTag)
	if hasError(w, err) {
		return
	}
	newTag, err := h.Repo.Add(rawTag)
	if hasError(w, err) {
		return
	}
	result, err := json.Marshal(newTag)
	if hasError(w, err) {
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (h *TagsHandler) Remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawTagId, has := vars["id"]
	if !has {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tagId, err := strconv.Atoi(rawTagId)
	if hasError(w, err) {
		return
	}
	err = h.Repo.Remove(tagId)
	if hasError(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TagsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawTagId, has := vars["id"]
	if !has {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tagId, err := strconv.Atoi(rawTagId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if hasError(w, err) {
		return
	}
	var editOps []activity.EditField
	err = json.Unmarshal(bodyBytes, &editOps)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	updatedTag, err := h.Repo.Edit(tagId, editOps)
	if hasError(w, err) {
		return
	}
	serializedTag, err := json.Marshal(updatedTag)
	if hasError(w, err) {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(serializedTag)
}
