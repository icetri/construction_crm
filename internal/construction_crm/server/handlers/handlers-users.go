package handlers

import (
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handlers) ListUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.s.UsersList()
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, users)
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {

	userId := mux.Vars(r)["id"]

	userID, err := strconv.Atoi(userId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	user, err := h.s.UserById(userID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, user)
}
