package methods

import (
	"fmt"
	"todo-app/internal/repository"
	storage "todo-app/pkg/cache/redis"
	"todo-app/pkg/email"
	"todo-app/pkg/email/layouts"
	"todo-app/pkg/middleware"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type EmailsService struct {
	repo    repository.Emails
	eClient *email.Email
	rdb     *storage.RedisCLientDB
}

func NewEmailService(
	repo repository.Emails,
	eClient *email.Email,
	rdb *redis.Client,
) *EmailsService {
	return &EmailsService{
		repo:    repo,
		eClient: eClient,
		rdb:     &storage.RedisCLientDB{Db: rdb},
	}
}

func (s *EmailsService) SendEmail(to string, userId int) error {

	code, err := middleware.CodeGeneration(6)
	if err != nil {
		return fmt.Errorf("failed to generate code")
	}

	logrus.Debug("Generated code for userId: ", userId, " code: ", code)

	err = s.repo.SaveCodeUser(fmt.Sprint(code), userId)
	if err != nil {
		return fmt.Errorf("failed to save code for user %d", userId)
	}

	if err := s.eClient.Send(to, "", &email.Content{Code: fmt.Sprint(code)}, layouts.СonfirmationCode); err != nil {
		logrus.Fatal("error send email: ", err.Error())
		return err
	}

	return nil
}

func (s *EmailsService) ConfirmationEmail(code string, userId int) error {
	codeUser, err := s.repo.GetCodeUser(userId)
	if err != nil {
		return fmt.Errorf("failed to get code user: %d", userId)
	}

	if codeUser.Code != code {
		return fmt.Errorf("code mismatch for user %d", userId)
	}

	err = s.repo.UpdateStatusUser(userId)
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("error update status")
	}

	err = s.rdb.DeleteRecord(userId, storage.Code_user)
	if err != nil {
		logrus.Error("Error: Delete Recording in Redis")
	}

	err = s.repo.DeleteRecord(codeUser.Id, userId)
	if err != nil {
		logrus.Error("Error: Delete Recording in Db")
	}

	return nil
}
