package utils

import (
	"net/http"
)

func ErrorHandler(err error, msg string, w http.ResponseWriter, httpStatus int) {
	http.Error(w, msg, httpStatus)
	Log.Error("Error ocurred!", "error", err)
}
