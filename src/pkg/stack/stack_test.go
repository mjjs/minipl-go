package stack

import "testing"

func TestPopReturnsNilForEmptyStack(t *testing.T) {
	s := New()
	actual := s.Pop()

	if actual != nil {
		t.Errorf("Expected nil, got %v", actual)
	}
}

func TestPopReturnsInLIFOOrder(t *testing.T) {
	s := New()

	s.Push(1)
	s.Push(123)
	s.Push(5959)
	x := s.Pop()
	if x != 5959 {
		t.Errorf("Expected %d, got %v", 5959, x)
	}

	x = s.Pop()
	if x != 123 {
		t.Errorf("Expected %d, got %v", 123, x)
	}

	x = s.Pop()
	if x != 1 {
		t.Errorf("Expected %d, got %v", 1, x)
	}
}
