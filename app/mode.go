package app

// NullReturn used for gen func return if don't have param to return
var NullReturn = []byte{1}

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
