package utils

import (
	"net/http"
)

func ErrorHandler(err error, w http.ResponseWriter, httpStatus int) {
	http.Error(w, err.Error(), httpStatus)
	Log.Error("Error ocurred!", "error", err)
}
