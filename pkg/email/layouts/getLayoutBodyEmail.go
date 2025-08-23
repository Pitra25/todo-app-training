package layouts

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

type TypeBodyMail int

const (
	СonfirmationUrl TypeBodyMail = iota
	СonfirmationCode
	Notification
)

func Get(typeBody TypeBodyMail) (string, error) {
	var (
		htmlBytes []byte
		err       error
	)

	logrus.Debug("start receiving the layout")

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	switch typeBody {
	case СonfirmationUrl:
		{
			fulPath := filepath.Join(dir, "verify_email.html")
			htmlBytes, err = os.ReadFile(fulPath)
			if err != nil {
				logrus.Fatal("get layout email. ", err.Error())
				os.Exit(1)
				return "", err
			}
		}
	case СonfirmationCode:
		{
			fulPath := filepath.Join(dir, "verify_code_email.html")
			htmlBytes, err = os.ReadFile(fulPath)
			if err != nil {
				logrus.Fatal("get layout code. ", err.Error())
				os.Exit(1)
				return "", err
			}
		}
	case Notification:
		{
			fulPath := filepath.Join(dir, "notification.html")
			htmlBytes, err = os.ReadFile(fulPath)
			if err != nil {
				logrus.Fatal("get layout notification. ", err.Error())
				os.Exit(1)
				return "", err
			}

		}
	}

	logrus.Debug("get layout")
	return string(htmlBytes), err
}
