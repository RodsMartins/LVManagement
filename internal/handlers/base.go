package handlers

import (
	"context"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Define the interface
type HandlerInterface interface {
	HasHxHeader(r *http.Request) bool
}

// Implement the interface in the BaseHandler struct
type BaseHandler struct{
    Request *http.Request
}

func (h *BaseHandler) GetContext() context.Context {
    return h.Request.Context()
}

// Function to check if the request has the hx- header
func (h BaseHandler) UsesHtmx(r *http.Request) bool {
	hxHeader := r.Header.Get("hx-request")
	return hxHeader == "true"
}

func (h BaseHandler) GetUrlUuid(urlUUID string, r *http.Request) (string, error) {
	// Get the URL parameter "id" using chi
	id := chi.URLParam(r, urlUUID)

	// Validate the ID using uuid
	err := uuid.Validate(id)
	if  err != nil {
		return "", err
	}

	return id, nil
}

func (h BaseHandler) GetUrlUuidOrEmpty(urlUUID string, r *http.Request) (string, error) {
	// Get the URL parameter "id" using chi
	id := chi.URLParam(r, urlUUID)

	if id == "" {
		return "", nil
		
	}

	// Validate the ID using uuid
	err := uuid.Validate(id)
	if  err != nil {
		return "", err
	}

	return id, nil
}

/*
func (r *BaseHandler) renderComponent(w http.ResponseWriter, componentName string, data ...interface{}) {
    component, err := r.getComponent(componentName, data...)
    if err != nil {
        http.Error(w, "Failed to render component", http.StatusInternalServerError)
        return
    }

    err = component.Render(r.GetContext(), w)
    if err != nil {
        http.Error(w, "Failed to render component", http.StatusInternalServerError)
    }
}

func (h *BaseHandler) renderPage(w http.ResponseWriter, pageName string, data ...interface{}) {
    page, err := h.getPage(pageName)
    if err != nil {
        http.Error(w, "Failed to render page", http.StatusInternalServerError)
        return
    }

    err = page.Render(h.GetContext(), w)
    if err != nil {
        http.Error(w, "Failed to render page", http.StatusInternalServerError)
    }
}

func (h *BaseHandler) getComponent(name string) (templ.Component, error) {
    switch name {
    case "components.SeedList":
        return components.SeedList, nil
    // Add other component cases as needed
    default:
        return nil, fmt.Errorf("unknown component: %s", name)
    }
}

func (h *BaseHandler) getPage(name string, data ...interface{}) (templ.Component, error) {
    switch name {
    case "pages.SeedListPage":
        return pages.SeedsPage, nil
    // Add other page cases as needed
    default:
        return nil, fmt.Errorf("unknown page: %s", name)
    }
}*/