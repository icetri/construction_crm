package config

type Config struct {
	ServerPort  string        `yaml:"server_port"`
	PostgresDsn string        `yaml:"postgres_dsn"`
	JWTKey      string        `yaml:"jwtkey"`
	Telegram    *Telegram     `yaml:"telegram"`
	Path        string        `yaml:"path_for_avatar"`
	FileStorage *Storage      `yaml:"file_storage"`
	FireBase    *Fire         `yaml:"firebase"`
	Email       *ForSendEmail `yaml:"server_email"`
	SMS         *Sms          `yaml:"sms"`
}

type Telegram struct {
	TelegramToken string `yaml:"telegram_token"`
	ChatID        string `yaml:"chat_id"`
}

type Storage struct {
	Host      string `yaml:"host"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	SSL       bool   `yaml:"ssl"`
}

type Fire struct {
	Path string `yaml:"path"`
}

type ForSendEmail struct {
	EmailHost        string `yaml:"host"`
	EmailPort        int    `yaml:"port"`
	EmailLogin       string `yaml:"login"`
	EmailFrom        string `yaml:"from"`
	EmailPass        string `yaml:"pass"`
	EmailUnsubscribe string `yaml:"email_unsubscribe"`
	NameSender       string `yaml:"email_name_sender"`
	EmailSender      string `yaml:"email_sender"`
}

type Sms struct {
	SMSURL    string `yaml:"sms_url"`
	SmsPass   string `yaml:"sms_pass"`
	SmsLog    string `yaml:"sms_login"`
	SmsSender string `yaml:"sms_sender"`
}
