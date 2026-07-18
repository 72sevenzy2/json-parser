package response

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// struct in which will hold req details
type Response struct {
	Data   json.RawMessage `json:"data"`
	Status int             `json:"status"`
	Error  string          `json:"error,omitempty"`
}

type ConfigOpts func(*Response) // any func which returns this type ONLY will use a pointer to the responses struct like used here

// status param func
func WithStatus(status int) ConfigOpts { // type of ConfigOpts which returns a func which points to the responses struct to modify its data which will then be used in my JSON helper func.
	return func(jo *Response) {
		jo.Status = status // modifies the responses struct (status) to the status passed in as the function param
	}
}

// with error param func
func WithError(msg string) ConfigOpts {
	return func(jo *Response) {
		jo.Error = msg
	}
}

// data param func
func WithData(data io.Reader) ConfigOpts { // param can be passed via strings.NewReader(..) or bytes.NewReader(..)
	return func(jo *Response) {
		b, err := io.ReadAll(data) // read data
		if err != nil {
			jo.Error = err.Error()
			return
		}

		// validate json
		if ok := json.Valid(b); !ok {
			jo.Error = "invalid json format."
			return
		}

		jo.Data = json.RawMessage(b) // typecast to type json.RawMessage for encoder.
	}
}

func JSON(w http.ResponseWriter, opts ...ConfigOpts) {
	// assigning the default values if i were to assign no params when calling the JSON func
	options := &Response{ // these options will be replaced if there were opts included when calling this func with the data in those opts
		Status: http.StatusOK,
		Data:   nil,
		Error:  "",
	}

	// initialising each opt to the appropriate param func
	for _, opt := range opts {
		opt(options) // each opt is a func that takes a pointer to the JsonOptions struct
	}

	response := &Response{
		Data:   options.Data,
		Status: options.Status,
		Error:  options.Error,
	} // initialising the response

	var bf bytes.Buffer // encoding the response with the buffer first, so if encoding fails when writing directly to http.ResponseWriter, the headers will be locked and it will be too late to change responses.

	if err := json.NewEncoder(&bf).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// writing headers if encoding succeeds (with buffer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(options.Status)

	// write the buffer to "w"
	_, err := w.Write(bf.Bytes())

	if err != nil {
		http.Error(w, "unable to write response.", http.StatusBadGateway)
		return
	}
}
