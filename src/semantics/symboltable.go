package semantics

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

type SymbolTable struct {
	symbols map[string]SymbolType
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{}
	st.symbols = make(map[string]SymbolType)
	return st
}

func (s *SymbolTable) Insert(name string, symbolType SymbolType) {
	s.symbols[name] = symbolType
}

func (s *SymbolTable) Get(name string) (SymbolType, bool) {
	x, ok := s.symbols[name]
	return x, ok
}
