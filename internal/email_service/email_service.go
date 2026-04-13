package email_service

type IEmailService interface {
	Send(email, text string) error
}

type EmailService struct{}

func CreateEmailService() IEmailService {
	return &EmailService{}
}

func (s *EmailService) Send(email, text string) error {
	return nil
}