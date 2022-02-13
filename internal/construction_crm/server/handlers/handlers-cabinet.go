package handlers

import (
	"encoding/json"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"

	"net/http"
)

func (h *Handlers) Cabinet(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	cabinet, err := h.s.GetCabinet(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cabinet)
}

func (h *Handlers) PutCabinetInfo(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	userInfo := types.PutUserInfo{ID: claims.UserID}

	if err = json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cabinet, err := h.s.PutCabinetUserInfo(&userInfo)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cabinet)
}

func (h *Handlers) PutCabinetValidatContact(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	userConct := types.PutUserCont{ID: claims.UserID}

	if err = json.NewDecoder(r.Body).Decode(&userConct); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cabinet, err := h.s.PutUserValidatePhone(&userConct)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cabinet)
}

func (h *Handlers) PutCabinetContact(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	user := types.PutUserPhone{ID: claims.UserID}

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.s.PutUserPhone(&user)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}
