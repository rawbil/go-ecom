package utils

import (
	"fmt"
	"net/http"
)

func ErrorHandler(err error, w http.ResponseWriter, httpStatus int) {
	http.Error(w, err.Error(), httpStatus)
	fmt.Println(err)
}
