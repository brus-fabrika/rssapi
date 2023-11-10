package main

import (
	"net/http"

	"github.com/brus-fabrika/rssapi/internal/auth"
	"github.com/brus-fabrika/rssapi/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User) error

func (apiCfg *apiConfig) middlewareAuth(next authedHandler) ApiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			return err
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			return err
		}

		return next(w, r, user)
	}
}
