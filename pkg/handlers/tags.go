package handlers

import (
	"cv-landing/pkg/tags"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TagsHandler struct {
	Repo tags.TagsHandler
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
