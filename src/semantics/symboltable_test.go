package semantics

import "testing"

func TestGetReturnsTrueWhenVariableExistsInTable(t *testing.T) {
	st := NewSymbolTable()

	st.Insert("x", INTEGER)

	symbol, exists := st.Get("x")
	if !exists {
		t.Error("Expected true, got false")
	}

	if symbol.Type() != INTEGER {
		t.Errorf("Expected %s, got %s", INTEGER, symbol.Type())
	}
}

func TestGetReturnsFalseWhenVariableDoesNotExistInTable(t *testing.T) {
	st := NewSymbolTable()
	_, exists := st.Get("foo")

	if exists {
		t.Errorf("Expected false, got true")
	}
}

func TestStringRepresentationsReturnsCorrectTypeNames(t *testing.T) {
	i := INTEGER.String()
	s := STRING.String()
	b := BOOLEAN.String()

	if i != "int" {
		t.Errorf("Expected int, got %s", i)
	}
	if s != "string" {
		t.Errorf("Expected string, got %s", s)
	}
	if b != "bool" {
		t.Errorf("Expected bool, got %s", b)
	}
}

func TestLockSetsVariableAsLocked(t *testing.T) {
	st := NewSymbolTable()
	st.Insert("foo", INTEGER)

	st.Lock("foo")

	foo, exists := st.Get("foo")
	if !exists {
		t.Error("Variable got deleted during locking")
	}

	if !foo.Locked() {
		t.Errorf("Expected foo to be locked")
	}
}

func TestUnLockRemovesLockFromVariable(t *testing.T) {
	st := NewSymbolTable()
	st.Insert("foo", INTEGER)

	st.Lock("foo")
	st.UnLock("foo")

	foo, exists := st.Get("foo")
	if !exists {
		t.Error("Variable got deleted during unlock")
	}

	if foo.Locked() {
		t.Errorf("Expected foo to be unlocked")
	}
}
