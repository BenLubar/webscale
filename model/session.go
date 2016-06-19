package model // import "github.com/BenLubar/webscale/model"

import "time"

// SessionID is the ID of a Session.
type SessionID struct{ UUID }

type Session struct {
	ID       SessionID
	User     UserID
	Address  IPs
	Browser  string
	LoggedIn time.Time
	LastSeen time.Time
}
