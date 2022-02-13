package postgres

import (
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func (p *Postgres) ProjectAdd(project *types.AddProject) (*types.Projects, error) {
	tx := p.db.MustBegin()
	defer tx.Rollback()

	var proj types.Projects
	if err := tx.QueryRowx(`INSERT INTO project (image, file, address, user_id, start_date, end_date, active, maker_id, material_costs_over_all, work_costs_over_all, material_cost_spent, work_cost_spent) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning *`, project.Image, project.File, project.Address, project.UserId, project.StartDate, project.EndDate, project.Active, project.Maker, project.MaterialCostsOverall, project.WorkCostsOverall, project.MaterialCostSpent, project.WorkCostSpent).StructScan(&proj); err != nil {
		return nil, err
	}

	_ = tx.MustExec(`INSERT INTO addresses (address, user_id) 
	VALUES ($1, $2)`, project.Address, project.UserId)

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "err with Commit in ProjectAdd")
	}

	return &proj, nil
}

func (p *Postgres) StageAdd(stage *types.AddStage, tx *sqlx.Tx) (*types.Stage, error) {

	var stag types.Stage
	if err := tx.QueryRowx(`INSERT INTO stages (name, project_id, phase, date) 
	VALUES ($1, $2, $3, $4) returning *`, stage.Name, stage.ProjectId, stage.Phase, stage.Date).StructScan(&stag); err != nil {
		return nil, err
	}

	return &stag, nil
}

func (p *Postgres) StageUpdateProject(tx *sqlx.Tx, stage *types.Stage) (*types.Stage, error) {

	stag := new(types.Stage)
	if err := tx.QueryRowx(`UPDATE stages set name = $1, project_id = $2, phase = $3, date = $4
	where id = $5 returning *`, stage.Name, stage.ProjectId, stage.Phase, stage.Date, stage.Id).StructScan(stag); err != nil {
		return nil, err
	}

	return stag, nil
}

func (p *Postgres) StageUpdateProjectDone(stage *types.Stage) (*types.Stage, error) {

	stag := new(types.Stage)
	if err := p.db.QueryRowx(`UPDATE stages set phase = $1
	where id = $2 returning *`, stage.Phase, stage.Id).StructScan(stag); err != nil {
		return nil, err
	}

	return stag, nil
}

func (p *Postgres) StageUpdate(stage *types.UpdateStage, tx *sqlx.Tx) (*types.Stages, error) {

	stag := new(types.Stages)
	if err := tx.QueryRowx(`UPDATE stages set name = $1, phase = $2, date = $3
	where id = $4 returning *`, stage.Name, stage.Phase, stage.Date, stage.Id).StructScan(stag); err != nil {
		return nil, err
	}

	return stag, nil
}

func (p *Postgres) CardAdd(card *types.AddCard) (*types.Card, error) {

	var car types.Card
	if err := p.db.QueryRowx(`INSERT INTO cards (images, title, deadline, stages_id, rating, description, left_to_pay, spent_on_materials, status, state) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $89, $10) returning *`, card.Image, card.Title, card.Deadline, card.StagesId, card.Rating, card.Description, card.LeftToPay, card.SpentOnMaterials, card.Status, card.State).StructScan(&car); err != nil {
		return nil, err
	}

	return &car, nil
}

func (p *Postgres) TaskAdd(task *types.AddTask) (*types.Task, error) {

	var tas types.Task
	if err := p.db.QueryRowx(`INSERT INTO tasks (images, title, card_id, complete, length) 
	VALUES ($1, $2, $3, $4, $5) returning *`, task.Image, task.Title, task.CardId, task.Complete, task.Length).StructScan(&tas); err != nil {
		return nil, err
	}

	return &tas, nil
}

func (p *Postgres) TaskUpdate(task *types.Task) (*types.Task, error) {

	tas := new(types.Task)
	if err := p.db.QueryRowx(`UPDATE tasks SET images = $1, title = $2, card_id = $3, complete = $4, length = $5
	WHERE id = $6 returning *`, task.Images, task.Title, task.CardId, task.Complete, task.Length, task.Id).StructScan(tas); err != nil {
		return nil, err
	}

	return tas, nil
}

func (p *Postgres) ChequeUpdate(tx *sqlx.Tx, cheque *types.Cheque) (*types.Cheque, error) {

	cheq := new(types.Cheque)
	if err := tx.QueryRowx(`UPDATE cheques SET file = $1, name = $2, cost = $3, type = $4, length = $5, weight = $6
	WHERE id = $7 returning *`, cheque.File, cheque.Name, cheque.Cost, cheque.Type, cheque.Length, cheque.Weight, cheque.Id).StructScan(cheq); err != nil {
		return nil, err
	}

	return cheq, nil
}

func (p *Postgres) ProjectUpdate(project *types.UpdateProject, tx *sqlx.Tx) (*types.ProjectManagerWithClient, error) {

	proj := new(types.ProjectManagerWithClient)

	err := tx.QueryRowx(`update project set image = $1,
                  file = $2,
                  address = $3,
                  user_id = $4,
                  start_date = $5,
                  end_date = $6,
                  active = $7,
                  maker_id = $8,
                  material_costs_over_all = $9,
                  work_costs_over_all = $10,
                  material_cost_spent = $11,
                  work_cost_spent = $12
where id = $13 returning *`, project.Image, project.File, project.Address, project.UserId, project.StartDate, project.EndDate, project.Active, project.Maker, project.MaterialCostsOverall, project.WorkCostsOverall, project.MaterialCostSpent, project.WorkCostSpent, project.Id).StructScan(proj)
	if err != nil {
		return nil, err
	}

	_ = tx.MustExec(`UPDATE addresses SET address = $1, user_id = $2
	WHERE project_id = $3`, project.Address, project.UserId, proj.Id)

	return proj, nil
}

func (p *Postgres) ProjectDone(project *types.Projects) (*types.ProjectManagerWithClient, error) {

	proj := new(types.ProjectManagerWithClient)

	err := p.db.QueryRowx(`UPDATE project set active = $1
WHERE id = $2 returning *`, project.Active, project.Id).StructScan(proj)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

func (p *Postgres) GetManagerAllProjects(userId int) ([]types.ProjectsWithClient, error) {

	projects := make([]types.ProjectsWithClient, 0)
	if err := p.db.Select(&projects, `SELECT id, image, file, address,user_id,
       start_date, end_date, active, maker_id, material_costs_over_all, work_costs_over_all, 
       material_cost_spent, work_cost_spent FROM project WHERE maker_id = $1`, userId); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *Postgres) GetUserAllProjects(userId int) ([]types.Projects, error) {

	projects := make([]types.Projects, 0)
	if err := p.db.Select(&projects, `SELECT id, image, file, address,user_id,
       start_date, end_date, active, maker_id, material_costs_over_all, work_costs_over_all, 
       material_cost_spent, work_cost_spent FROM project WHERE user_id = $1`, userId); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *Postgres) GetProject(projectId int) (*types.Project, error) {

	project := &types.Project{}
	if err := p.db.Get(project, `SELECT id, image, file, address,user_id,
       start_date, end_date, active, maker_id, material_costs_over_all, work_costs_over_all, 
       material_cost_spent, work_cost_spent FROM project WHERE id = $1`, projectId); err != nil {
		return nil, err
	}

	return project, nil
}

func (p *Postgres) GetProjectManager(projectId int) (*types.ProjectManager, error) {

	project := &types.ProjectManager{}
	if err := p.db.Get(project, `SELECT id, image, file, address,user_id,
       start_date, end_date, active, maker_id, material_costs_over_all, work_costs_over_all, 
       material_cost_spent, work_cost_spent FROM project WHERE id = $1`, projectId); err != nil {
		return nil, err
	}

	return project, nil
}

func (p *Postgres) GetStages(projectId int) ([]types.Stages, error) {

	stages := make([]types.Stages, 0)
	if err := p.db.Select(&stages, `SELECT id, name, project_id, phase, date FROM stages WHERE project_id = $1`, projectId); err != nil {
		return nil, err
	}

	return stages, nil
}

func (p *Postgres) GetStageByStageID(stageID int) (*types.Stage, error) {

	stage := new(types.Stage)
	if err := p.db.QueryRowx(`SELECT id, name, project_id, phase, date FROM stages WHERE id = $1`, stageID).StructScan(stage); err != nil {
		return nil, err
	}

	return stage, nil
}

func (p *Postgres) GetStagesWithCards(projectId int) ([]types.StagesWithCards, error) {

	stages := make([]types.StagesWithCards, 0)
	if err := p.db.Select(&stages, `SELECT id, name, project_id, phase, date FROM stages WHERE project_id = $1`, projectId); err != nil {
		return nil, err
	}

	return stages, nil
}

func (p *Postgres) GetCardByStageId(stageId int) ([]types.Card, error) {

	cards := make([]types.Card, 0)
	if err := p.db.Select(&cards, `SELECT id, title, images,
       deadline, stages_id, rating, description, left_to_pay, spent_on_materials, status, state FROM cards WHERE stages_id = $1`, stageId); err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *Postgres) GetCardByStageIdWithCount(stageId int) ([]types.CardWithCount, error) {

	cards := make([]types.CardWithCount, 0)
	if err := p.db.Select(&cards, `SELECT id, title, images,
       deadline, stages_id, rating, description, left_to_pay, spent_on_materials, status, state FROM cards WHERE stages_id = $1`, stageId); err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *Postgres) GetCardById(cardId int) (*types.Cards, error) {

	cards := &types.Cards{}
	if err := p.db.Get(cards, `select id, title, deadline, stages_id, rating, description, status, state, left_to_pay, spent_on_materials, images from cards where id = $1`, cardId); err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *Postgres) GetTasksByCardId(cardId int) ([]types.Task, error) {

	tasks := make([]types.Task, 0)
	if err := p.db.Select(&tasks, `select id, title, card_id, complete, images, length from tasks where card_id = $1`, cardId); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *Postgres) GetChequesByCardId(cardId int) ([]types.Cheque, error) {

	cheques := make([]types.Cheque, 0)
	if err := p.db.Select(&cheques, `select id, created_at, file, name, cost, type, card_id, user_id, project_id, length, weight from cheques where card_id = $1`, cardId); err != nil {
		return nil, err
	}

	return cheques, nil
}

func (p *Postgres) GetChequesByProjectId(projectId int) ([]types.Cheque, error) {

	cheques := make([]types.Cheque, 0)
	if err := p.db.Select(&cheques, `select id, created_at, file, name, cost, type, card_id, user_id, project_id, length, weight from cheques where project_id = $1`, projectId); err != nil {
		return nil, err
	}

	return cheques, nil
}

func (p *Postgres) AddCheque(tx *sqlx.Tx, cheques *types.AddCheque) (*types.Cheque, error) {

	cheque := new(types.Cheque)
	if err := tx.QueryRowx(`INSERT INTO cheques (file, name, cost, type, card_id, user_id, project_id, length, weight) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *`, cheques.File, cheques.Name, cheques.Cost, cheques.Type, cheques.CardId, cheques.UserId, cheques.ProjectId, cheques.Length, cheques.Weight).StructScan(cheque); err != nil {
		return nil, err
	}

	return cheque, nil
}

func (p *Postgres) DeleteCheque(tx *sqlx.Tx, cheque *types.DeleteCheque) error {

	_ = tx.MustExec(`DELETE FROM cheques WHERE id = $1 and card_id = $2 and project_id = $3`, cheque.Id, cheque.CardId, cheque.ProjectId)

	return nil
}

func (p *Postgres) GetIntProjects(projectId int) (int, error) {

	var count int
	if err := p.db.Get(&count, `SELECT count(*) from project where id = $1`, projectId); err != nil {
		return 0, err
	}

	return count, nil
}

func (p *Postgres) GetIntStage(stageId, projectId int) (int, error) {

	var count int
	if err := p.db.Get(&count, `SELECT count(*) from stages where id = $1 and project_id = $2`, stageId, projectId); err != nil {
		return 0, err
	}

	return count, nil
}

//// ДОБАВИТЬ ПРОЕКТ ///////
func (p *Postgres) ProjectAddManager(project *types.AddProjectManager, tx *sqlx.Tx) (*types.AddProjectManager, error) {

	var proj types.AddProjectManager
	if err := tx.QueryRowx(`INSERT INTO project (image, file, address, user_id, start_date, end_date, active, maker_id, material_costs_over_all, work_costs_over_all, material_cost_spent, work_cost_spent) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) returning *`, project.Image, project.File, project.Address, project.UserId, project.StartDate, project.EndDate, project.Active, project.Maker, project.MaterialCostsOverall, project.WorkCostsOverall, project.MaterialCostSpent, project.WorkCostSpent).StructScan(&proj); err != nil {
		return nil, err
	}

	_ = tx.MustExec(`INSERT INTO addresses (address, user_id, project_id) 
	VALUES ($1, $2, $3)`, project.Address, project.UserId, proj.Id)

	return &proj, nil
}

func (p *Postgres) StageAddManager(stage *types.AddStageManager, tx *sqlx.Tx) (*types.AddStageManager, error) {

	stag := new(types.AddStageManager)
	if err := tx.QueryRowx(`INSERT INTO stages (name, project_id, phase, date) 
	VALUES ($1, $2, $3, $4) returning *`, stage.Name, stage.ProjectId, stage.Phase, stage.Date).StructScan(stag); err != nil {
		return nil, err
	}

	return stag, nil
}

func (p *Postgres) CardAddManager(card *types.AddCardManager, tx *sqlx.Tx) (*types.AddCardManager, error) {

	car := new(types.AddCardManager)
	if err := tx.QueryRowx(`INSERT INTO cards (images, title, deadline, stages_id, rating, description, left_to_pay, spent_on_materials, status, state) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning *`, card.Image, card.Title, card.Deadline, card.StagesId, card.Rating, card.Description, card.LeftToPay, card.SpentOnMaterials, card.Status, card.State).StructScan(car); err != nil {
		return nil, err
	}

	return car, nil
}

func (p *Postgres) TaskAddManager(task *types.AddTaskManager, tx *sqlx.Tx) (*types.AddTaskManager, error) {

	var tas types.AddTaskManager
	if err := tx.QueryRowx(`INSERT INTO tasks (images, title, card_id, complete, length) 
	VALUES ($1, $2, $3, $4, $5) returning *`, task.Image, task.Title, task.CardId, task.Complete, task.Length).StructScan(&tas); err != nil {
		return nil, err
	}

	return &tas, nil
}

func (p *Postgres) GetCardByStageIdManager(stageId int) ([]types.AddCardManager, error) {

	cards := make([]types.AddCardManager, 0)
	if err := p.db.Select(&cards, `SELECT id, images, title, 
       deadline, stages_id, rating, description, left_to_pay, spent_on_materials, status, state FROM cards WHERE stages_id = $1`, stageId); err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *Postgres) GetTasksByCardIdManager(cardId int) ([]types.AddTaskManager, error) {

	tasks := make([]types.AddTaskManager, 0)
	if err := p.db.Select(&tasks, `select id, images, title, card_id, complete, images, length from tasks where card_id = $1`, cardId); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *Postgres) CardAddManagerProject(card *types.CardsAdd, tx *sqlx.Tx) (*types.CardsAdd, error) {

	car := new(types.CardsAdd)
	if err := tx.QueryRowx(`INSERT INTO cards (images,title, deadline, stages_id, rating, description, left_to_pay, spent_on_materials, status, state) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning *`, card.Image, card.Title, card.Deadline, card.StagesId, card.Rating, card.Description, card.LeftToPay, card.SpentOnMaterials, card.Status, card.State).StructScan(car); err != nil {
		return nil, err
	}

	return car, nil
}

func (p *Postgres) CardUpdateManagerProject(card *types.UpdateCard) (*types.Cards, error) {

	car := new(types.Cards)
	if err := p.db.QueryRowx(`UPDATE cards SET images = $1, title = $2, deadline = $3, stages_id = $4, rating = $5, description = $6, left_to_pay = $7, spent_on_materials = $8, status = $9, state = $10
	WHERE id = $11 returning *`, card.Images, card.Title, card.Deadline, card.StagesId, card.Rating, card.Description, card.LeftToPay, card.SpentOnMaterials, card.Status, card.State, card.Id).StructScan(car); err != nil {
		return nil, err
	}

	return car, nil
}

func (p *Postgres) CardUpdateManagerProjectDone(card *types.Card) (*types.Card, error) {

	car := new(types.Card)
	if err := p.db.QueryRowx(`UPDATE cards SET status = $1
	WHERE id = $2 returning *`, card.Status, card.Id).StructScan(car); err != nil {
		return nil, err
	}

	return car, nil
}

func (p *Postgres) ChequeAddManager(cheques *types.AddChequeManager, tx *sqlx.Tx) (*types.AddChequeManager, error) {

	cheque := new(types.AddChequeManager)
	if err := tx.QueryRowx(`INSERT INTO cheques (file, name, cost, type, card_id, user_id, project_id, length, weight) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *`, cheques.File, cheques.Name, cheques.Cost, cheques.Type, cheques.CardId, cheques.UserId, cheques.ProjectId, cheques.Length, cheques.Weight).StructScan(cheque); err != nil {
		return nil, err
	}

	return cheque, nil
}

func (p *Postgres) GetCardByCardId(cardId int) (*types.Card, error) {

	car := new(types.Card)
	if err := p.db.QueryRowx(`SELECT id, images, title, deadline, stages_id, rating, description, left_to_pay, spent_on_materials, status, state FROM cards
WHERE id = $1`, cardId).StructScan(car); err != nil {
		return nil, err
	}

	return car, nil
}

func (p *Postgres) GetCheque(chequeId int) (*types.Cheque, error) {

	cheq := new(types.Cheque)
	if err := p.db.QueryRowx(`SELECT id, created_at, file, name, cost, type, card_id, user_id, project_id, length, weight FROM cheques WHERE id = $1`, chequeId).StructScan(cheq); err != nil {
		return nil, err
	}

	return cheq, nil
}

func (p *Postgres) GetIntCardCheques(cardID int) (int, error) {

	var count int
	if err := p.db.Get(&count, `SELECT count(*) from cheques where card_id = $1`, cardID); err != nil {
		return 0, err
	}

	return count, nil
}

func (p *Postgres) DeleteCard(tx *sqlx.Tx, stageID, cardID int) error {

	if _, err := tx.Exec(`DELETE FROM cards WHERE id = $1 and stages_id = $2`, cardID, stageID); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteStage(tx *sqlx.Tx, stageID, projectID int) error {

	if _, err := tx.Exec(`DELETE FROM stages WHERE id = $1 and project_id = $2`, stageID, projectID); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteTask(tx *sqlx.Tx, cardID, taskID int) error {

	if _, err := tx.Exec(`DELETE FROM tasks WHERE id = $1 and card_id = $2`, taskID, cardID); err != nil {
		return err
	}

	return nil
}
