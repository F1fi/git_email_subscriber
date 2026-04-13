package subscriptions

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"errors"
	"github.com/jackc/pgx/v5"
)

type ISubscriptionRepository interface {
	Create(ctx context.Context, sub Subscription) error
	Exists(ctx context.Context, email, repo string) (bool, error)

	FindByConfirmToken(ctx context.Context, token string) (*Subscription, error)
	Confirm(ctx context.Context, token string) error

	DeleteByUnsubscribeToken(ctx context.Context, token string) error

	GetByEmail(ctx context.Context, email string) ([]Subscription, error)

	GetActive(ctx context.Context) ([]Subscription, error)
	UpdateLastSeenTag(ctx context.Context, id string, tag string) error
}

type SubscriptionRepository struct{
	db *pgxpool.Pool
}

func CreateRepo(db *pgxpool.Pool) ISubscriptionRepository{
	return &SubscriptionRepository{db}
}

func (r *SubscriptionRepository) Exists(ctx context.Context, repo, email string) (bool, error) {
	return false, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub Subscription) error {
	query := `
		INSERT INTO subscriptions (id, email, repo, confirmed, confirm_token, unsubscribe_token)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		sub.ID,
		sub.Email,
		sub.Repo,
		sub.Confirmed,
		sub.ConfirmToken,
		sub.UnsubscribeToken,
	)

	return err
}

func (r *SubscriptionRepository) Confirm(ctx context.Context, token string) error {
	query := `
		UPDATE subscriptions
		SET confirmed = true
		WHERE confirm_token = $1
	`

	_, err := r.db.Exec(ctx, query, token)
	return err
}

func (r *SubscriptionRepository) FindByConfirmToken(
	ctx context.Context,
	token string,
) (*Subscription, error) {

	query := `
		SELECT id, email, repo, confirmed, confirm_token, unsubscribe_token, last_seen_tag
		FROM subscriptions
		WHERE confirm_token = $1
	`

	row := r.db.QueryRow(ctx, query, token)

	var sub Subscription

	err := row.Scan(
		&sub.ID,
		&sub.Email,
		&sub.Repo,
		&sub.Confirmed,
		&sub.ConfirmToken,
		&sub.UnsubscribeToken,
		&sub.LastSeenTag,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &sub, nil
}

func (r *SubscriptionRepository) DeleteByUnsubscribeToken(ctx context.Context, token string) error {
	query := `
		DELETE FROM subscriptions
		WHERE unsubscribe_token = $1
	`

	_, err := r.db.Exec(ctx, query, token)
	return err
}

func (r *SubscriptionRepository) GetByEmail(ctx context.Context, email string) ([]Subscription, error) {
	query := `
		SELECT id, email, repo, confirmed, confirm_token, unsubscribe_token, last_seen_tag
		FROM subscriptions
		WHERE email=$1 AND confirmed=true
	`

	rows, err := r.db.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Subscription

	for rows.Next() {
		var sub Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.Email,
			&sub.Repo,
			&sub.Confirmed,
			&sub.ConfirmToken,
			&sub.UnsubscribeToken,
			&sub.LastSeenTag,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, sub)
	}

	return result, nil
}

func (r *SubscriptionRepository) GetActive(ctx context.Context) ([]Subscription, error) {
	query := `
		SELECT id, email, repo, confirmed, confirm_token, unsubscribe_token, last_seen_tag
		FROM subscriptions
		WHERE confirmed = true
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Subscription

	for rows.Next() {
		var sub Subscription

		err := rows.Scan(
			&sub.ID,
			&sub.Email,
			&sub.Repo,
			&sub.Confirmed,
			&sub.ConfirmToken,
			&sub.UnsubscribeToken,
			&sub.LastSeenTag,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, sub)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (r *SubscriptionRepository) UpdateLastSeenTag(ctx context.Context, id string, tag string) error {
	query := `
		UPDATE subscriptions
		SET last_seen_tag = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, tag, id)
	return err
}