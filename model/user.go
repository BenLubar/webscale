//go:generate go run mkid.go User

package model // import "github.com/BenLubar/webscale/model"

import (
	"net"
	"time"

	"github.com/BenLubar/webscale/db"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        UserID
	Name      string
	Slug      string
	password  []byte
	Email     string
	JoinDate  time.Time
	LastSeen  pq.NullTime
	Birthday  pq.NullTime
	Signature string
	Bio       string
	Location  string
	Website   string
	Avatar    string
}

const userFields = `u.id, u.name, u.slug, u.password, u.email, u.join_date, u.last_seen, u.birthday, u.signature, u.bio, u.location, u.website, u.avatar`

func scanUser(s scanner) (*User, error) {
	var u User
	if err := s.Scan(&u.ID, &u.Name, &u.Slug, &u.password, &u.Email, &u.JoinDate, &u.LastSeen, &u.Birthday, &u.Signature, &u.Bio, &u.Location, &u.Website, &u.Avatar); err != nil {
		return nil, err
	}
	return &u, nil
}

var idGetUser = db.Prepare(`select ` + userFields + ` from users as u where can_user($1::bigint, 'user-meta', $2::boolean, u.id) and u.id = $3::bigint order by u.id asc;`)
var idsGetUser = db.Prepare(`select ` + userFields + ` from users as u where can_user($1::bigint, 'user-meta', $2::boolean, u.id) and array[u.id] <@ $3::bigint[] order by u.id asc;`)

const PasswordCost = bcrypt.DefaultCost

var userSetPassword = db.Prepare(`update users set password = $1::varchar where id = $2::bigint and (password = $3::varchar or ($3::varchar is null and password is null));`)

func (u *User) SetPassword(tx *db.Tx, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return errors.Wrap(err, "hash user password")
	}

	rows, err := tx.Exec(userSetPassword, hash, u.ID, u.password)
	if err != nil {
		return errors.Wrap(err, "update user password")
	}

	if rows == 0 {
		return errors.New("password change failed: password in database did not match")
	}

	u.password = hash

	return nil
}

func (u *User) CheckPassword(tx *db.Tx, password string) error {
	if u.password == nil {
		return errors.New("user has no password")
	}

	if err := bcrypt.CompareHashAndPassword(u.password, []byte(password)); err != nil {
		return errors.Wrap(err, "password did not match")
	}

	if cost, err := bcrypt.Cost(u.password); err != nil {
		return errors.Wrap(err, "check password cost")
	} else if cost < PasswordCost {
		return errors.Wrap(u.SetPassword(tx, password), "re-hashing password with higher cost")
	}

	return nil
}

type UserIP struct {
	User     UserID
	IP       net.IP
	LastSeen time.Time
}

var userIPs = db.Prepare(`select ip, last_seen from user_ips where user_id = $1 order by last_seen desc;`)

func (u *User) IPs(ctx *Context) ([]UserIP, error) {
	rows, err := ctx.Tx.Query(userIPs, u.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "find IP addresses of user %d", u.ID)
	}
	defer rows.Close()

	var ips []UserIP

	for rows.Next() {
		var s string
		var t time.Time
		if err = rows.Scan(&s, &t); err != nil {
			return nil, errors.Wrapf(err, "scan IP addresses of user %d", u.ID)
		}

		ip := net.ParseIP(s)
		if ip == nil {
			return nil, errors.Errorf("invalid IP address: %q", s)
		}

		ips = append(ips, UserIP{User: u.ID, IP: ip, LastSeen: t})
	}

	return ips, errors.Wrapf(rows.Close(), "find IP addresses of user %d", u.ID)
}

var usersByIP = db.Prepare(`select user_id, ip, last_seen from user_ips where ip <<= $1::cidr order by last_seen desc;`)

func UsersByIP(ctx *Context, cidr *net.IPNet) ([]UserIP, error) {
	if _, bits := cidr.Mask.Size(); bits == 0 {
		return nil, errors.Errorf("invalid CIDR mask: %q", cidr)
	}

	rows, err := ctx.Tx.Query(usersByIP, cidr.String())
	if err != nil {
		return nil, errors.Wrapf(err, "find users with IP address %q", cidr)
	}
	defer rows.Close()

	var ips []UserIP

	for rows.Next() {
		var u UserID
		var s string
		var t time.Time
		if err = rows.Scan(&u, &s, &t); err != nil {
			return nil, errors.Wrapf(err, "scan users with IP address %q", cidr)
		}

		ip := net.ParseIP(s)
		if ip == nil {
			return nil, errors.Errorf("invalid IP address: %q", s)
		}

		ips = append(ips, UserIP{User: u, IP: ip, LastSeen: t})
	}

	return ips, errors.Wrapf(rows.Close(), "find users with IP address %q", cidr)
}
