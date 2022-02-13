package server

import (
	"github.com/construction_crm/internal/construction_crm/server/handlers"
	"github.com/construction_crm/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
)

func StartServer(handlers *handlers.Handlers, port string) {

	router := NewRouter(handlers)
	logger.LogInfo("Restart service")
	if err := http.ListenAndServe(port, router); err != nil {
		logger.LogFatal(errors.Wrap(err, "err with NewRouter"))
	}
}
