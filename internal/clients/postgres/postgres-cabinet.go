package postgres

import (
	"github.com/construction_crm/internal/construction_crm/types"
)

func (p *Postgres) GetCabinet(id int) (*types.CabinetNew, error) {

	cabinet := &types.CabinetNew{}
	err := p.db.Get(cabinet, "SELECT id, image, phone,email,last_name,first_name,middle_name FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return cabinet, nil
}

func (p *Postgres) GetAddresses(id int) ([]types.Address, error) {

	addresses := make([]types.Address, 0)
	err := p.db.Select(&addresses, "SELECT id, address, user_id, project_id FROM addresses WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (p *Postgres) PutCabinet(user *types.PutUserInfo) error {

	_, err := p.db.Exec("UPDATE users SET last_name = $1, first_name = $2, middle_name = $3, email = $4, image = $5, updated_at = NOW() WHERE id = $6", user.LastName, user.FirstName, user.MiddleName, user.Email, user.Image, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdatePhone(user *types.User) error {

	_, err := p.db.Exec("UPDATE users SET phone = $1, updated_at = NOW() WHERE id = $2", user.Phone, user.ID)
	if err != nil {
		return err
	}

	return nil
}
