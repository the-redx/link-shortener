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

	link, appErr := ch.service.GetLinkByID(linkId, r.Context())
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	http.Redirect(w, r, link.Url, http.StatusPermanentRedirect)
}

func (ch *LinkHandler) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	links, appErr := ch.service.GetAllLinks(r.Context())
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, links)
}

func (ch *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId := vars["link_id"]

	link, appErr := ch.service.GetLinkByID(linkId, r.Context())
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, link)
}

func (ch *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var link domain.CreateLinkDTO

	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		writeError(w, errs.NewBadRequestError("Invalid link data"))
		return
	}

	if err := validate.Struct(link); err != nil {
		writeError(w, errs.NewBadRequestError(err.Error()))
		return
	}

	newLink, appErr := ch.service.CreateLink(&link, r.Context())
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

	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		writeError(w, errs.NewBadRequestError("Invalid body"))
		return
	}

	if err := validate.Struct(link); err != nil {
		writeError(w, errs.NewBadRequestError(err.Error()))
		return
	}

	newLink, appErr := ch.service.UpdateLinkByID(linkId, &link, r.Context())
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, newLink)
}

func (ch *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId := vars["link_id"]

	link, appErr := ch.service.DeleteLinkByID(linkId, r.Context())
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, link)
}

func NewLinkHandler(service services.LinkService) *LinkHandler {
	return &LinkHandler{service}
}
