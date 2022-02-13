package handlers

import (
	"encoding/json"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

/////MANAGER/////
func (h *Handlers) AddProjectManager(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	project := types.AddProjectManager{Maker: claims.UserID}

	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.s.AddProjectManager(&project)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) AddCardManager(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	card := types.CardsAdd{}

	err = json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cardProject, err := h.s.AddCardManager(&card, claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cardProject)
}

func (h *Handlers) AddProject(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	project := types.AddProject{Maker: claims.UserID}

	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	proj, err := h.s.AddProject(&project)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, proj)
}

func (h *Handlers) UpdateProject(w http.ResponseWriter, r *http.Request) {

	project := types.UpdateProject{}

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	proj, err := h.s.UpdateProject(&project)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, proj)
}

func (h *Handlers) UpdateProjectDone(w http.ResponseWriter, r *http.Request) {

	project := types.Projects{}

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	proj, err := h.s.UpdateProjectDone(&project)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, proj)
}

func (h *Handlers) AddProjectStage(w http.ResponseWriter, r *http.Request) {

	addStage := []types.AddStage{}

	err := json.NewDecoder(r.Body).Decode(&addStage)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	stages, err := h.s.AddProjectStage(addStage)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, stages)
}

func (h *Handlers) UpdateProjectStage(w http.ResponseWriter, r *http.Request) {

	stage := types.Stage{}

	err := json.NewDecoder(r.Body).Decode(&stage)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	stages, err := h.s.UpdateProjectStage(&stage)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, stages)
}

func (h *Handlers) UpdateProjectStageDone(w http.ResponseWriter, r *http.Request) {

	stage := types.Stage{}

	err := json.NewDecoder(r.Body).Decode(&stage)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	stages, err := h.s.UpdateProjectStageDone(&stage)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, stages)
}

func (h *Handlers) AddProjectCard(w http.ResponseWriter, r *http.Request) {

	addCards := []types.AddCard{}

	err := json.NewDecoder(r.Body).Decode(&addCards)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cards, err := h.s.AddProjectCard(addCards)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cards)
}

func (h *Handlers) UpdateProjectCard(w http.ResponseWriter, r *http.Request) {

	card := types.UpdateCard{}

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cards, err := h.s.UpdateProjectCard(&card)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cards)
}

func (h *Handlers) UpdateProjectCardDone(w http.ResponseWriter, r *http.Request) {

	card := types.Card{}

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cards, err := h.s.UpdateProjectCardDone(&card)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cards)
}

func (h *Handlers) AddProjectTask(w http.ResponseWriter, r *http.Request) {

	addTask := []types.AddTask{}

	err := json.NewDecoder(r.Body).Decode(&addTask)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	task, err := h.s.AddProjectTask(addTask)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, task)
}

func (h *Handlers) UpdateProjectTask(w http.ResponseWriter, r *http.Request) {

	task := types.Task{}

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	tas, err := h.s.UpdateProjectTask(&task)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, tas)
}

func (h *Handlers) UpdateProjectCheque(w http.ResponseWriter, r *http.Request) {

	cheque := types.Cheque{}

	err := json.NewDecoder(r.Body).Decode(&cheque)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cheq, err := h.s.UpdateProjectCheque(&cheque)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cheq)
}

func (h *Handlers) RegistredUserByManagerProject(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	managerRegisterUser := types.RegistredUserByManagerProject{}

	err = json.NewDecoder(r.Body).Decode(&managerRegisterUser)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	userInProject, err := h.s.RegisterUserByManagerProject(claims.UserID, &managerRegisterUser)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, userInProject)
}

func (h *Handlers) GetManagerAllProjects(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	projects, err := h.s.GetManagerAllProjects(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, projects)
}

func (h *Handlers) GetManagerProject(w http.ResponseWriter, r *http.Request) {

	projID := mux.Vars(r)["id"]

	projectID, err := strconv.Atoi(projID)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	project, err := h.s.GetManagerProject(projectID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, project)
}

func (h *Handlers) GetManagerCard(w http.ResponseWriter, r *http.Request) {

	cardId := mux.Vars(r)["id"]

	cardID, err := strconv.Atoi(cardId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	projectId := mux.Vars(r)["idProject"]

	projectID, err := strconv.Atoi(projectId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	card, err := h.s.GetManagerCard(cardID, projectID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, card)
}

//////USER////////
func (h *Handlers) GetUserAllProjects(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	projects, err := h.s.GetUserAllProjects(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, projects)
}

func (h *Handlers) GetProject(w http.ResponseWriter, r *http.Request) {

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

	project, err := h.s.GetUserProject(projectID, claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, project)
}

func (h *Handlers) GetCard(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	cardId := mux.Vars(r)["id"]

	cardID, err := strconv.Atoi(cardId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	projectId := mux.Vars(r)["idProject"]

	projectID, err := strconv.Atoi(projectId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	card, err := h.s.GetUserCard(cardID, projectID, claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, card)
}

func (h *Handlers) AddCheque(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	cardId := mux.Vars(r)["id"]

	cardID, err := strconv.Atoi(cardId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	projectId := mux.Vars(r)["idProject"]

	projectID, err := strconv.Atoi(projectId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	addCheque := types.AddCheque{CardId: cardID, UserId: claims.UserID, ProjectId: projectID}

	if err = json.NewDecoder(r.Body).Decode(&addCheque); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cheque, err := h.s.AddCardCheque(&addCheque)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, cheque)
}

func (h *Handlers) DeleteCheque(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	cardId := mux.Vars(r)["id"]

	cardID, err := strconv.Atoi(cardId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	projectId := mux.Vars(r)["idProject"]

	projectID, err := strconv.Atoi(projectId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	cheque := types.DeleteCheque{CardId: cardID, ProjectId: projectID}

	if err = json.NewDecoder(r.Body).Decode(&cheque); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.s.DeleteCardCheque(&cheque, claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) DeleteProjectCard(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	cardId := mux.Vars(r)["id"]

	cardID, err := strconv.Atoi(cardId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	stageId := mux.Vars(r)["idStage"]

	stageID, err := strconv.Atoi(stageId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	projectId := mux.Vars(r)["idProject"]

	projectID, err := strconv.Atoi(projectId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	deleteCard := &types.DeleteCard{
		ProjectID: projectID,
		StageId:   stageID,
		CardId:    cardID,
		UserID:    claims.UserID,
	}

	err = h.s.DeleteProjectCard(deleteCard)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) DeleteProjectStage(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.JWTKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	stageId := mux.Vars(r)["idStage"]

	stageID, err := strconv.Atoi(stageId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	projectId := mux.Vars(r)["idProject"]

	projectID, err := strconv.Atoi(projectId)
	if err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	deleteStage := &types.DeleteStage{
		ProjectID: projectID,
		StageId:   stageID,
		UserID:    claims.UserID,
	}

	err = h.s.DeleteProjectStage(deleteStage)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}
