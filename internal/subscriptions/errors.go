package subscriptions

import "errors"

var (
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidRepoFormat = errors.New("invalid repo format")
	ErrInvalidRepo      = errors.New("repository is invalid")
	ErrAlreadySubscribed = errors.New("already subscribed")
	ErrEmailConfirmation = errors.New("email confirmation failed")
	ErrSubscriptionConfirmed = errors.New("subscription already confirmed")
	ErrInvalidToken = errors.New("Invalid token")
	ErrNotFound = errors.New("Not found")
)