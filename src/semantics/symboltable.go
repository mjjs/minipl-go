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

type Symbol struct {
	symbolType SymbolType
	locked     bool
}

func (s Symbol) Type() SymbolType { return s.symbolType }
func (s Symbol) Locked() bool     { return s.locked }

type SymbolTable struct {
	symbols map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{}
	st.symbols = make(map[string]Symbol)

	return st
}

func (s *SymbolTable) Insert(name string, symbolType SymbolType) {
	s.symbols[name] = Symbol{symbolType: symbolType}
}

func (s *SymbolTable) Get(name string) (Symbol, bool) {
	x, ok := s.symbols[name]
	return x, ok
}

func (s *SymbolTable) Lock(name string) bool {
	symbol, ok := s.symbols[name]
	if !ok {
		return false
	}

	symbol.locked = true
	s.symbols[name] = symbol

	return true
}

func (s *SymbolTable) UnLock(name string) bool {
	symbol, ok := s.symbols[name]
	if !ok {
		return false
	}

	symbol.locked = false
	s.symbols[name] = symbol

	return true
}
