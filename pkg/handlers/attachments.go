package handlers

import (
	"cv-landing/pkg/attachments"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AttachmentHandler struct {
	Repo attachments.AttachmentHandler
}

func (h *AttachmentHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	foundAttachments, err := h.Repo.Get(activityId)
	if hasError(w, err) {
		return
	}
	result, err := json.Marshal(foundAttachments)
	if hasError(w, err) {
		return
	}
	w.Write(result)
}
