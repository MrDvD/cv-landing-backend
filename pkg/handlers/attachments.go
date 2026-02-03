package handlers

import (
	"cv-landing-backend/pkg/attachments"
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

func (h *AttachmentHandler) Add(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var rawAttachment attachments.Attachment
	err := decoder.Decode(&rawAttachment)
	if hasError(w, err) {
		return
	}
	newAttachment, err := h.Repo.Add(rawAttachment)
	if hasError(w, err) {
		return
	}
	result, err := json.Marshal(newAttachment)
	if hasError(w, err) {
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}
