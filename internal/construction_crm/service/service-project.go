package service

import (
	"database/sql"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/construction_crm/pkg/logger"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

////MANAGER/////
func (s *Service) AddProjectManager(proj *types.AddProjectManager) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	if proj.WorkCostsOverall < 0 || proj.MaterialCostsOverall < 0 {
		return infrastruct.ErrorProjectBudgetNegative
	}

	timeCheck, err := s.timeCheckProject(proj.StartDate, proj.EndDate)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.timeCheckProject in AddProjectManager"))
		return infrastruct.ErrorProjectTimeOut
	}

	if !timeCheck {
		return infrastruct.ErrorProjectHotCoffee
	}

	project, err := s.db.ProjectAddManager(proj, tx)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ProjectAddManager in AddProjectManager"))
		return infrastruct.ErrorProject
	}

	stages := make([]types.AddStageManager, 0)

	for _, val := range proj.AddStagesManager {

		check, err := s.timeCheck(project.StartDate, project.EndDate, val.Date)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.timeCheck in AddProjectManager"))
			return infrastruct.ErrorProjectStageTimeOut
		}

		if !check {
			return infrastruct.ErrorProjectStageHotCoffee
		}

		sta, err := s.db.CheckStageDate(tx, project.Id, val.Date)
		if err != nil && err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckStageDate in AddProjectManager"))
			return infrastruct.ErrorInternalServerError
		}

		if sta != nil {
			return infrastruct.ErrorStageDate
		}

		val.ProjectId = project.Id

		stage, err := s.db.StageAddManager(&val, tx)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with StageAddManager in AddProjectManager"))
			return infrastruct.ErrorProjectStage
		}

		stages = append(stages, *stage)

		for _, valu := range val.AddCardsManager {

			if valu.SpentOnMaterials < 0 || valu.LeftToPay < 0 {
				return infrastruct.ErrorProjectBudgetNegative
			}

			//3
			checkTime, err := s.timeCheck(stage.Date, project.EndDate, valu.Deadline)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with s.timeCheck in AddProjectManager"))
				return infrastruct.ErrorProjectCardTimeOut
			}

			if !checkTime {
				return infrastruct.ErrorProjectCardHotCoffee
			}

			valu.StagesId = stage.Id

			if valu.Image == nil {
				valu.Image = pq.StringArray{}
			}

			card, err := s.db.CardAddManager(&valu, tx)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with CardAddManager in AddProjectManager"))
				return infrastruct.ErrorProjectCard
			}

			for _, value := range valu.AddTasks {

				value.CardId = card.Id

				if value.Image == nil {
					value.Image = pq.StringArray{}
				}

				_, err := s.db.TaskAddManager(&value, tx)
				if err != nil {
					logger.LogError(errors.Wrap(err, "err with TaskAddManager in AddProjectManager"))
					return infrastruct.ErrorProjectTask
				}
			}

			for _, cheques := range valu.AddCheque {

				cheques.CardId = card.Id
				cheques.UserId = project.Maker
				cheques.ProjectId = project.Id

				if cheques.File == nil {
					cheques.File = pq.StringArray{}
				}

				if cheques.Weight == nil {
					cheques.Weight = pq.Int64Array{}
				}

				if cheques.Cost < 0 {
					cheques.Cost = 0
				}

				cheque, err := s.db.ChequeAddManager(&cheques, tx)
				if err != nil {
					logger.LogError(errors.Wrap(err, "err with ChequeAddManager in AddProjectManager"))
					return infrastruct.ErrorProjectCheques
				}

				if cheque.Type == "material" {
					project.MaterialCostSpent = project.MaterialCostSpent + cheque.Cost
					if err = s.db.UpdateProjectMaterialCostSpent(tx, project.MaterialCostSpent, project.Id); err != nil {
						logger.LogError(errors.Wrap(err, "err with UpdateProjectMaterialCostSpent in AddProjectManager"))
						return infrastruct.ErrorInternalServerError
					}
				}

				if cheque.Type == "work" {
					project.WorkCostSpent = project.WorkCostSpent + cheque.Cost
					if err = s.db.UpdateProjectWorkCostSpent(tx, project.WorkCostSpent, project.Id); err != nil {
						logger.LogError(errors.Wrap(err, "err with UpdateProjectWorkCostSpent in AddProjectManager"))
						return infrastruct.ErrorInternalServerError
					}
				}
			}
		}
	}

	//project.AddStagesManager = append(project.AddStagesManager, stages...)
	//
	//for i, val := range project.AddStagesManager {
	//	cards, err := s.db.GetCardByStageIdManager(val.Id)
	//	if err != nil {
	//		logger.LogError(errors.Wrap(err, "err with GetCardByStageIdManager in AddProjectManager"))
	//		return nil, infrastruct.ErrorInternalServerError
	//	}
	//	project.AddStagesManager[i].AddCardsManager = cards
	//
	//	for in, valu := range project.AddStagesManager[i].AddCardsManager {
	//		tasks, err := s.db.GetTasksByCardIdManager(valu.Id)
	//		if err != nil {
	//			logger.LogError(errors.Wrap(err, "err with GetTasksByCardIdManager in AddProjectManager"))
	//			return nil, infrastruct.ErrorInternalServerError
	//		}
	//		project.AddStagesManager[i].AddCardsManager[in].AddTasks = tasks
	//	}
	//}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "err with Commit in AddProjectManager")
	}

	return nil
}

func (s *Service) AddCardManager(card *types.CardsAdd, userId int) (*types.CardsAdd, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	if card.SpentOnMaterials < 0 || card.LeftToPay < 0 {
		return nil, infrastruct.ErrorBadRequest
	}

	project, err := s.db.GetProject(card.ProjectId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetProject in AddCardManager"))
		return nil, infrastruct.ErrorInternalServerError
	}

	stage, err := s.db.GetStageByStageID(card.StagesId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStageByStageID in AddProjectCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	check, err := s.timeCheck(stage.Date, project.EndDate, card.Deadline)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.timeCheck in AddCardManager"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if !check {
		return nil, infrastruct.ErrorTeapot
	}

	if card.Image == nil {
		card.Image = pq.StringArray{}
	}

	cardProject, err := s.db.CardAddManagerProject(card, tx)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CardAddManager in AddCardManager"))
		return nil, infrastruct.ErrorInternalServerError
	}

	tasks := make([]types.AddTaskManager, 0)
	cheques := make([]types.AddChequeManager, 0)

	for _, val := range card.AddTasks {

		val.CardId = cardProject.Id

		if val.Image == nil {
			val.Image = pq.StringArray{}
		}

		task, err := s.db.TaskAddManager(&val, tx)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with TaskAddManager in AddCardManager"))
			return nil, infrastruct.ErrorInternalServerError
		}

		tasks = append(tasks, *task)
	}

	for _, val := range card.AddCheque {

		val.CardId = cardProject.Id
		val.UserId = userId

		if val.File == nil {
			val.File = pq.StringArray{}
		}

		if val.Weight == nil {
			val.Weight = pq.Int64Array{}
		}

		cheque, err := s.db.ChequeAddManager(&val, tx)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with ChequeAddManager in AddCardManager"))
			return nil, infrastruct.ErrorInternalServerError
		}

		proj, err := s.db.GetProject(cheque.ProjectId)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetProject in AddCardManager"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if cheque.Type == "material" {
			proj.MaterialCostSpent = proj.MaterialCostSpent + cheque.Cost
			if err = s.db.UpdateProjectMaterialCostSpent(tx, proj.MaterialCostSpent, proj.Id); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateProjectMaterialCostSpent in AddCardManager"))
				return nil, infrastruct.ErrorInternalServerError
			}
		}

		if cheque.Type == "work" {
			proj.WorkCostSpent = proj.WorkCostSpent + cheque.Cost
			if err = s.db.UpdateProjectWorkCostSpent(tx, proj.WorkCostSpent, proj.Id); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateProjectWorkCostSpent in AddCardManager"))
				return nil, infrastruct.ErrorInternalServerError
			}
		}

		cheques = append(cheques, *cheque)
	}

	cardProject.AddTasks = append(cardProject.AddTasks, tasks...)
	cardProject.AddCheque = append(cardProject.AddCheque, cheques...)

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with Commit in AddCardManager")
	}

	return cardProject, nil
}

func (s *Service) AddProject(proj *types.AddProject) (*types.Projects, error) {

	if proj.WorkCostsOverall < 0 || proj.MaterialCostsOverall < 0 {
		return nil, infrastruct.ErrorBadRequest
	}

	project, err := s.db.ProjectAdd(proj)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ProjectAdd in AddProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return project, nil
}

func (s *Service) AddProjectStage(stages []types.AddStage) ([]types.Stage, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	allStage := make([]types.Stage, 0)
	for _, stage := range stages {

		project, err := s.db.GetProject(stage.ProjectId)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetProject in AddProjectStage"))
			return nil, infrastruct.ErrorInternalServerError
		}

		check, err := s.timeCheck(project.StartDate, project.EndDate, stage.Date)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.timeCheck in AddProjectStage"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if !check {
			return nil, infrastruct.ErrorTeapot
		}

		sta, err := s.db.CheckStageDate(tx, stage.ProjectId, stage.Date)
		if err != nil && err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckStageDate in AddProjectStage"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if sta != nil {
			return nil, infrastruct.ErrorStageDate
		}

		stag, err := s.db.StageAdd(&stage, tx)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with StageAdd in AddProjectStage"))
			return nil, infrastruct.ErrorInternalServerError
		}
		allStage = append(allStage, *stag)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with Commit in AddProjectStage")
	}

	return allStage, nil
}

func (s *Service) UpdateProjectStage(stage *types.Stage) (*types.Stage, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	project, err := s.db.GetProject(stage.ProjectId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProject in UpdateProjectStage"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	check, err := s.timeCheck(project.StartDate, project.EndDate, stage.Date)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.timeCheck in UpdateProjectStage"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if !check {
		return nil, infrastruct.ErrorTeapot
	}

	sta, err := s.db.CheckStageOldDate(tx, project.Id, stage.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CheckStageOldDate in UpdateProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if sta.Date != stage.Date {
		staProj, err := s.db.CheckStageDate(tx, project.Id, stage.Date)
		if err != nil && err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckStageDate in UpdateProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if staProj != nil {
			return nil, infrastruct.ErrorStageDate
		}
	}

	stages, err := s.db.StageUpdateProject(tx, stage)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with StageUpdateProject in UpdateProjectStage"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with tx.Commit()")
	}

	return stages, nil
}

func (s *Service) UpdateProjectStageDone(stage *types.Stage) (*types.Stage, error) {

	project, err := s.db.GetProjectByStageId(stage.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetProjectByStageId in UpdateProjectStageDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	stage.Phase = true
	stages, err := s.db.StageUpdateProjectDone(stage)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with StageUpdateProjectDone in UpdateProjectStageDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	token, err := s.db.GetUserToken(project.UserId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserToken in UpdateProjectStageDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	message := &messaging.Message{
		Data: map[string]string{
			"project_id": strconv.Itoa(project.Id),
			"stage_id":   strconv.Itoa(stages.Id),
		},
		Notification: &messaging.Notification{
			Title: "Завершен этап проекта",
			Body:  fmt.Sprintf("Завершен этап проекта №%d по адресу:%s", project.Id, project.Address),
		},
		Token: token.Token,
	}

	if err = s.fb.Send(message); err != nil {
		return stages, nil
	}

	return stages, nil
}

func (s *Service) AddProjectCard(cards []types.AddCard) ([]types.Card, error) {
	tx := s.db.Begin()
	tx.Rollback()

	allCards := make([]types.Card, 0)
	for _, card := range cards {

		if card.SpentOnMaterials < 0 || card.LeftToPay < 0 {
			return nil, infrastruct.ErrorBadRequest
		}

		project, err := s.db.GetProject(card.ProjectId)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetProject in AddProjectCard"))
			return nil, infrastruct.ErrorInternalServerError
		}

		//1
		stage, err := s.db.GetStageByStageID(card.StagesId)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetStageByStageID in AddProjectCard"))
			return nil, infrastruct.ErrorInternalServerError
		}

		check, err := s.timeCheck(stage.Date, project.EndDate, card.Deadline)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.timeCheck in AddProjectCard"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if !check {
			return nil, infrastruct.ErrorTeapot
		}

		if card.Image == nil {
			card.Image = pq.StringArray{}
		}

		car, err := s.db.CardAdd(&card)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with CardAdd in AddProjectCard"))
			return nil, infrastruct.ErrorInternalServerError
		}
		allCards = append(allCards, *car)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with Commit in AddProjectCard")
	}

	return allCards, nil
}

func (s *Service) UpdateProjectCard(card *types.UpdateCard) (*types.Cards, error) {

	if card.SpentOnMaterials < 0 || card.LeftToPay < 0 {
		return nil, infrastruct.ErrorBadRequest
	}

	if card.Images == nil {
		card.Images = pq.StringArray{}
	}

	project, err := s.db.GetProject(card.ProjectId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProject in UpdateProjectCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	stage, err := s.db.GetStageByStageID(card.StagesId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetStageByStageID in AddProjectCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if stage == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	check, err := s.timeCheck(stage.Date, project.EndDate, card.Deadline)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.timeCheck in UpdateProjectCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if !check {
		return nil, infrastruct.ErrorTeapot
	}

	cardOld, err := s.db.GetCardByCardId(card.Id)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetCardByCardId in UpdateProjectCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if cardOld == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	cardProject, err := s.db.CardUpdateManagerProject(card)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CardUpdateManagerProject in UpdateProjectCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	tasks, err := s.db.GetTasksByCardId(card.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetTasksByCardId in UpdateProjectCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cheques, err := s.db.GetChequesByCardId(card.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetChequesByCardId in UpdateProjectCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cardProject.Tasks = tasks
	cardProject.Cheques = cheques

	if cardOld.Status != card.Status || cardOld.State != card.State {

		token, err := s.db.GetUserToken(project.UserId)
		if err != nil && err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserToken in UpdateProjectCard"))
			return nil, infrastruct.ErrorInternalServerError
		}

		message := &messaging.Message{
			Data: map[string]string{
				"project_id": strconv.Itoa(project.Id),
				"card_id":    strconv.Itoa(cardProject.Id),
			},
			Notification: &messaging.Notification{
				Title: "Изменение карточки этапа",
				Body:  fmt.Sprintf("Изменена карточка проекта по адресу %s,  №%d, текущая стадия:%s -%s", project.Address, project.Id, cardProject.Status, cardProject.State),
			},
			Token: token.Token,
		}

		if err = s.fb.Send(message); err != nil {
			return cardProject, nil
		}
	}

	return cardProject, nil
}

func (s *Service) UpdateProjectCardDone(card *types.Card) (*types.Card, error) {

	card.Status = "Выполнено"
	cardProject, err := s.db.CardUpdateManagerProjectDone(card)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CardUpdateManagerProjectDone in UpdateProjectCardDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project, err := s.db.GetProjectByCardId(cardProject.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetProjectByCardId in UpdateProjectCardDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	token, err := s.db.GetUserToken(project.UserId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserToken in UpdateProjectCardDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	message := &messaging.Message{
		Data: map[string]string{
			"project_id": strconv.Itoa(project.Id),
			"card_id":    strconv.Itoa(cardProject.Id),
		},
		Notification: &messaging.Notification{
			Title: "Завершена карточка этапа",
			Body:  fmt.Sprintf("Завершена карточка этапа по проекту №%d по адресу:%s", project.Id, project.Address),
		},
		Token: token.Token,
	}

	if err = s.fb.Send(message); err != nil {
		return cardProject, nil
	}

	return cardProject, nil
}

func (s *Service) AddProjectTask(tasks []types.AddTask) ([]types.Task, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	allTasks := make([]types.Task, 0)
	for _, task := range tasks {

		if task.Image == nil {
			task.Image = pq.StringArray{}
		}

		tas, err := s.db.TaskAdd(&task)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with TaskAdd in AddProjectTask"))
			return nil, infrastruct.ErrorInternalServerError
		}
		allTasks = append(allTasks, *tas)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with Commit in AddProjectTask")
	}

	return allTasks, nil
}

func (s *Service) UpdateProjectTask(tasks *types.Task) (*types.Task, error) {

	if tasks.Images == nil {
		tasks.Images = pq.StringArray{}
	}

	task, err := s.db.TaskUpdate(tasks)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with TaskUpdate in UpdateProjectTask"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return task, nil
}

func (s *Service) UpdateProjectCheque(cheque *types.Cheque) (*types.Cheque, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	cheques, err := s.db.GetCheque(cheque.Id)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetCheque in UpdateProjectCheque"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if cheques == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	project, err := s.db.GetProject(cheque.ProjectId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetProject in UpdateProjectCheque"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if cheques.ProjectId != project.Id {
		return nil, infrastruct.ErrorBadRequest
	}

	if cheque.File == nil {
		cheque.File = pq.StringArray{}
	}

	if cheque.Weight == nil {
		cheque.Weight = pq.Int64Array{}
	}

	if cheque.Cost < 0 {
		return nil, infrastruct.ErrorBadRequest
	}

	cheq, err := s.db.ChequeUpdate(tx, cheque)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ChequeUpdate in UpdateProjectCheque"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if cheq.Type != cheques.Type {
		if cheq.Type == "material" {
			project.WorkCostSpent = project.WorkCostSpent - cheques.Cost
			project.MaterialCostSpent = project.MaterialCostSpent + cheques.Cost
			if project.WorkCostSpent < 0 {
				project.WorkCostSpent = 0
			}
			if err = s.db.UpdateProjectCostSpent(tx, project.WorkCostSpent, project.MaterialCostSpent, project.Id); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateProjectCostSpent in UpdateProjectCheque material"))
				return nil, infrastruct.ErrorInternalServerError
			}
		} else {
			project.MaterialCostSpent = project.MaterialCostSpent - cheques.Cost
			project.WorkCostSpent = project.WorkCostSpent + cheques.Cost
			if project.MaterialCostSpent < 0 {
				project.MaterialCostSpent = 0
			}
			if err = s.db.UpdateProjectCostSpent(tx, project.WorkCostSpent, project.MaterialCostSpent, project.Id); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateProjectCostSpent in UpdateProjectCheque work"))
				return nil, infrastruct.ErrorInternalServerError
			}
		}
	}

	if cheq.Type == "material" {
		project.MaterialCostSpent = project.MaterialCostSpent - cheques.Cost + cheq.Cost
		if project.MaterialCostSpent < 0 {
			project.MaterialCostSpent = 0
		}
		if err = s.db.UpdateProjectMaterialCostSpent(tx, project.MaterialCostSpent, project.Id); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateProjectMaterialCostSpent in UpdateProjectCheque"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	if cheq.Type == "work" {
		project.WorkCostSpent = project.WorkCostSpent - cheques.Cost + cheq.Cost
		if project.WorkCostSpent < 0 {
			project.WorkCostSpent = 0
		}
		if err = s.db.UpdateProjectWorkCostSpent(tx, project.WorkCostSpent, project.Id); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateProjectWorkCostSpent in UpdateProjectCheque"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with Commit in UpdateProjectCheque")
	}

	return cheq, nil
}

func (s *Service) UpdateProject(proj *types.UpdateProject) (*types.ProjectManagerWithClient, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	if proj.WorkCostsOverall < 0 || proj.MaterialCostsOverall < 0 {
		return nil, infrastruct.ErrorBadRequest
	}

	timeCheck, err := s.timeCheckProject(proj.StartDate, proj.EndDate)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.timeCheckProject in UpdateProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if !timeCheck {
		return nil, infrastruct.ErrorTeapot
	}

	userClient := new(types.User)
	if proj.Client != nil {

		client, err := s.db.UpdateUserById(proj.Client, tx)
		if err != nil {
			return nil, infrastruct.ErrorPhoneOrEmailIsIncorrect
		}

		user, err := s.db.GetUserById(client.Id)
		if err != nil {
			return nil, infrastruct.ErrorInternalServerError
		}

		userClient = user
	}

	proj.UserId = userClient.ID

	project, err := s.db.ProjectUpdate(proj, tx)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ProjectUpdate in UpdateProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for _, stage := range proj.Stages {

		check, err := s.timeCheck(project.StartDate, project.EndDate, stage.Date)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.timeCheck in UpdateProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if !check {
			return nil, infrastruct.ErrorTeapot
		}

		sta, err := s.db.CheckStageOldDate(tx, project.Id, stage.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with CheckStageOldDate in UpdateProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if sta.Date != stage.Date {
			staProj, err := s.db.CheckStageDate(tx, project.Id, stage.Date)
			if err != nil && err != sql.ErrNoRows {
				logger.LogError(errors.Wrap(err, "err with CheckStageDate in UpdateProject"))
				return nil, infrastruct.ErrorInternalServerError
			}

			if staProj != nil {
				return nil, infrastruct.ErrorStageDate
			}
		}

		stag, err := s.db.StageUpdate(&stage, tx)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with StageUpdate in UpdateProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		cards, err := s.db.GetCardByStageId(stag.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCardByStageId in UpdateProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		stag.Cards = cards

		project.Stages = append(project.Stages, *stag)
	}

	cheques, err := s.db.GetChequesByProjectId(project.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetChequesByProjectId in UpdateProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project.Cheques = append(project.Cheques, cheques...)

	project.Client = userClient

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with Commit in UpdateProject")
	}

	return project, nil
}

func (s *Service) UpdateProjectDone(proj *types.Projects) (*types.ProjectManagerWithClient, error) {

	user, err := s.db.GetUserById(proj.UserId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserById in UpdateProjectDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	proj.Active = false
	project, err := s.db.ProjectDone(proj)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ProjectDone in UpdateProjectDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project.Client = user

	stages, err := s.db.GetStages(project.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStages in UpdateProjectDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project.Stages = append(project.Stages, stages...)

	for i, stage := range project.Stages {
		cards, err := s.db.GetCardByStageId(stage.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCardsByStageId in UpdateProjectDone"))
			return nil, infrastruct.ErrorInternalServerError
		}
		project.Stages[i].Cards = cards
	}

	cheques, err := s.db.GetChequesByProjectId(project.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetChequesByProjectId in UpdateProjectDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project.Cheques = append(project.Cheques, cheques...)

	token, err := s.db.GetUserToken(project.UserId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserToken in UpdateProjectCardDone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	message := &messaging.Message{
		Data: map[string]string{
			"project_id": strconv.Itoa(project.Id),
		},
		Notification: &messaging.Notification{
			Title: "Проект завершен",
			Body:  fmt.Sprintf("Завершен проект №%d по адресу:%s", project.Id, project.Address),
		},
		Token: token.Token,
	}

	if err = s.fb.Send(message); err != nil {
		return project, nil
	}

	return project, nil
}

func (s *Service) RegisterUserByManagerProject(managerId int, user *types.RegistredUserByManagerProject) (*types.User, error) {

	userInDB, err := s.db.GetUserById(user.Id)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserById in RegisterUserByManagerProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if userInDB == nil {

		user.Email = strings.ToLower(user.Email)

		if err := s.checkDoubleUserByPhone(user.Phone); err != nil {
			return nil, err
		}

		if err := s.checkDoubleUserByEmail(user.Email); err != nil {
			return nil, err
		}

		user.RegManager = true
		err = s.db.CreateUserByManager(managerId, user)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with CreateUserByManager in RegisterUserByManagerProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		userInProject, err := s.db.GetUserByPhone(user.Phone)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetUserByPhone in RegisterUserByManagerProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		return userInProject, nil
	}

	if err = s.db.AddManagerList(managerId, userInDB.ID); err != nil {
		return nil, err
	}

	return userInDB, nil
}

func (s *Service) GetManagerAllProjects(userId int) ([]types.ProjectsWithClient, error) {

	projects, err := s.db.GetManagerAllProjects(userId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetManagerAllProjects in GetManagerAllProjects"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i, val := range projects {
		user, err := s.db.GetUserById(val.UserId)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetUserById in GetManagerAllProjects"))
			return nil, infrastruct.ErrorInternalServerError
		}

		projects[i].Client = user
	}

	return projects, nil
}

func (s *Service) GetManagerProject(projectId int) (*types.ProjectManager, error) {

	project, err := s.db.GetProjectManager(projectId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProject in GetManagerProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	user, err := s.db.GetUserByIdManager(project.UserId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByIdManager in GetManagerProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	addresses, err := s.db.GetAddresses(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAddresses in GetManagerProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	user.Address = addresses

	project.Client = user

	stages, err := s.db.GetStagesWithCards(projectId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStages in GetManagerProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project.Stages = append(project.Stages, stages...)

	for i, stage := range project.Stages {
		cards, err := s.db.GetCardByStageIdWithCount(stage.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCardsByStageId in GetManagerProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		for in, card := range cards {

			count, err := s.db.GetIntCardCheques(card.Id)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with GetIntCardCheques in GetManagerProject"))
				return nil, infrastruct.ErrorInternalServerError
			}

			cards[in].Cheques = count
		}

		project.Stages[i].Cards = cards
	}

	cheques, err := s.db.GetChequesByProjectId(project.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetChequesByUserId in GetManagerProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project.Cheques = append(project.Cheques, cheques...)

	return project, nil
}

func (s *Service) GetManagerCard(cardId, projectID int) (*types.Cards, error) {

	count, err := s.db.GetIntProjects(projectID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetIntProjects in GetManagerCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if count == 0 {
		return nil, infrastruct.ErrorBadRequest
	}

	card, err := s.db.GetCardById(cardId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetCardById in GetManagerCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if card == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	tasks, err := s.db.GetTasksByCardId(card.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetTasksByCardId in GetManagerCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cheques, err := s.db.GetChequesByCardId(card.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetChequesByCardId in GetUserCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	card.Tasks = append(card.Tasks, tasks...)
	card.Cheques = append(card.Cheques, cheques...)

	return card, nil
}

////USER/////
func (s *Service) GetUserAllProjects(userId int) ([]types.Projects, error) {

	projects, err := s.db.GetUserAllProjects(userId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserAllProjects in GetUserAllProjects"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return projects, nil
}

func (s *Service) GetUserProject(projectId, userID int) (*types.Project, error) {

	project, err := s.db.GetProject(projectId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProject in GetUserProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	if project.UserId != userID {
		return nil, infrastruct.ErrorPermissionDenied
	}

	stages, err := s.db.GetStages(projectId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStages in GetUserProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project.Stages = append(project.Stages, stages...)

	for i, stage := range project.Stages {
		cards, err := s.db.GetCardByStageId(stage.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCardsByStageId in GetUserProject"))
			return nil, infrastruct.ErrorInternalServerError
		}
		project.Stages[i].Cards = cards
	}

	return project, nil
}

func (s *Service) GetUserCard(cardId, projectID, userID int) (*types.Cards, error) {

	count, err := s.db.GetIntProjects(projectID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetIntProjects in GetUserCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if count == 0 {
		return nil, infrastruct.ErrorBadRequest
	}

	project, err := s.db.GetProject(projectID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetProject in GetUserCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if project.UserId != userID {
		return nil, infrastruct.ErrorPermissionDenied
	}

	card, err := s.db.GetCardById(cardId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetCardById in GetUserCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if card == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	tasks, err := s.db.GetTasksByCardId(card.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetTasksByCardId in GetUserCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cheques, err := s.db.GetChequesByCardId(card.Id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetChequesByCardId in GetUserCard"))
		return nil, infrastruct.ErrorInternalServerError
	}

	card.Tasks = append(card.Tasks, tasks...)
	card.Cheques = append(card.Cheques, cheques...)

	return card, nil
}

func (s *Service) AddCardCheque(cheques *types.AddCheque) (*types.Cheque, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	if cheques.File == nil {
		cheques.File = pq.StringArray{}
	}

	if cheques.Weight == nil {
		cheques.Weight = pq.Int64Array{}
	}

	if cheques.Cost < 0 {
		return nil, infrastruct.ErrorBadRequest
	}

	cheque, err := s.db.AddCheque(tx, cheques)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with AddCheque in AddCardCheque"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project, err := s.db.GetProject(cheque.ProjectId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProject in AddCardCheque"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	if cheques.Type == "material" {
		project.MaterialCostSpent = project.MaterialCostSpent + cheque.Cost
		if err = s.db.UpdateProjectMaterialCostSpent(tx, project.MaterialCostSpent, project.Id); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateProjectMaterialCostSpent in AddCardCheque"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	if cheques.Type == "work" {
		project.WorkCostSpent = project.WorkCostSpent + cheque.Cost
		if err = s.db.UpdateProjectWorkCostSpent(tx, project.WorkCostSpent, project.Id); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateProjectWorkCostSpent in AddCardCheque"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with Commit in AddCardCheque")
	}

	return cheque, nil
}

func (s *Service) DeleteCardCheque(cheque *types.DeleteCheque, userID int) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	cheques, err := s.db.GetCheque(cheque.Id)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetCheque in AddCardCheque"))
		return infrastruct.ErrorInternalServerError
	}

	if cheques == nil {
		return infrastruct.ErrorBadRequest
	}

	project, err := s.db.GetProject(cheque.ProjectId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProject in AddCardCheque"))
		return infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return infrastruct.ErrorBadRequest
	}

	if project.UserId != userID {
		return infrastruct.ErrorPermissionDenied
	}

	if cheques.Type == "material" {
		project.MaterialCostSpent = project.MaterialCostSpent - cheques.Cost
		if project.MaterialCostSpent < 0 {
			project.MaterialCostSpent = 0
		}
		if err = s.db.UpdateProjectMaterialCostSpent(tx, project.MaterialCostSpent, project.Id); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateProjectMaterialCostSpent in AddCardCheque"))
			return infrastruct.ErrorInternalServerError
		}
	}

	if cheques.Type == "work" {
		project.WorkCostSpent = project.WorkCostSpent - cheques.Cost
		if project.WorkCostSpent < 0 {
			project.WorkCostSpent = 0
		}
		if err = s.db.UpdateProjectWorkCostSpent(tx, project.WorkCostSpent, project.Id); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateProjectWorkCostSpent in AddCardCheque"))
			return infrastruct.ErrorInternalServerError
		}
	}

	if err := s.db.DeleteCheque(tx, cheque); err != nil {
		logger.LogError(errors.Wrap(err, "err with DeleteCheque in DeleteCardCheque"))
		return infrastruct.ErrorInternalServerError
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "err with Commit in DeleteCardCheque")
	}

	return nil
}

func (s *Service) DeleteProjectCard(deleteCard *types.DeleteCard) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	project, err := s.db.GetProjectManager(deleteCard.ProjectID)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProjectManager in DeleteProjectCard"))
		return infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return infrastruct.ErrorProjectEmpty
	}

	if project.Maker != deleteCard.UserID {
		return infrastruct.ErrorPermissionDenied
	}

	stages, err := s.db.GetStagesWithCards(deleteCard.ProjectID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStagesWithCards in DeleteProjectCard"))
		return infrastruct.ErrorInternalServerError
	}

	var stageIsExist bool

	for _, stage := range stages {
		if stage.Id == deleteCard.StageId && stage.ProjectId == deleteCard.ProjectID {
			stageIsExist = true
		}
	}

	if stageIsExist != true {
		return infrastruct.ErrorWrongStage
	}

	cards, err := s.db.GetCardByStageId(deleteCard.StageId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCardByStageId in DeleteProjectCard"))
		return infrastruct.ErrorInternalServerError
	}

	var cardIsExist bool
	for _, card := range cards {
		if card.Id == deleteCard.CardId && card.StagesId == deleteCard.StageId {
			cardIsExist = true
		}
	}

	if cardIsExist != true {
		return infrastruct.ErrorWrongCard
	}

	tasks, err := s.db.GetTasksByCardId(deleteCard.CardId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetTasksByCardId in DeleteProjectCard"))
		return infrastruct.ErrorInternalServerError
	}

	for _, task := range tasks {
		if err = s.db.DeleteTask(tx, deleteCard.CardId, task.Id); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteTask in DeleteProjectCard"))
			return infrastruct.ErrorInternalServerError
		}
	}

	cheques, err := s.db.GetChequesByCardId(deleteCard.CardId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetChequesByCardId in DeleteProjectCard"))
		return infrastruct.ErrorInternalServerError
	}

	for _, cheque := range cheques {
		if err = s.db.DeleteCheque(tx, &types.DeleteCheque{
			Id:        cheque.Id,
			CardId:    deleteCard.CardId,
			ProjectId: deleteCard.ProjectID,
		}); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteCheque in DeleteProjectCard"))
			return infrastruct.ErrorInternalServerError
		}
	}

	if err = s.db.DeleteCard(tx, deleteCard.StageId, deleteCard.CardId); err != nil {
		logger.LogError(errors.Wrap(err, "err with DeleteCard in DeleteProjectCard"))
		return infrastruct.ErrorInternalServerError
	}

	if err = tx.Commit(); err != nil {
		logger.LogError(errors.Wrap(err, "err with tx.Commit()"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) DeleteProjectStage(deleteStage *types.DeleteStage) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	project, err := s.db.GetProjectManager(deleteStage.ProjectID)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProjectManager in DeleteProjectStage"))
		return infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return infrastruct.ErrorProjectEmpty
	}

	if project.Maker != deleteStage.UserID {
		return infrastruct.ErrorPermissionDenied
	}

	stages, err := s.db.GetStagesWithCards(deleteStage.ProjectID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStagesWithCards in DeleteProjectStage"))
		return infrastruct.ErrorInternalServerError
	}

	var stageIsExist bool
	for _, stage := range stages {
		if stage.Id == deleteStage.StageId && stage.ProjectId == deleteStage.ProjectID {
			stageIsExist = true
		}
	}

	if stageIsExist != true {
		return infrastruct.ErrorWrongStage
	}

	cards, err := s.db.GetCardByStageId(deleteStage.StageId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCardByStageId in DeleteProjectStage"))
		return infrastruct.ErrorInternalServerError
	}

	for _, card := range cards {

		tasks, err := s.db.GetTasksByCardId(card.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetTasksByCardId in DeleteProjectStage"))
			return infrastruct.ErrorInternalServerError
		}

		for _, task := range tasks {
			if err = s.db.DeleteTask(tx, card.Id, task.Id); err != nil {
				logger.LogError(errors.Wrap(err, "err with DeleteTask in DeleteProjectStage"))
				return infrastruct.ErrorInternalServerError
			}
		}

		cheques, err := s.db.GetChequesByCardId(card.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetChequesByCardId in DeleteProjectStage"))
			return infrastruct.ErrorInternalServerError
		}

		for _, cheque := range cheques {
			if err = s.db.DeleteCheque(tx, &types.DeleteCheque{
				Id:        cheque.Id,
				CardId:    card.Id,
				ProjectId: deleteStage.ProjectID,
			}); err != nil {
				logger.LogError(errors.Wrap(err, "err with DeleteCheque in DeleteProjectStage"))
				return infrastruct.ErrorInternalServerError
			}
		}

		if err = s.db.DeleteCard(tx, deleteStage.StageId, card.Id); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteCard in DeleteProjectStage"))
			return infrastruct.ErrorInternalServerError
		}
	}

	if err = s.db.DeleteStage(tx, deleteStage.StageId, deleteStage.ProjectID); err != nil {
		logger.LogError(errors.Wrap(err, "err with DeleteStage in DeleteProjectStage"))
		return infrastruct.ErrorInternalServerError
	}

	if err = tx.Commit(); err != nil {
		logger.LogError(errors.Wrap(err, "err with tx.Commit()"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}
