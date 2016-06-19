package model // import "github.com/BenLubar/webscale/model"

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type scanner interface {
	Scan(dest ...interface{}) error
}

// ID provides support for the PostgreSQL type bigint, with 0 standing for null.
type ID int64

// Scan implements sql.Scanner.
func (id *ID) Scan(value interface{}) error {
	if value == nil {
		*id = 0
		return nil
	}

	if i, ok := value.(int64); ok {
		if i == 0 {
			return errors.New("ID cannot be 0")
		}
		*id = ID(i)
		return nil
	}

	return errors.Errorf("unexpected ID type %T", value)
}

// Value implements driver.Valuer.
func (id ID) Value() (driver.Value, error) {
	if id == 0 {
		return nil, nil
	}
	return int64(id), nil
}

// IDs provides support for the PostgreSQL type bigint[].
type IDs []ID

// Scan implements sql.Scanner.
func (ids *IDs) Scan(value interface{}) error {
	var src string
	if b, ok := value.([]byte); ok {
		src = string(b)
	} else {
		return errors.Errorf("unexpected IDs type %T", value)
	}

	if len(src) < 2 || src[0] != '{' || src[len(src)-1] != '}' {
		return errors.Errorf("invalid array: %q", src)
	}

	s := strings.Split(src[1:len(src)-1], ",")
	decoded := make(IDs, len(s))

	for i, id := range s {
		if id == "NULL" {
			continue
		}

		n, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "invalid ID: %q", id)
		}

		if n == 0 {
			return errors.New("ID cannot be 0")
		}

		decoded[i] = ID(n)
	}

	*ids = decoded
	return nil
}

// Value implements driver.Valuer.
func (ids IDs) Value() (driver.Value, error) {
	s := make([]string, len(ids))
	for i, id := range ids {
		if id == 0 {
			s[i] = "NULL"
		}

		s[i] = strconv.FormatInt(int64(id), 10)
	}

	return "{" + strings.Join(s, ",") + "}", nil
}

// UUID provides support for the PostgreSQL type uuid.
type UUID [16]byte

// Scan implements sql.Scanner.
func (uuid *UUID) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("unexpected UUID type %T", value)
	}

	if len(b) != 36 {
		return errors.Errorf("unexpected UUID length %d %q", len(b), b)
	}

	// 1092ff5e-26d5-4da7-9dff-214de62f42ae
	// 000000000011111111112222222222333333
	// 012345678901234567890123456789012345

	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return errors.Errorf("unexpected UUID format: %q", b)
	}

	var decoded UUID
	if i, err := hex.Decode(decoded[0:4], b[0:8]); err != nil {
		return errors.Wrapf(err, "at position %d of UUID %q", 0+i, b)
	}
	if i, err := hex.Decode(decoded[4:6], b[9:13]); err != nil {
		return errors.Wrapf(err, "at position %d of UUID %q", 9+i, b)
	}
	if i, err := hex.Decode(decoded[6:8], b[14:18]); err != nil {
		return errors.Wrapf(err, "at position %d of UUID %q", 14+i, b)
	}
	if i, err := hex.Decode(decoded[8:10], b[19:23]); err != nil {
		return errors.Wrapf(err, "at position %d of UUID %q", 19+i, b)
	}
	if i, err := hex.Decode(decoded[10:16], b[24:36]); err != nil {
		return errors.Wrapf(err, "at position %d of UUID %q", 24+i, b)
	}

	*uuid = decoded
	return nil
}

// Value implements driver.Valuer.
func (uuid UUID) Value() (driver.Value, error) {
	return []byte(uuid.String()), nil
}

// String implements fmt.Stringer.
func (uuid UUID) String() string {
	return fmt.Sprintf("%02x-%02x-%02x-%02x-%02x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// IPs provides support for the PostgreSQL type inet[].
type IPs []net.IP

// Scan implements sql.Scanner.
func (ips *IPs) Scan(value interface{}) error {
	var src string
	if b, ok := value.([]byte); ok {
		src = string(b)
	} else {
		return errors.Errorf("unexpected IPs type %T", value)
	}

	if len(src) < 2 || src[0] != '{' || src[len(src)-1] != '}' {
		return errors.Errorf("invalid array: %q", src)
	}

	s := strings.Split(src[1:len(src)-1], ",")
	addrs := make(IPs, len(s))

	for i, ip := range s {
		addr := net.ParseIP(ip)
		if addr == nil {
			return errors.Errorf("invalid IP: %q", ip)
		}

		addrs[i] = addr
	}

	*ips = addrs
	return nil
}

// Value implements driver.Valuer.
func (ips IPs) Value() (driver.Value, error) {
	s := make([]string, len(ips))
	for i, ip := range ips {
		s[i] = ip.String()
	}

	return "{" + strings.Join(s, ",") + "}", nil
}

// Strings provides support for the PostgreSQL types text[] and citext[].
type Strings []string

// Scan implements sql.Scanner.
func (strs *Strings) Scan(value interface{}) error {
	var src string
	if b, ok := value.([]byte); ok {
		src = string(b)
	} else {
		return errors.Errorf("unexpected IDs type %T", value)
	}

	if len(src) < 2 || src[0] != '{' || src[len(src)-1] != '}' {
		return errors.Errorf("invalid array: %q", src)
	}

	if len(src) == 2 {
		*strs = nil
		return nil
	}

	if len(src) < 4 || src[1] != '"' || src[len(src)-2] != '"' {
		return errors.Errorf("invalid string array: %q", src)
	}

	values := strings.Split(src[2:len(src)-2], `","`)
	decoded := make(Strings, len(values))

	for i, v := range values {
		s, err := strconv.Unquote(`"` + v + `"`)
		if err != nil {
			return errors.Wrapf(err, "invalid string: %q", v)
		}

		decoded[i] = s
	}

	*strs = decoded
	return nil
}

// Value implements driver.Valuer.
func (strs Strings) Value() (driver.Value, error) {
	s := make([]string, len(strs))
	for i, v := range strs {
		s[i] = strconv.Quote(v)
	}

	return "{" + strings.Join(s, ",") + "}", nil
}

// Use $1 and $2 without making the query actually do anything.
const noopPermission = `(coalesce(0, $1::bigint) = 0 and coalesce(false, $2::boolean) = false)`
