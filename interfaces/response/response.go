package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// printDebugf behaves like log.Printf only in the debug env
func printDebugf(format string, args ...interface{}) {
	if env := os.Getenv("GO_SERVER_DEBUG"); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// ErrorResponse is Error response template
type ErrorResponse struct {
	Message string `json:"reason"`
	Error   error  `json:"-"`
}

func (e *ErrorResponse) String() string {
	return fmt.Sprintf("reason: %s, error: %s", e.Message, e.Error.Error())
}

// Respond is response write to ResponseWriter
func Respond(writer http.ResponseWriter, code int, src interface{}) {
	var body []byte
	var err error

	switch s := src.(type) {
	case []byte:
		if !json.Valid(s) {
			Error(writer, http.StatusInternalServerError, err, "invalid json")
			return
		}
		body = s
	case string:
		body = []byte(s)
	case *ErrorResponse, ErrorResponse:
		// avoid infinite loop
		if body, err = json.Marshal(src); err != nil {
			writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("{\"reason\":\"failed to parse json\"}"))
			return
		}
	default:
		if body, err = json.Marshal(src); err != nil {
			Error(writer, http.StatusInternalServerError, err, "failed to parse json")
			return
		}
	}
	writer.WriteHeader(code)
	writer.Write(body)
}

// Error is wrapped Respond when error response
func Error(writer http.ResponseWriter, code int, err error, msg string) {
	e := &ErrorResponse{
		Message: msg,
		Error:   err,
	}
	printDebugf("%s", e.String())
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	Respond(writer, code, e)
}

// JSON is wrapped Respond when success response
func JSON(writer http.ResponseWriter, code int, src interface{}) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	Respond(writer, code, src)
}
