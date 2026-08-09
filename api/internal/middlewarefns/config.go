package middlewarefns

import "net/http"

type Mdlw func(http.Handler) http.Handler