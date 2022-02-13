package service

import (
	"database/sql"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/construction_crm/pkg/logger"
	"github.com/pkg/errors"
	"strings"
)

func (s *Service) CheckRegister(userPhone string) (*types.User, error) {
	userMan, err := s.checkDoubleUserManager(userPhone)
	if err != nil {
		return userMan, err
	}
	return nil, err
}

func (s *Service) RegisterNewUser(user *types.Register) error {

	user.Email = strings.ToLower(user.Email)
	delSpaceRegister(user)

	if err := s.checkDoubleUserByPhone(user.Phone); err != nil {
		return err
	}

	if err := s.checkDoubleUserByEmail(user.Email); err != nil {
		return err
	}

	err := s.db.CreateUser(user)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterNewUser"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) AuthUserPhone(auth *types.AuthCodePhone) (*types.Token, error) {

	user, err := s.db.GetUserByCodePhone(auth.Code, auth.Phone)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserByCodePhone in AuthUserPhone"))
			return nil, infrastruct.ErrorInternalServerError
		}
		return nil, infrastruct.ErrorCodeIsIncorrect
	}

	if user.Code != auth.Code {
		return nil, infrastruct.ErrorCodeIsIncorrect
	}

	token, err := infrastruct.GenerateJWT(user.ID, user.Role, s.JWTkey)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with infrastruct.GenerateJWT"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{Token: token}, nil
}

func (s *Service) AuthUserEmail(auth *types.AuthCodeEmail) (*types.Token, error) {

	user, err := s.db.GetUserByCodeEmail(auth.Code, auth.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserByCodeEmail in AuthUserEmail"))
			return nil, infrastruct.ErrorInternalServerError
		}
		return nil, infrastruct.ErrorCodeIsIncorrect
	}

	if user.Code != auth.Code {
		return nil, infrastruct.ErrorCodeIsIncorrect
	}

	token, err := infrastruct.GenerateJWT(user.ID, user.Role, s.JWTkey)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with infrastruct.GenerateJWT"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{Token: token}, nil
}

func (s *Service) RegisterUserManager(user *types.Register) error {

	user.Email = strings.ToLower(user.Email)
	delSpaceRegister(user)

	if err := s.checkDoubleUserByPhone(user.Phone); err != nil {
		return err
	}

	if err := s.checkDoubleUserByEmail(user.Email); err != nil {
		return err
	}

	user.RegManager = true
	err := s.db.CreateUser(user)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterUserManager"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) checkDoubleUserByPhone(phone string) error {

	user, err := s.db.GetUserByPhone(phone)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserByPhone in checkDoubleUserByPhone"))
		return infrastruct.ErrorInternalServerError
	}

	if user != nil {
		return infrastruct.ErrorPhoneIsExist
	}

	return nil
}

func (s *Service) checkDoubleUserManager(phone string) (*types.User, error) {

	user, err := s.db.GetUserByPhone(phone)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserByPhone in checkDoubleUserManager"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if user == nil {
		return nil, nil
	}

	if user.RegManager == true {
		return user, infrastruct.ErrorPhoneIsExistManager
	}

	return nil, nil
}

func (s *Service) checkDoubleUserByEmail(email string) error {

	user, err := s.db.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetUserByPhone in checkDoubleUserByEmail"))
		return infrastruct.ErrorInternalServerError
	}

	if user != nil {
		return infrastruct.ErrorEmailIsExist
	}

	return nil
}

func (s *Service) CodePhone(auth *types.AuthPhone) error {

	auth.Phone = strings.ReplaceAll(auth.Phone, " ", "")

	user, err := s.db.GetUserByPhone(auth.Phone)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserByPhone in CodePhone"))
			return infrastruct.ErrorInternalServerError
		}
		return infrastruct.ErrorPhoneIsIncorrect
	}

	if err = s.makeConfirmationPhoneCode(user); err != nil {
		return err
	}

	return nil
}

func (s *Service) CodeEmail(auth *types.AuthEmail) error {

	auth.Email = strings.ReplaceAll(auth.Email, " ", "")

	user, err := s.db.GetUserByEmail(auth.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserByEmail in CodeEmail"))
			return infrastruct.ErrorInternalServerError
		}
		return infrastruct.ErrorEmailIsIncorrect
	}

	if err = s.makeConfirmationEmailCode(user); err != nil {
		return err
	}

	return nil
}

func (s *Service) AuthManagerEmail(auth *types.AuthEmailManager) (*types.Token, error) {

	auth.Email = strings.ToLower(auth.Email)
	auth.Email = strings.ReplaceAll(auth.Email, " ", "")

	user, err := s.db.GetManagerByEmail(auth.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetManagerByEmail in AuthManagerEmail"))
			return nil, infrastruct.ErrorInternalServerError
		}
		return nil, infrastruct.ErrorPasswordOrEmailIsIncorrect
	}

	if user.Password != auth.Password {
		return nil, infrastruct.ErrorPasswordOrEmailIsIncorrect
	}

	token, err := infrastruct.GenerateJWT(user.ID, user.Role, s.JWTkey)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with infrastruct.GenerateJWT"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{Token: token}, nil
}
