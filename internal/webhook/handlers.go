package handlers

import "net/http"

type Validate interface {
	Robokassa()
	// Yookassa()
	// Some shit()
}

func RobokassaHandler() (pattern string, handler http.HandlerFunc) {
	return "/api/robokassa", func(w http.ResponseWriter, r *http.Request) {

		// write code here

	}
}

// func Yookassa() http.HandlerFunc{}
