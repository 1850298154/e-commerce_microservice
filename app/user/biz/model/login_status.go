package model

import (
	"database/sql/driver"
)

type LoginStatus int8

const (
	LoginFailed  = 0
	LoginSuccess = 1
)

func (l *LoginStatus) String() string {
	return [...]string{"LoginFailed", "LoginSuccess"}[*l]
}

func (l *LoginStatus) Value() (driver.Value, error) {
	return *l, nil
}

func (l *LoginStatus) Scan(value int64) error {
	*l = LoginStatus(value)
	return nil
}
