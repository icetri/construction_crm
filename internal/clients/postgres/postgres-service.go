package postgres

import (
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/jmoiron/sqlx"
)

func (p *Postgres) GetUserByPhone(phone string) (*types.User, error) {

	user := &types.User{}
	if err := p.db.QueryRowx("SELECT id, image, created_at, updated_at, phone, email, role, last_name, first_name, middle_name, registered_manager, code FROM users WHERE phone = $1", phone).StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Postgres) GetUserByCode(code string) (*types.User, error) {

	user := &types.User{}
	if err := p.db.QueryRowx("SELECT id, created_at, updated_at, phone, email, role, last_name, first_name, middle_name, registered_manager, code FROM users WHERE code = $1", code).StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Postgres) GetUserByEmail(email string) (*types.User, error) {

	user := &types.User{}
	if err := p.db.QueryRowx("SELECT id, created_at, updated_at, phone, email, role, last_name, first_name, middle_name, registered_manager, code FROM users WHERE email = $1", email).StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Postgres) GetUserByCodePhone(code, phone string) (*types.User, error) {

	user := &types.User{}
	if err := p.db.QueryRowx("SELECT id, created_at, updated_at, phone, email, role, last_name, first_name, middle_name, registered_manager, code FROM users WHERE code = $1 and phone = $2", code, phone).StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Postgres) GetUserByCodeEmail(code, email string) (*types.User, error) {

	user := &types.User{}
	if err := p.db.QueryRowx("SELECT id, created_at, updated_at, phone, email, role, last_name, first_name, middle_name, registered_manager, code FROM users WHERE code = $1 and email = $2", code, email).StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Postgres) GetManagerByEmail(email string) (*types.Manager, error) {

	user := &types.Manager{}
	if err := p.db.QueryRowx("SELECT id, image, email, password, last_name, first_name, middle_name, phone, country, city, role, created_at, updated_at FROM managers WHERE email = $1", email).StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *Postgres) UpdateProjectMaterialCostSpent(tx *sqlx.Tx, cost int, projectId int) error {

	_ = tx.MustExec(`UPDATE project SET material_cost_spent = $1 WHERE id = $2`, cost, projectId)

	return nil
}

func (p *Postgres) UpdateProjectWorkCostSpent(tx *sqlx.Tx, cost int, projectId int) error {

	_ = tx.MustExec(`UPDATE project SET work_cost_spent = $1 WHERE id = $2`, cost, projectId)

	return nil
}

func (p *Postgres) UpdateProjectCostSpent(tx *sqlx.Tx, costWork, costMaterial, projectId int) error {

	_ = tx.MustExec(`UPDATE project SET work_cost_spent = $1, material_cost_spent = $2 WHERE id = $3`, costWork, costMaterial, projectId)

	return nil
}

func (p *Postgres) CheckStageDate(tx *sqlx.Tx, projectId int, date string) (*types.Stage, error) {

	stage := new(types.Stage)
	if err := tx.QueryRowx(`SELECT id, name, project_id, phase, date FROM stages WHERE project_id = $1 AND date = $2`, projectId, date).StructScan(stage); err != nil {
		return nil, err
	}

	return stage, nil
}

func (p *Postgres) CheckStageDateForUP(tx *sqlx.Tx, projectId int) ([]types.Stage, error) {

	stage := make([]types.Stage, 0)
	if err := tx.Select(&stage, `SELECT id, name, project_id, phase, date FROM stages WHERE project_id = $1`, projectId); err != nil {
		return nil, err
	}

	return stage, nil
}

func (p *Postgres) CheckStageOldDate(tx *sqlx.Tx, projectId int, id int) (*types.Stage, error) {

	stage := new(types.Stage)
	if err := tx.QueryRowx(`SELECT id, name, project_id, phase, date FROM stages WHERE project_id = $1 AND id = $2`, projectId, id).StructScan(stage); err != nil {
		return nil, err
	}

	return stage, nil
}
