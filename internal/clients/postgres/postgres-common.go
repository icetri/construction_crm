package postgres

import (
	"github.com/construction_crm/internal/construction_crm/service/minio"
	"github.com/construction_crm/internal/construction_crm/types"
)

func (p *Postgres) UploadImageUser(minio *minio.File) (*types.FileInfo, error) {

	fileInfo := new(types.FileInfo)
	err := p.db.QueryRowx(`INSERT INTO files (url, length, mime, bucket, object, role, user_id, tag) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) returning id, created_at, url, length, mime, object, role, user_id, tag`, minio.Url, minio.Length,
		minio.MimeType, minio.Bucket, minio.Object, minio.Role, minio.UserId, minio.Tag).Scan(&fileInfo.ID, &fileInfo.Date, &fileInfo.Name, &fileInfo.Length, &fileInfo.MimeType, &fileInfo.Object, &fileInfo.Role, &fileInfo.UserId, &fileInfo.Tag)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

func (p *Postgres) GetProjectCalendar(projectId int) (*types.ProjectCalendar, error) {

	project := &types.ProjectCalendar{}
	if err := p.db.Get(project, `SELECT id, address, start_date, end_date, active, maker_id FROM project WHERE id = $1`, projectId); err != nil {
		return nil, err
	}

	return project, nil
}

func (p *Postgres) GetUserByProjectIdCalendar(id int) (*types.UserCalendar, error) {

	users := &types.UserCalendar{}
	if err := p.db.Get(users, `SELECT u.id, last_name, first_name, middle_name FROM users u left join project p on u.id = p.user_id where p.id = $1`, id); err != nil {
		return nil, err
	}

	return users, nil
}

func (p *Postgres) GetStagesCalendar(projectId int) ([]types.StagesCalendar, error) {

	stages := make([]types.StagesCalendar, 0)
	if err := p.db.Select(&stages, `SELECT id, name, project_id, phase, date FROM stages WHERE project_id = $1`, projectId); err != nil {
		return nil, err
	}

	return stages, nil
}

func (p *Postgres) GetCardByStageIdCalendar(stageId int) ([]types.CardCalendar, error) {

	cards := make([]types.CardCalendar, 0)
	if err := p.db.Select(&cards, `SELECT id, title, 
       deadline, stages_id, status FROM cards WHERE stages_id = $1`, stageId); err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *Postgres) GetManagerAllProjectsCalendar(userId int) ([]types.ProjectCalendar, error) {

	projects := make([]types.ProjectCalendar, 0)
	if err := p.db.Select(&projects, `SELECT id, address, start_date, end_date, active, maker_id 
FROM project WHERE maker_id = $1`, userId); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *Postgres) DeviceUpdateUser(token string, userId int) error {

	if _, err := p.db.Exec(`UPDATE users SET token = $1 WHERE id = $2`, token, userId); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeviceUpdateManager(token string, userId int) error {

	if _, err := p.db.Exec(`UPDATE managers SET token = $1 WHERE id = $2`, token, userId); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetProjectByCardId(cardId int) (*types.Projects, error) {

	// p.id, p.address, p.user_id, p.start_date, p.end_date, p.active, p.maker_id, p.material_costs_over_all, p.work_costs_over_all, p.material_cost_spent, p.work_cost_spent, p.image, p.file

	project := new(types.Projects)
	if err := p.db.QueryRowx(`SELECT DISTINCT p.id, p.address, p.user_id, p.start_date, p.end_date, p.active, p.maker_id, p.material_costs_over_all, p.work_costs_over_all, p.material_cost_spent, p.work_cost_spent, p.image, p.file from project p
left join stages s on p.id = s.project_id
left join cards c on s.id = c.stages_id
WHERE c.id = $1`, cardId).StructScan(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (p *Postgres) GetProjectByStageId(stageId int) (*types.Projects, error) {

	// p.id, p.address, p.user_id, p.start_date, p.end_date, p.active, p.maker_id, p.material_costs_over_all, p.work_costs_over_all, p.material_cost_spent, p.work_cost_spent, p.image, p.file

	project := new(types.Projects)
	if err := p.db.QueryRowx(`SELECT DISTINCT p.id, p.address, p.user_id, p.start_date, p.end_date, p.active, p.maker_id, p.material_costs_over_all, p.work_costs_over_all, p.material_cost_spent, p.work_cost_spent, p.image, p.file from project p
left join stages s on p.id = s.project_id
WHERE s.id = $1`, stageId).StructScan(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (p *Postgres) GetUserToken(userId int) (*types.DeviceToken, error) {

	token := new(types.DeviceToken)
	if err := p.db.QueryRowx(`SELECT token FROM users WHERE id = $1`, userId).StructScan(token); err != nil {
		return nil, err
	}

	return token, nil
}

func (p *Postgres) AddLogRespSMS(text string, userID int) error {

	_, err := p.db.Exec("INSERT INTO log_sms (resp, user_id) VALUES ($1, $2)", text, userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) InsertProblemEmail(email, data string) error {

	_, err := p.db.Exec("INSERT INTO email_problems (email_name, email_data) VALUES ($1, $2)", email, data)
	if err != nil {
		return err
	}

	return nil
}
