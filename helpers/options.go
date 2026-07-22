package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/72sevenzy2/json-parser/response"
)

/* Ok helper. -> pass in a http.ResponseWriter, and some data with type []byte. 

Example:
b, err := json.Marshal(v)
if err != nil {
...
}

Ok(http.ResponseWriter, b)
*/
func Ok(w http.ResponseWriter, data []byte) { 
	// check if data is valid json
	if ok := json.Valid(data); !ok {
		http.Error(w, "invalid json.", http.StatusNotAcceptable)
		return
	}

	response.JSON(w, response.WithData(bytes.NewReader(data)), response.WithStatus(http.StatusOK))
}

// Failed helper. -> pass in a http.ResponseWriter, a status code (int), and a message to display to the user (string).
func Failed(w http.ResponseWriter) {
	response.JSON(w, response.WithStatus(http.StatusBadRequest), response.WithError(http.StatusText(http.StatusBadRequest))
}
