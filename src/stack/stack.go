package stack

import "container/list"

type Stack struct {
	list *list.List
}

func New() *Stack {
	return &Stack{list: list.New()}
}

func (s *Stack) Push(value interface{}) {
	s.list.PushBack(value)
}

func (s *Stack) Pop() interface{} {
	if s.list.Len() == 0 {
		return nil
	}

	tail := s.list.Back()
	s.list.Remove(tail)
	return tail.Value
}
