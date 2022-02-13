package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/construction_crm/pkg/logger"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func (s *Service) UploadFile(file *multipart.FileHeader, bucket string, claims *infrastruct.CustomClaims) (*types.FileInfo, error) {

	info, err := s.fs.Add(file, bucket, claims)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Add in UploadFile"))
		return nil, infrastruct.ErrorInternalServerError
	}

	url := s.GetURLFile(path, info.Bucket, info.Tag)

	info.Url = url

	fileInfo, err := s.db.UploadImageUser(info)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with UploadFile in UploadFile"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return fileInfo, nil
}

func (s *Service) GetUserCalendar(projectId, userID int) (*types.ProjectCalendar, error) {

	project, err := s.db.GetProjectCalendar(projectId)
	if err != nil && err != sql.ErrNoRows {
		logger.LogError(errors.Wrap(err, "err with GetProject in GetUserCalendar"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if project == nil {
		return nil, infrastruct.ErrorBadRequest
	}

	user, err := s.db.GetUserByProjectIdCalendar(projectId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByIdCalendar in GetUserCalendar"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if user.ID != userID {
		return nil, infrastruct.ErrorBadRequest
	}

	project.User = user

	stages, err := s.db.GetStagesCalendar(projectId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStages in GetUserProject"))
		return nil, infrastruct.ErrorInternalServerError
	}

	project.Stages = append(project.Stages, stages...)

	for i, stage := range project.Stages {
		cards, err := s.db.GetCardByStageIdCalendar(stage.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetCardsByStageId in GetUserProject"))
			return nil, infrastruct.ErrorInternalServerError
		}
		project.Stages[i].Cards = cards
	}

	return project, nil
}

func (s *Service) GetManagerCalendar(userId int) ([]types.ProjectCalendar, error) {

	proj, err := s.db.GetManagerAllProjectsCalendar(userId)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetProject in GetUserCalendar"))
		return nil, infrastruct.ErrorInternalServerError
	}

	projects := make([]types.ProjectCalendar, 0)

	for _, project := range proj {
		user, err := s.db.GetUserByProjectIdCalendar(project.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetUserByIdCalendar in GetUserCalendar"))
			return nil, infrastruct.ErrorInternalServerError
		}

		project.User = user

		stages, err := s.db.GetStagesCalendar(project.Id)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetStages in GetUserProject"))
			return nil, infrastruct.ErrorInternalServerError
		}

		project.Stages = append(project.Stages, stages...)

		for i, stage := range project.Stages {
			cards, err := s.db.GetCardByStageIdCalendar(stage.Id)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with GetCardsByStageId in GetUserProject"))
				return nil, infrastruct.ErrorInternalServerError
			}
			project.Stages[i].Cards = cards
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (s *Service) UpdateDeviceInfo(userToken *types.DeviceToken, user *infrastruct.CustomClaims) error {

	if user.Role == "USER" {
		err := s.db.DeviceUpdateUser(userToken.Token, user.UserID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with DeviceUpdateUser in UpdateDeviceInfo"))
			return infrastruct.ErrorInternalServerError
		}
	}

	if user.Role == "MANAGER" {
		err := s.db.DeviceUpdateManager(userToken.Token, user.UserID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with DeviceUpdateManager in UpdateDeviceInfo"))
			return infrastruct.ErrorInternalServerError
		}
	}

	return nil
}

func (s *Service) SendSMS(mes []types.SMSMessages) ([]byte, error) {
	smsBody := types.Sms{Login: s.sms.SmsLog,
		Password: s.sms.SmsPass,
		Messages: mes,
	}

	jsonStr, err := json.Marshal(smsBody)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Marshal in SendSMS"))
		return nil, err
	}

	req, err := http.NewRequest("POST", s.sms.SMSURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logger.LogError(errors.Wrap(err, "Err with http.NewRequest in SendSMS"))
		return nil, infrastruct.ErrorInternalServerError
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//логируем респонс
	respLogerBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with read body in ChoiceCertificate for logger"))
	}

	defer resp.Body.Close()

	return respLogerBytes, nil
}
