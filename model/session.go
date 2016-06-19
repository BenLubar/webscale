package model // import "github.com/BenLubar/webscale/model"

import (
	"time"

	"github.com/BenLubar/webscale/db"
)

type Session struct {
	ID       SessionID
	User     UserID
	Address  IPs
	Browser  string
	LoggedIn time.Time
	LastSeen time.Time
}

const sessionFields = `s.id, s.user_id, s.address, s.browser, s.logged_in, s.last_seen`

func scanSession(s scanner) (*Session, error) {
	var v Session
	if err := s.Scan(&v.ID, &v.User, &v.Address, &v.Browser, &v.LoggedIn, &v.LastSeen); err != nil {
		return nil, err
	}
	return &v, nil
}

var idGetSession = db.Prepare(`select ` + sessionFields + ` from sessions as s where can_user($1::bigint, 'user-view-sessions', $2::boolean, s.user_id) and s.id = $3::uuid order by s.id asc;`)
