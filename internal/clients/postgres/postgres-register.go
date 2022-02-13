package postgres

import (
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/pkg/errors"
)

func (p *Postgres) CreateUser(user *types.Register) error {

	_, err := p.db.Exec(`INSERT INTO users (phone, email, last_name, first_name, middle_name, registered_manager) 
	VALUES ($1, $2, $3, $4, $5, $6)`, user.Phone, user.Email, user.LastName, user.FirstName, user.MiddleName, user.RegManager)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateUserByManager(managerId int, userReg *types.RegistredUserByManagerProject) error {
	tx, err := p.db.Begin()
	if err != nil {
		return errors.Wrap(err, "err with Begin in CreateUserByManager")
	}
	defer tx.Rollback()

	var userId int
	err = tx.QueryRow(`INSERT INTO users (phone, email, last_name, first_name, middle_name, registered_manager) 
	VALUES ($1, $2, $3, $4, $5, $6) returning id`, userReg.Phone, userReg.Email, userReg.LastName, userReg.FirstName, userReg.MiddleName, userReg.RegManager).Scan(&userId)
	if err != nil {
		return errors.Wrap(err, "err with Exec users in CreateUserByManager")
	}

	_, err = tx.Exec("INSERT INTO manager_list (manager_id,user_id) VALUES ($1,$2)", managerId, userId)
	if err != nil {
		return errors.Wrap(err, "err with Exec manager_list in CreateUserByManager")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "err with Commit in CreateUserByManager")
	}

	return nil
}

func (p *Postgres) AddManagerList(managerId, userId int) error {

	_, err := p.db.Exec("INSERT INTO manager_list (manager_id,user_id) VALUES ($1,$2)", managerId, userId)
	if err != nil {
		return errors.Wrap(err, "err with Exec manager_list in AddManagerList")
	}

	return nil
}

func (p *Postgres) DeleteConfPhone(phoneCode string, userID int) error {

	phoneCode = ""
	_, err := p.db.Exec(`UPDATE users SET code = $1 WHERE id = $2`, phoneCode, userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) AddCodeForConfPhone(phoneCode string, userID int) error {

	_, err := p.db.Exec(`UPDATE users SET code = $1 WHERE id = $2`, phoneCode, userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteConfEmail(emailCode, email string) error {

	emailCode = ""
	_, err := p.db.Exec(`UPDATE users SET code = $1 WHERE email = $2`, emailCode, email)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) AddCodeForConfEmail(emailCode, email string) error {

	_, err := p.db.Exec(`UPDATE users SET code = $1 WHERE email = $2`, emailCode, email)
	if err != nil {
		return err
	}

	return nil
}
