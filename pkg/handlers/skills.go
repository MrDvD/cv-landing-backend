package handlers

import (
	"fmt"
	"net/http"
)

func GetSkills(w http.ResponseWriter, r *http.Request) {
	skillsType := r.PathValue("type")
	switch skillsType {
	case "hard":
		fmt.Fprintln(w, "hard! wooow!")
	case "soft":
		fmt.Fprintln(w, "soft! boooo!")
	default:
		fmt.Fprintln(w, "oof! not found this one(")
	}
}
