package handlers

import (
	"encoding/json"
	"github.com/construction_crm/internal/construction_crm/service"
	"github.com/construction_crm/internal/construction_crm/types/config"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/construction_crm/pkg/logger"

	"net/http"
)

type Handlers struct {
	s      *service.Service
	JWTKey string
}

func NewHandlers(s *service.Service, cfg *config.Config) *Handlers {
	return &Handlers{
		s:      s,
		JWTKey: cfg.JWTKey,
	}
}

func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong 1102 1216 THIS IN PROD BRANCH"))
}

func apiErrorEncode(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if customError, ok := err.(*infrastruct.CustomError); ok {
		w.WriteHeader(customError.Code)
	}

	result := struct {
		Err string `json:"error"`
	}{
		Err: err.Error(),
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogError(err)
	}
}

func apiResponseEncoder(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.LogError(err)
	}
}

func (h *Handlers) TestMail(w http.ResponseWriter, _ *http.Request) {
	err := h.s.TestMail()
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}
