package postgres

import (
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/jmoiron/sqlx"
)

func (p *Postgres) GetUserList() ([]types.PaginationListUsers, error) {

	users := make([]types.PaginationListUsers, 0)
	if err := p.db.Select(&users, `SELECT id, image, created_at,updated_at, phone,email, role, last_name, first_name, middle_name, registered_manager, code FROM users`); err != nil {
		return nil, err
	}

	return users, nil
}

func (p *Postgres) GetUserById(id int) (*types.User, error) {

	users := &types.User{}
	if err := p.db.Get(users, `SELECT id, image, created_at, updated_at, phone, email, role, last_name, first_name, middle_name, registered_manager, code FROM users WHERE id = $1`, id); err != nil {
		return nil, err
	}

	return users, nil
}

func (p *Postgres) GetUserByIdManager(id int) (*types.PaginationListUsers, error) {

	pagList := new(types.PaginationListUsers)
	if err := p.db.Get(&pagList.Hits, `SELECT count(*) from project where user_id = $1`, id); err != nil {
		return nil, err
	}

	if err := p.db.Get(pagList, `SELECT id, image, created_at, updated_at, phone, email, role, last_name, first_name, middle_name, registered_manager, code FROM users WHERE id = $1`, id); err != nil {
		return nil, err
	}

	return pagList, nil
}

func (p *Postgres) UpdateUserById(cli *types.Client, tx *sqlx.Tx) (*types.Client, error) {

	client := new(types.Client)
	if err := tx.QueryRowx(`UPDATE users SET phone = $1, email = $2, last_name = $3, first_name = $4, middle_name = $5 
WHERE id = $6 returning id, last_name, first_name, middle_name, phone, email`, cli.Phone, cli.Email, cli.LastName, cli.FirstName, cli.MiddleName, cli.Id).StructScan(client); err != nil {
		return nil, err
	}

	return cli, nil
}

func (p *Postgres) GetCountProjectUserByIdManager(id int) (int64, error) {
	var pagList int64
	if err := p.db.Get(&pagList, `SELECT count(*) from project where user_id = $1`, id); err != nil {
		return 0, err
	}

	return pagList, nil
}
