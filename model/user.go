//go:generate go run mkid.go User

package model // import "github.com/BenLubar/webscale/model"

import (
	"time"

	"github.com/BenLubar/webscale/db"
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
	LastSeen  time.Time
	Address   IPs
	Birthday  time.Time
	Signature string
	Bio       string
	Location  string
	Website   string
	Avatar    string
}

const PasswordCost = bcrypt.DefaultCost

var userSetPassword = db.Prepare(`update users set password = $1::varchar where id = $2::bigint and (password = $3::varchar or ($3::varchar is null and password is null));`)

func (u *User) SetPassword(tx *db.Tx, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return errors.Wrap(err, "hash user password")
	}

	result, err := tx.Exec(userSetPassword, hash, u.ID, u.password)
	if err != nil {
		return errors.Wrap(err, "update user password")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "check update user password")
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
