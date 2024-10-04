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

	link, err := ch.service.GetLinkByID(linkId)

	if err != nil {
		writeError(w, err)
		return
	}

	if link.Status != domain.Active {
		err := errs.NewNotFoundError("Link is not found")
		writeError(w, err)
		return
	}

	http.Redirect(w, r, link.Url, http.StatusPermanentRedirect)
}

func (ch *LinkHandler) GetAllLinks(w http.ResponseWriter, r *http.Request) {
	links, err := ch.service.GetAllLinks()
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, links)
}

func (ch *LinkHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId := vars["link_id"]

	link, err := ch.service.GetLinkByID(linkId)
	if err != nil {
		writeError(w, err)
		return
	}

	writeResponse(w, http.StatusOK, link)
}

func (ch *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var link domain.CreateLinkDTO

	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		appErr := errs.NewBadRequestError("Invalid link data")
		writeError(w, appErr)
		return
	}

	err = validate.Struct(link)
	if err != nil {
		appErr := errs.NewBadRequestError(err.Error())
		writeError(w, appErr)
		return
	}

	newLink, appErr := ch.service.CreateLink(&link)
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

	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		appErr := errs.NewBadRequestError("Invalid body")
		writeError(w, appErr)
		return
	}

	err = validate.Struct(link)
	if err != nil {
		appErr := errs.NewBadRequestError(err.Error())
		writeError(w, appErr)
		return
	}

	newLink, appErr := ch.service.UpdateLinkByID(linkId, &link)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, newLink)
}

func (ch *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	linkId := vars["link_id"]

	link, appErr := ch.service.DeleteLinkByID(linkId)
	if appErr != nil {
		writeError(w, appErr)
		return
	}

	writeResponse(w, http.StatusOK, link)
}

func NewLinkHandler(service services.LinkService) *LinkHandler {
	return &LinkHandler{service}
}
