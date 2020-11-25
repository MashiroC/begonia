package berr

import (
	"fmt"
)

type BErr struct {
	Package string
	Typ     string
	Msg     string
}

func (err BErr) Error() string {
	// fmt.Sprintf("%s %s error: %s")
	// example: dispatch link error: you are link in a black hole

	return err.Package + " " + err.Typ + " error: " + err.Msg
}

func New(pkg, typ, msg string) error {
	if msg == "" {
		return nil
	}
	return BErr{
		Package: pkg,
		Typ:     typ,
		Msg:     msg,
	}
}

func NewF(pkg, typ, msgF string, param ...interface{}) error {
	return BErr{
		Package: pkg,
		Typ:     typ,
		Msg:     fmt.Sprintf(msgF, param...),
	}
}

func Warp(pkg, typ string, err error) error {
	return BErr{
		Package: pkg,
		Typ:     typ,
		Msg:     err.Error(),
	}
}
