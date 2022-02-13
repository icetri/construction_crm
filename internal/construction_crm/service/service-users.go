package service

import (
	"database/sql"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/construction_crm/pkg/logger"
	"github.com/pkg/errors"
)

func (s *Service) UsersList() ([]types.PaginationListUsers, error) {

	users, err := s.db.GetUserList()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserList in UsersList"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i, val := range users {
		addresses, err := s.db.GetAddresses(val.ID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetAddresses in UsersList"))
			return nil, infrastruct.ErrorInternalServerError
		}
		users[i].Address = addresses

		count, err := s.db.GetCountProjectUserByIdManager(val.ID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCountProjectUserByIdManager in UsersList"))
			return nil, infrastruct.ErrorInternalServerError
		}
		users[i].Hits = count
	}

	return users, nil
}

func (s *Service) UserById(userID int) (*types.PaginationListUsers, error) {

	user, err := s.db.GetUserByIdManager(userID)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserByIdManager in UserById"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if user == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	addresses, err := s.db.GetAddresses(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAddresses in UserById"))
		return nil, infrastruct.ErrorInternalServerError
	}

	user.Address = addresses

	return user, nil
}
