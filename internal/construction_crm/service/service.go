package service

import (
	"fmt"
	"github.com/construction_crm/internal/clients/postgres"
	"github.com/construction_crm/internal/construction_crm/service/firebase"
	"github.com/construction_crm/internal/construction_crm/service/mail"
	"github.com/construction_crm/internal/construction_crm/service/minio"
	"github.com/construction_crm/internal/construction_crm/types"
	"github.com/construction_crm/internal/construction_crm/types/config"
	"github.com/construction_crm/pkg/infrastruct"
	"github.com/construction_crm/pkg/logger"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	db     *postgres.Postgres
	JWTkey string
	path   string
	fs     *minio.FileStorage
	fb     *firebase.FireBase
	email  *mail.Mail
	sms    *config.Sms
}

func NewService(db *postgres.Postgres, cfg *config.Config, minio *minio.FileStorage, firebase *firebase.FireBase, email *mail.Mail, sms *config.Sms) *Service {
	return &Service{
		db:     db,
		JWTkey: cfg.JWTKey,
		path:   cfg.Path,
		fs:     minio,
		fb:     firebase,
		email:  email,
		sms:    sms,
	}
}

const (
	path string = "https://minio.4k-remont.ru/%s/%s"
)

func delSpaceRegister(user *types.Register) {
	user.Phone = strings.ReplaceAll(user.Phone, " ", "")
	user.FirstName = strings.ReplaceAll(user.FirstName, " ", "")
	user.LastName = strings.ReplaceAll(user.LastName, " ", "")
	user.MiddleName = strings.ReplaceAll(user.MiddleName, " ", "")
	user.Email = strings.ReplaceAll(user.Email, " ", "")
}

func (s *Service) makeConfirmationPhoneCode(user *types.User) error {

	if err := s.db.DeleteConfPhone(user.Code, user.ID); err != nil {
		logger.LogError(errors.Wrap(err, "err with DeleteConfPhone in makeConfirmationPhoneCode"))
		return infrastruct.ErrorInternalServerError
	}

	//генерируем код
	rand.Seed(time.Now().UnixNano())
	num := 999 + rand.Intn(9000)
	code := strconv.Itoa(num)

	if err := s.db.AddCodeForConfPhone(code, user.ID); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddCodeForConfPhone in makeConfirmationPhoneCode"))
		return infrastruct.ErrorInternalServerError
	}

	mess := make([]types.SMSMessages, 0)
	mes := types.SMSMessages{Phone: user.Phone, Sender: s.sms.SmsSender, ClientID: user.ID, Text: code}
	mess = append(mess, mes)

	resp, err := s.SendSMS(mess)
	if err != nil {
		return infrastruct.ErrorInternalServerError
	}

	if err := s.db.AddLogRespSMS(string(resp), user.ID); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogRespSMS in makeConfirmationPhoneCode"))
	}

	//	logger.LogInfo(code)

	return nil
}

func (s *Service) makeConfirmationEmailCode(user *types.User) error {

	if err := s.db.DeleteConfEmail(user.Code, user.Email); err != nil {
		logger.LogError(errors.Wrap(err, "err with DeleteConfPhone in makeConfirmationEmailCode"))
		return infrastruct.ErrorInternalServerError
	}

	//генерируем код
	rand.Seed(time.Now().UnixNano())
	num := 999 + rand.Intn(9000)
	code := strconv.Itoa(num)

	if err := s.db.AddCodeForConfEmail(code, user.Email); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddCodeForConfPhone in makeConfirmationEmailCode"))
		return infrastruct.ErrorInternalServerError
	}

	go func() {

		if err := s.email.SendMail(mail.SubjectEmail, fmt.Sprintf(mail.BodyRegisterText, code), user.Email); err != nil {
			logger.LogError(errors.Wrap(err, fmt.Sprintf("err with SendMail in makeConfirmationEmailCode for user_email = %s", user.Email)))
		}

	}()

	return nil
}

func (s *Service) GetURLFile(path, bucket, name string) string {
	return fmt.Sprintf(path, bucket, name)
}

func (s *Service) timeCheck(be, af, cur string) (bool, error) {

	start := strings.Split(be, ".")
	end := strings.Split(af, ".")
	current := strings.Split(cur, ".")

	if len(start) != 3 || len(end) != 3 || len(current) != 3 {
		return false, infrastruct.ErrorBadRequest
	}

	x := start[2] + start[1] + start[0]
	y := end[2] + end[1] + end[0]
	z := current[2] + current[1] + current[0]

	st, err := strconv.Atoi(x)
	if err != nil {
		return false, err
	}

	en, err := strconv.Atoi(y)
	if err != nil {
		return false, err
	}

	cu, err := strconv.Atoi(z)
	if err != nil {
		return false, err
	}

	if st <= cu && cu <= en {
		return true, nil
	}

	return false, nil
}

func (s *Service) timeCheckProject(be, af string) (bool, error) {

	start := strings.Split(be, ".")
	end := strings.Split(af, ".")

	if len(start) != 3 || len(end) != 3 {
		return false, infrastruct.ErrorBadRequest
	}

	x := start[2] + start[1] + start[0]
	y := end[2] + end[1] + end[0]

	st, err := strconv.Atoi(x)
	if err != nil {
		return false, err
	}

	en, err := strconv.Atoi(y)
	if err != nil {
		return false, err
	}

	if st < en {
		return true, nil
	}

	return false, nil
}

func (s *Service) TestMail() error {

	go func() {

		if err := s.email.SendMail(mail.SubjectEmail, "Ваш код", "icri31@gmail.com"); err != nil {
			logger.LogError(errors.Wrap(err, fmt.Sprintf("err with SendMail in makeConfirmationEmailCode for user_email = %s", "icetri31@gmail.com")))
		}

	}()

	return nil
}
