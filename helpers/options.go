package helpers

import (
	"github.com/72sevenzy2/json-parser/response"
	"net/http"
)

// Ok helper. -> pass in a http.ResponseWriter, and some data. Example: map[string]string{ "User": "Name" }
func Ok(w http.ResponseWriter, data any) {
	response.JSON(w, response.WithData(data), response.WithStatus(http.StatusOK))
}

// Failed helper. -> pass in a http.ResponseWriter, a status code (int), and a message to display to the user (string).
func Failed(w http.ResponseWriter, status int, msg string) {
	response.JSON(w, response.WithStatus(status), response.WithError(msg))
}