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

func ParseMode(optionMap map[string]interface{}) (mode ServiceAppModeTyp) {
	modeTmp, ok := optionMap["mode"]
	if ok {
		mode = modeTmp.(ServiceAppModeTyp)
	} else {
		mode = ServiceAppMode
	}
	return
}
