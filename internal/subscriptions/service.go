package subscriptions

import (
	"context"
	"git_email_subscriber/internal/github_api"
	"git_email_subscriber/internal/email_service"
	"net/mail"
	"github.com/google/uuid"
)

type ISubscriptionService interface {
	Subscribe(ctx context.Context, repo, email string) (*Subscription, error)
	Confirm(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error
	GetSubscriptions(ctx context.Context, email string) ([]Subscription, error)
}

type SubscriptionService struct {
	gitHubApi github_api.IGitHubApi
	repo ISubscriptionRepository
	emailService email_service.IEmailService
}

func NewSubscriptionService(gitHubApi github_api.IGitHubApi, repo ISubscriptionRepository, service email_service.IEmailService) ISubscriptionService{
	return &SubscriptionService{gitHubApi, repo, service}
}

func (s *SubscriptionService) Subscribe(ctx context.Context, repo, email string) (*Subscription, error) {
	// 1. Validate Email
	isEmailValid := isValidEmail(email)
	if !isEmailValid {
		return nil, ErrInvalidEmail
	}

	// 2. Validate Repo
	repoExists, err := s.gitHubApi.RepoExists(ctx, repo)

	if err != nil {
		return nil, err
	}

	if !repoExists {
		return nil, ErrInvalidRepo
	}

	// 3. Check Duplicates
	isCreated, err := s.repo.Exists(ctx, repo, email)

	if isCreated {
		return nil, ErrAlreadySubscribed
	}

	// 4. Create Sub
	sub := Subscription{
		ID:               uuid.NewString(),
		Email:            email,
		Repo:             repo,
		Confirmed:        false,
		ConfirmToken:     uuid.NewString(),
		UnsubscribeToken: uuid.NewString(),
	}

	// 5. Save
	err = s.repo.Create(ctx, sub)
	if err != nil {
		return nil, err
	}

	// 6. Send Email
	err = s.emailService.Send(email, sub.ConfirmToken)
	if err != nil {
		return nil, ErrEmailConfirmation
	}

	return nil, nil
}

func isValidEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

func (s *SubscriptionService) Confirm(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	sub, err := s.repo.FindByConfirmToken(ctx, token)
	if err != nil {
		return ErrNotFound
	}

	if sub.Confirmed {
		return ErrSubscriptionConfirmed
	}

	err = s.repo.Confirm(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) error {
	if token == "" {
		return ErrInvalidToken
	}

	err := s.repo.DeleteByUnsubscribeToken(ctx, token)
	if err != nil {
		return ErrNotFound
	}

	return nil
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, email string) ([]Subscription, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}

	subs, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return subs, nil
}