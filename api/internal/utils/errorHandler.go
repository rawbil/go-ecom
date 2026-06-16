package utils

import (
	"fmt"
	"net/http"
)

func ErrorHandler(err error, w http.ResponseWriter) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	fmt.Println(err)
}
