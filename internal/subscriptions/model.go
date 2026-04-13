package subscriptions

type Subscription struct {
	ID               string
	Email            string
	Repo             string
	Confirmed        bool
	ConfirmToken     string
	UnsubscribeToken string
	LastSeenTag      string
}