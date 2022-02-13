package handlers

import (
	"encoding/json"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"net/http"
)

func (h *Handlers) SingUp(w http.ResponseWriter, r *http.Request) {

	user := types.Register{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	userMan, err := h.s.CheckRegister(user.Phone)
	if err != nil {
		apiResponseEncoder(w, userMan)
		return
	}

	err = h.s.RegisterNewUser(&user)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

}

func (h *Handlers) SingInPhone(w http.ResponseWriter, r *http.Request) {

	user := types.AuthCodePhone{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	token, err := h.s.AuthUserPhone(&user)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, token)
}

func (h *Handlers) SingInEmail(w http.ResponseWriter, r *http.Request) {

	user := types.AuthCodeEmail{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	token, err := h.s.AuthUserEmail(&user)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, token)
}

func (h *Handlers) RegisterUserManager(w http.ResponseWriter, r *http.Request) {

	managerRegisterUser := types.Register{}

	err := json.NewDecoder(r.Body).Decode(&managerRegisterUser)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	if err = h.s.RegisterUserManager(&managerRegisterUser); err != nil {
		apiErrorEncode(w, err)
		return
	}

}

func (h *Handlers) SingInCodePhone(w http.ResponseWriter, r *http.Request) {

	user := types.AuthPhone{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.s.CodePhone(&user)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

}

func (h *Handlers) SingInCodeEmail(w http.ResponseWriter, r *http.Request) {

	user := types.AuthEmail{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.s.CodeEmail(&user)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

}

func (h *Handlers) SingInManager(w http.ResponseWriter, r *http.Request) {

	user := types.AuthEmailManager{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	token, err := h.s.AuthManagerEmail(&user)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, token)
}
