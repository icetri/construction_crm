package mail

import (
	"fmt"
	"github.com/construction_crm/internal/clients/postgres"
	"github.com/construction_crm/internal/construction_crm/types/config"
	"github.com/construction_crm/pkg/logger"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"time"
)

type Mail struct {
	email *config.ForSendEmail
	db    *postgres.Postgres
}

func NewMail(cnf *config.ForSendEmail, db *postgres.Postgres) (*Mail, error) {
	return &Mail{
		email: cnf,
		db:    db,
	}, nil
}

func (r *Mail) SendMail(subject, text, to string) error {

	m := gomail.NewMessage()

	m.SetAddressHeader("From", r.email.EmailSender, r.email.NameSender)
	m.SetAddressHeader("To", to, to)

	m.SetHeader("From", fmt.Sprintf("%s <%s>", r.email.NameSender, r.email.EmailSender))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetHeader("MIME-Version:", "1.0")
	m.SetHeader("Reply-To", r.email.EmailUnsubscribe)

	m.SetBody("text/plain", text)

	d := gomail.NewDialer(r.email.EmailHost, r.email.EmailPort, r.email.EmailLogin, r.email.EmailPass)

	stopMail := 0
	for stopMail < 1 {
		stopMail++

		if err := d.DialAndSend(m); err != nil {

			time.Sleep(time.Second * 30)
			if stopMail < 1 {
				continue
			}

			errInsert := r.db.InsertProblemEmail(to, err.Error())
			if errInsert != nil {
				logger.LogError(errors.Wrap(errInsert, "err with errInsert"))
			}

			return errors.Wrap(err, "err with DialAndSend mail")
		}

		break
	}

	return nil
}
