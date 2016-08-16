package wot_paser

import (
	"appengine"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", stats)

	http.HandleFunc("/_ah/start", start)
	http.HandleFunc("/_ah/stop", start)
}

func stats(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	toInit(c)

	//getAccountStats(c)
	getTankStats(c)
	getTankInfo(c)
	calculateStats(account, exp_tanks)
}

func start(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}
