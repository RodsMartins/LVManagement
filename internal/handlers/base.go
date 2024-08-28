package handlers

import (
	"net/http"
)

// Define the interface
type HandlerInterface interface {
	HasHxHeader(r *http.Request) bool
}

// Implement the interface in the BaseHandler struct
type BaseHandler struct{
}

// Function to check if the request has the hx- header
func (h BaseHandler) UsesHtmx(r *http.Request) bool {
	hxHeader := r.Header.Get("hx-request")
	return hxHeader == "true"
}
