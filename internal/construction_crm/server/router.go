package server

import (
	"github.com/construction_crm/internal/construction_crm/server/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(h *handlers.Handlers) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	router.Methods(http.MethodGet).Path("/ping").HandlerFunc(h.Ping)
	router.Methods(http.MethodPost).Path("/test").HandlerFunc(h.TestMail)

	auth := router.PathPrefix("/auth").Subrouter()
	auth.Methods(http.MethodPost).Path("/sing-up").HandlerFunc(h.SingUp)
	auth.Methods(http.MethodPost).Path("/sing-in/phone").HandlerFunc(h.SingInPhone)
	auth.Methods(http.MethodPost).Path("/sing-in/email").HandlerFunc(h.SingInEmail)
	auth.Methods(http.MethodPost).Path("/sing-in/phone/confirm").HandlerFunc(h.SingInCodePhone)
	auth.Methods(http.MethodPost).Path("/sing-in/email/confirm").HandlerFunc(h.SingInCodeEmail)

	upload := router.PathPrefix("/upload").Subrouter()
	upload.Use(h.CheckUser)
	upload.Methods(http.MethodPost).Path("").HandlerFunc(h.Upload)
	upload.Methods(http.MethodPost).Path("/info").HandlerFunc(h.DeviceInfo)

	cabinet := router.PathPrefix("/cabinet").Subrouter()
	cabinet.Use(h.CheckRoleUser)

	cabinet.Methods(http.MethodGet).Path("").HandlerFunc(h.Cabinet)
	cabinet.Methods(http.MethodPut).Path("").HandlerFunc(h.PutCabinetInfo)
	cabinet.Methods(http.MethodPut).Path("/user/phone").HandlerFunc(h.PutCabinetContact)
	cabinet.Methods(http.MethodPost).Path("/user/phone/confirm").HandlerFunc(h.PutCabinetValidatContact)

	project := router.PathPrefix("").Subrouter()
	project.Use(h.CheckRoleUser)

	project.Methods(http.MethodGet).Path("/projects").HandlerFunc(h.GetUserAllProjects)
	project.Methods(http.MethodGet).Path("/project/{id:[0-9]+}").HandlerFunc(h.GetProject)
	project.Methods(http.MethodGet).Path("/project/{id:[0-9]+}/calendar").HandlerFunc(h.GetCalendar)
	project.Methods(http.MethodGet).Path("/project/{idProject:[0-9]+}/stage/card/{id:[0-9]+}").HandlerFunc(h.GetCard)
	project.Methods(http.MethodPost).Path("/project/{idProject:[0-9]+}/stage/card/{id:[0-9]+}/cheque").HandlerFunc(h.AddCheque)
	project.Methods(http.MethodDelete).Path("/project/{idProject:[0-9]+}/stage/card/{id:[0-9]+}/cheque").HandlerFunc(h.DeleteCheque)

	manager := router.PathPrefix("/manager").Subrouter()
	manager.Use(h.CheckRoleManager)

	router.Methods(http.MethodPost).Path("/auth/manager/sing-in").HandlerFunc(h.SingInManager)

	manager.Methods(http.MethodPost).Path("/register/user").HandlerFunc(h.RegisterUserManager)
	manager.Methods(http.MethodGet).Path("/users").HandlerFunc(h.ListUsers)
	manager.Methods(http.MethodGet).Path("/user/{id:[0-9]+}").HandlerFunc(h.GetUser)

	manager.Methods(http.MethodGet).Path("/projects").HandlerFunc(h.GetManagerAllProjects)
	manager.Methods(http.MethodGet).Path("/projects/calendar").HandlerFunc(h.GetManagerCalendar)
	manager.Methods(http.MethodGet).Path("/projects/{id:[0-9]+}").HandlerFunc(h.GetManagerProject)
	manager.Methods(http.MethodGet).Path("/project/{idProject:[0-9]+}/stage/card/{id:[0-9]+}").HandlerFunc(h.GetManagerCard)

	manager.Methods(http.MethodPost).Path("/project/register/user").HandlerFunc(h.RegistredUserByManagerProject)

	manager.Methods(http.MethodPost).Path("/project").HandlerFunc(h.AddProject)
	manager.Methods(http.MethodPut).Path("/project").HandlerFunc(h.UpdateProject)
	manager.Methods(http.MethodPost).Path("/project/done").HandlerFunc(h.UpdateProjectDone)

	manager.Methods(http.MethodPost).Path("/project/stages").HandlerFunc(h.AddProjectStage)
	manager.Methods(http.MethodPut).Path("/project/stages").HandlerFunc(h.UpdateProjectStage)
	manager.Methods(http.MethodDelete).Path("/project/{idProject:[0-9]+}/stage/{idStage:[0-9]+}").HandlerFunc(h.DeleteProjectStage)
	manager.Methods(http.MethodPost).Path("/project/stages/done").HandlerFunc(h.UpdateProjectStageDone)

	manager.Methods(http.MethodPost).Path("/project/stage/cards").HandlerFunc(h.AddProjectCard)
	manager.Methods(http.MethodPut).Path("/project/stage/cards").HandlerFunc(h.UpdateProjectCard)
	manager.Methods(http.MethodDelete).Path("/project/{idProject:[0-9]+}/stage/{idStage:[0-9]+}/card/{id:[0-9]+}").HandlerFunc(h.DeleteProjectCard)
	manager.Methods(http.MethodPost).Path("/project/stage/cards/done").HandlerFunc(h.UpdateProjectCardDone)

	manager.Methods(http.MethodPost).Path("/project/stage/cards/tasks").HandlerFunc(h.AddProjectTask)
	manager.Methods(http.MethodPut).Path("/project/stage/cards/tasks").HandlerFunc(h.UpdateProjectTask)

	manager.Methods(http.MethodPost).Path("/project/{idProject:[0-9]+}/stage/card/{id:[0-9]+}/cheque").HandlerFunc(h.AddCheque)
	manager.Methods(http.MethodDelete).Path("/project/{idProject:[0-9]+}/stage/card/{id:[0-9]+}/cheque").HandlerFunc(h.DeleteCheque)
	manager.Methods(http.MethodPut).Path("/project/stage/cards/tasks/cheque").HandlerFunc(h.UpdateProjectCheque)

	manager.Methods(http.MethodPost).Path("/project/add").HandlerFunc(h.AddProjectManager)
	manager.Methods(http.MethodPost).Path("/project/stage/card/add").HandlerFunc(h.AddCardManager)

	return router
}
