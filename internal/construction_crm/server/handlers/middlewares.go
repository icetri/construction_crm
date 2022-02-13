package handlers

import (
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"net/http"
)

func (h *Handlers) CheckRoleAdmin(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleAdmin {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckRoleManager(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleManager {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckRoleUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleUser {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleUser && claims.Role != types.RoleManager {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
