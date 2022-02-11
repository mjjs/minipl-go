package symboltable

type SymbolType int

const (
	INTEGER SymbolType = iota
	STRING
	BOOLEAN
)

func (s SymbolType) String() string {
	switch s {
	case INTEGER:
		return "int"
	case STRING:
		return "string"
	default:
		return "bool"
	}
}

type Symbol struct {
	symbolType SymbolType
}

func (s Symbol) Type() SymbolType { return s.symbolType }
