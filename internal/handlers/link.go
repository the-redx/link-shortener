package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/the-redx/link-shortener/internal/domain"
	"github.com/the-redx/link-shortener/internal/services"
	"github.com/the-redx/link-shortener/pkg/errs"
)

type LinkHandler struct {
	service services.LinkService
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func (ch *LinkHandler) RedirectToLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId := vars["link_id"]

	link, appErr := ch.service.GetLinkByID("", linkId)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	http.Redirect(w, r, link.Url, http.StatusPermanentRedirect)
}

func (ch *LinkHandler) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	userId, appErr := getUserIdFromContext(r)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	links, appErr := ch.service.GetAllLinks(userId)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, links)
}

func (ch *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId := vars["link_id"]

	userId, appErr := getUserIdFromContext(r)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	link, appErr := ch.service.GetLinkByID(userId, linkId)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, link)
}

func (ch *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var link domain.CreateLinkDTO

	userId, appErr := getUserIdFromContext(r)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		writeError(w, errs.NewBadRequestError("Invalid link data"))
		return
	}

	err = validate.Struct(link)
	if err != nil {
		writeError(w, errs.NewBadRequestError(err.Error()))
		return
	}

	newLink, appErr := ch.service.CreateLink(userId, &link)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, newLink)
}

func (ch *LinkHandler) UpdateLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId := vars["link_id"]

	var link domain.UpdateLinkDTO

	userId, appErr := getUserIdFromContext(r)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		writeError(w, errs.NewBadRequestError("Invalid body"))
		return
	}

	err = validate.Struct(link)
	if err != nil {
		writeError(w, errs.NewBadRequestError(err.Error()))
		return
	}

	newLink, appErr := ch.service.UpdateLinkByID(userId, linkId, &link)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, newLink)
}

func (ch *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId := vars["link_id"]

	userId, appErr := getUserIdFromContext(r)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	link, appErr := ch.service.DeleteLinkByID(userId, linkId)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, link)
}

func NewLinkHandler(service services.LinkService) *LinkHandler {
	return &LinkHandler{service}
}
