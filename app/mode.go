package app

var ServiceAppMode = Reflect

type ServiceAppModeTyp int

const (
	invalid ServiceAppModeTyp = iota
	Ast
	Reflect
)

func (s ServiceAppModeTyp) String() string {
	switch s {
	case Ast:
		return "ast"
	case Reflect:
		return "reflect"
	}
	return ""
}
