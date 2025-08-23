package email

import (
	"bytes"
	"net/smtp"
	"strings"
	"todo-app/pkg/email/layouts"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/html"
)

type Email struct {
	SmtpClient *smtp.Client
}

func NewSmtpClient(SmtpClient *smtp.Client) *Email {
	return &Email{SmtpClient: SmtpClient}
}

type ConfigSendEmail struct {
	To      string   `yml:"to"`
	Subject string   `yml:"subject"`
	Body    []string `yml:"body"`
}

type Content struct {
	Code string
	Name string
	Body string
}

func (s *Email) Send(to, subject string, c *Content, typeBody layouts.TypeBodyMail) error {
	logrus.Debug("Start push email")

	if err := s.SmtpClient.Rcpt(to); err != nil {
		logrus.Error(err.Error())
		return err
	}

	bodyMail, err := layouts.Get(typeBody)
	if err != nil {
		logrus.Fatal("error get layouts: ", typeBody)
		return err
	}

	readyBody, err := insertCodeIntoHTML(bodyMail, c, typeBody)
	if err != nil {
		logrus.Error("error insert code into HTML")
		return err
	}

	// Debug
	// err = os.WriteFile("modified_template.html", []byte(readyBody), 0644)
	// if err != nil {
	// 	panic(err)
	// }

	w, err := s.SmtpClient.Data()
	if err != nil {
		logrus.Error("Error opening data connection: ", err.Error())
		return err
	}
	defer w.Close()

	msg := messageFormation(viper.GetString("smtp.from"), to, string(readyBody), typeBody)
	_, err = w.Write(msg)
	if err != nil {
		logrus.Error("Error writing data: ", err)
	}

	if err = w.Close(); err != nil {
		logrus.Fatal("Writer close error: ", err)
		return err
	}

	logrus.Debug("email sent")

	return nil
}

/*
	Additional Methods
*/

func insertCodeIntoHTML(htmlStr string, c *Content, typeClass layouts.TypeBodyMail) ([]byte, error) {
	logrus.Debugf("start insert code into html")

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	var (
		f        func(*html.Node)
		clasName = ""
		data     = ""
	)

	switch typeClass {
	case layouts.Notification:
		{
			clasName = "name_notification"
			data = c.Name
		}
	case layouts.СonfirmationCode:
		{
			clasName = "code"
			data = c.Code
		}
	case layouts.СonfirmationUrl:
		{
			clasName = "Body_notification"
			data = c.Body
		}
	}

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && hasClass(attr.Val, clasName) {
					for n.FirstChild != nil {
						n.RemoveChild(n.FirstChild)
					}

					n.FirstChild = &html.Node{
						Type: html.TextNode,
						Data: data,
					}

					logrus.Debug("chek")
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	var buf bytes.Buffer
	if err = html.Render(&buf, doc); err != nil {
		return nil, err
	}

	logrus.Debug("end of insert code into html")

	return buf.Bytes(), err
}

func hasClass(classAttr, className string) bool {
	classes := strings.Fields(classAttr)
	for _, c := range classes {
		if c == className {
			return true
		}
	}
	return false
}

func messageFormation(from, to, body string, typeSybject layouts.TypeBodyMail) []byte {
	// logrus.Debug(to, body)
	if to == "" || body == "" {
		logrus.Error("fields are undefined")
		return []byte{}
	}
	// logrus.Debug("body: ", body)

	subject := ""
	switch typeSybject {
	case layouts.Notification:
		subject = "Confirmation of registration"
	case layouts.СonfirmationCode:
		subject = "Confirmation of registration"
	case layouts.СonfirmationUrl:
		subject = "Confirmation of registration"
	}

	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "," + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	return msg
}
