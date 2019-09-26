package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func Respond(w http.ResponseWriter, status int, payload interface{}) {
	b, err := json.Marshal(payload)
	if err != nil {
		RespondError(w, err.Error())
		return
	}
	setHeaders(w, status)
	fmt.Fprintf(w, string(b))
}

func RespondMessage(w http.ResponseWriter, status int, message string) {
	setHeaders(w, status)
	fmt.Fprintf(w, `{"code": `+strconv.Itoa(status)+`, "message": "`+message+`"}`)
}

func RespondError(w http.ResponseWriter, message string) {
	setHeaders(w, http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error": "`+message+`"}`)
}

func setHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}
