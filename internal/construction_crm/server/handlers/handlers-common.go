package handlers

import (
	"encoding/json"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (h *Handlers) Upload(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	bucketName := r.FormValue("bucket")

	filesBucket := make([]types.FileInfo, 0)
	files := r.MultipartForm.File["file"]
	for _, file := range files {
		path, err := h.s.UploadFile(file, bucketName, claims)
		if err != nil {
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
			return
		}
		filesBucket = append(filesBucket, *path)
	}
	apiResponseEncoder(w, filesBucket)
}

func (h *Handlers) GetCalendar(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	projID := mux.Vars(r)["id"]

	projectID, err := strconv.Atoi(projID)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	calendar, err := h.s.GetUserCalendar(projectID, claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, calendar)
}

func (h *Handlers) GetManagerCalendar(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	calendar, err := h.s.GetManagerCalendar(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, calendar)
}

func (h *Handlers) DeviceInfo(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	token := types.DeviceToken{}

	err = json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.s.UpdateDeviceInfo(&token, claims)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}
