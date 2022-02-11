package symboltable

type SymbolTable struct {
	symbols map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{}
	st.symbols = make(map[string]Symbol)

	return st
}

func (s *SymbolTable) Insert(name string, symbolType SymbolType) *SymbolTable {
	s.symbols[name] = Symbol{symbolType: symbolType}
	return s
}

func (s *SymbolTable) Get(name string) (Symbol, bool) {
	x, ok := s.symbols[name]
	return x, ok
}
