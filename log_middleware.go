package main

import (
	"log"
	"net/http"
)

func (apiCfg *apiConfig) logIt(next ApiFunc) ApiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println(r.Method, r.URL.Path)
		return next(w, r)
	}
}
