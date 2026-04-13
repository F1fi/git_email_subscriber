package scanner

import (
	"context"
	"log"
	"time"
	"git_email_subscriber/internal/subscriptions"
	"git_email_subscriber/internal/github_api"
	"git_email_subscriber/internal/email_service"
)

type IScanner interface{
	Start(ctx context.Context)
}

type Scanner struct {
	repo     subscriptions.ISubscriptionRepository
	github   github_api.IGitHubApi
	notifier email_service.IEmailService

	interval time.Duration
}

func NewScanner(
	repo subscriptions.ISubscriptionRepository,
	github github_api.IGitHubApi,
	notifier email_service.IEmailService,
	interval time.Duration,
) *Scanner {
	return &Scanner{
		repo,
		github,
		notifier,
		interval,
	}
}

func (s *Scanner) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.run(ctx)
			case <-ctx.Done():
				log.Println("scanner stopped")
				return
			}
		}
	}()
}

func (s *Scanner) run(ctx context.Context) {
	subs, err := s.repo.GetActive(ctx)
	if err != nil {
		log.Println("failed to fetch subscriptions:", err)
		return
	}

	for _, sub := range subs {
		s.processSubscription(ctx, sub)
	}
}

func (s *Scanner) processSubscription(ctx context.Context, sub subscriptions.Subscription) {
	tag, err := s.github.GetLatestRelease(ctx, sub.Repo)
	if err != nil {
		log.Println("github error:", err)
		return
	}

	if sub.LastSeenTag == "" {
		_ = s.repo.UpdateLastSeenTag(ctx, sub.ID, tag)
		return
	}

	if tag != sub.LastSeenTag {
		err := s.sendReleaseNote(sub.Email, sub.Repo, tag)
		if err != nil {
			log.Println("email error:", err)
			return
		}

		_ = s.repo.UpdateLastSeenTag(ctx, sub.ID, tag)
	}
}

func (s *Scanner) sendReleaseNote(email, repo, tag string) error {
	///TODO: Add text
	err := s.notifier.Send(email, "repo, tag")
	return err
}