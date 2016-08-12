package base

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", start)

	http.HandleFunc("/_ah/start", start)
	http.HandleFunc("/_ah/stop", start)
}

var (
        application_id = "76d0d8bd88a4c872b17404c408332284"
        access_token = ""
)

func start(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}
