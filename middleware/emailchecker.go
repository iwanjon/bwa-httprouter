package middleware

import (
	"bwahttprouter/helper"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func EmailChecker(rr httprouter.Handle) httprouter.Handle {
	// var w http.ResponseWriter

	fmt.Println("test dong")

	xxx := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		defer func() {
			fmt.Println("test dong defer")
			err := recover()
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				webResponse := helper.APIResponse("success", http.StatusOK, "success", true)
				helper.WriteToResponseBody(w, webResponse)
			}
		}()
		rr(w, r, params)
	}
	return xxx
}
