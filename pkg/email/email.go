package email

import (
	"crypto/tls"
	"log"
	"net/smtp"

	"github.com/sirupsen/logrus"
)

type ConfigSMTP struct {
	Host     string `yml:"host"`
	Port     string `yml:"port"`
	Username string `yml:"username"`
	Password string `yml:"password"`
	From     string `yml:"from"`
}

func New(cfg ConfigSMTP) *smtp.Client {

	client, err := smtp.Dial(cfg.Host + ":" + cfg.Port)
	if err != nil {
		logrus.Fatal("error conn. ", err.Error())
		return &smtp.Client{}
	}

	if err = client.Hello("localhost"); err != nil {
		logrus.Fatal("Hello error: ", err)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         cfg.Host,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		logrus.Fatal("TLS: ", err)
		return &smtp.Client{}
	}

	// conn, err := tls.Dial("tcp", cfg.Host+":"+cfg.Port, tlsConfig)
	// // conn, err := smtp.Dial(cfg.Host + ":" + cfg.Port)
	// if err != nil {
	// 	logrus.Fatal("error conn tls. ", err.Error())
	// 	return &smtp.Client{}
	// }
	// defer conn.Close()
	// client, err := smtp.NewClient(conn, cfg.Host)
	// if err != nil {
	// 	logrus.Error(err.Error())
	// 	return &smtp.Client{}
	// }
	// defer client.Quit()

	auth := smtp.PlainAuth(
		"",
		cfg.Username,
		cfg.Password,
		cfg.Host,
	)
	logrus.Debug("Authorization smt")

	if err = client.Auth(auth); err != nil {
		logrus.Error(err.Error())
		return &smtp.Client{}
	}

	if err = client.Mail(cfg.From); err != nil {
		logrus.Error(err.Error())
		return &smtp.Client{}
	}

	if err = client.Rcpt(cfg.From); err != nil {
		log.Fatal("RCPT TO error: ", err)
	}

	return client
}

// type ConfigEmail struct {
// 	Addr              string        `yml:"addr"`
// 	Domain            string        `yml:"domain"`
// 	WriteTimeout      time.Duration `yml:"write_timeout"`
// 	ReadTimeout       time.Duration `yml:"read_timeout"`
// 	MaxMessageBytes   int64         `yml:"max_message_bytes"`
// 	MaxRecipients     int           `yml:"max_recipients"`
// 	AllowInsecureAuth bool          `yml:"allow_insecure_auth"`
// }

// func NewServer(cfg ConfigEmail) (*smtp.Server, error) {
// 	s := smtp.NewServer(&backend.BackendSMTP{})
// 	s.Addr = ":" + cfg.Addr
// 	s.Domain = cfg.Domain
// 	s.WriteTimeout = cfg.WriteTimeout * time.Second
// 	s.ReadTimeout = cfg.ReadTimeout * time.Second
// 	s.MaxMessageBytes = cfg.MaxMessageBytes
// 	s.MaxRecipients = cfg.MaxRecipients
// 	s.AllowInsecureAuth = cfg.AllowInsecureAuth
// 	if err := s.ListenAndServe(); err != nil {
// 		return nil, err
// 	}
// 	return s, nil
// }
