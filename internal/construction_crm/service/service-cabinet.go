package service

import (
	"database/sql"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/construction_crm/pkg/logger"
	"github.com/pkg/errors"
)

func (s *Service) GetCabinet(id int) (*types.CabinetNew, error) {

	cabinet, err := s.db.GetCabinet(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinet in GetCabinet"))
		return nil, infrastruct.ErrorInternalServerError
	}

	addresses, err := s.db.GetAddresses(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAddresses in GetCabinet"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cabinet.Addresses = append(cabinet.Addresses, addresses...)

	return cabinet, nil
}

func (s *Service) PutCabinetUserInfo(user *types.PutUserInfo) (*types.CabinetNew, error) {

	err := s.db.PutCabinet(user)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with PutCabinet in PutCabinetUserInfo"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cabinet, err := s.db.GetCabinet(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinet in PutCabinetUserInfo"))
		return nil, infrastruct.ErrorInternalServerError
	}

	addresses, err := s.db.GetAddresses(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAddresses in GetCabinet"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cabinet.Addresses = append(cabinet.Addresses, addresses...)

	return cabinet, nil
}

func (s *Service) PutUserPhone(user *types.PutUserPhone) error {

	userInDb, err := s.db.GetUserById(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserById in PutUserPhone"))
		return infrastruct.ErrorInternalServerError
	}

	userInDb.Phone = user.NewPhone

	if err := s.makeConfirmationPhoneCode(userInDb); err != nil {
		logger.LogError(errors.Wrap(err, "err with makeConfirmationPhoneCode in PutUserPhone"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) PutUserValidatePhone(user *types.PutUserCont) (*types.CabinetNew, error) {

	userId, err := s.db.GetUserByCode(user.Code)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserByCode in PutUserValidatePhone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if userId == nil {
		return nil, infrastruct.ErrorCodeIsIncorrect
	}

	userId.Phone = user.Phone

	if err := s.db.UpdatePhone(userId); err != nil {
		logger.LogError(errors.Wrap(err, "err with UpdatePhone in PutUserValidatePhone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cabinet, err := s.db.GetCabinet(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCabinet in PutUserValidatePhone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	addresses, err := s.db.GetAddresses(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAddresses in GetCabinet"))
		return nil, infrastruct.ErrorInternalServerError
	}

	cabinet.Addresses = append(cabinet.Addresses, addresses...)

	return cabinet, nil
}
