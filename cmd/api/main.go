//go:generate swagger generate spec
package main

import (
	"net/http"
	"runtime/debug"
	"se-challenge-payroll/internal/rest"
	"se-challenge-payroll/pkg/log"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.ZL.Fatal().Msgf("%s: %s", r, string(debug.Stack()))
		}
	}()

	server := &http.Server{
		Addr:           ":8080",
		Handler:        rest.New(newService()),
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		log.ZL.Fatal().Msgf("http server startup failed: %v", err)
	}

}
