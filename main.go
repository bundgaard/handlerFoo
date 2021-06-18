package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func handleError(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rErr := recover()
			if rErr != nil {
				var err error
				switch t := rErr.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				// TODO add something to notify the developer
				http.Error(w, fmt.Sprintf("received: %v", err), http.StatusInternalServerError)
			}

		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(f)
}
func main() {

	root := http.NewServeMux()

	root.Handle("/one", handleError(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Something happened on purpose")
	})))
	root.Handle("/two", handleError(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("Second purpose of error"))
	})))
	if err := http.ListenAndServe(":8080", root); err != nil {
		log.Fatal(err)
	}
}
