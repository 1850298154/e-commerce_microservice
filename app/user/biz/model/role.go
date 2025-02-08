package model

import (
	"database/sql/driver"
)

type Role int64

const (
	AdminRole = 1
	UserRole  = 2
)

func (r *Role) String() string {
	return [...]string{"Admin", "User", "Guest"}[*r]
}

func (r *Role) Value() (driver.Value, error) {
	return *r, nil
}

func (r *Role) Scan(value int64) error {
	*r = Role(value)
	return nil
}
