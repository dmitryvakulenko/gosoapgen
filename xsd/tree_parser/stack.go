package tree_parser

type typesStack struct {
	s []*Type
}

func (s *typesStack) Push(t *Type) {
	s.s = append(s.s, t)
}

func (s *typesStack) Pop() *Type {
	lastElem := len(s.s) - 1

	if lastElem == -1 {
		return nil
	}

	res := s.s[lastElem]
	s.s = s.s[0:lastElem]

	return res
}

func (s *typesStack) GetLast() *Type {
	lastElem := len(s.s) - 1

	if lastElem == -1 {
		return nil
	}

	return s.s[lastElem]
}


type stringsStack struct {
	s []string
}

func (s *stringsStack) Push(t string) {
	s.s = append(s.s, t)
}

func (s *stringsStack) Pop() string {
	lastElem := len(s.s) - 1

	if lastElem == -1 {
		return ""
	}

	res := s.s[lastElem]
	s.s = s.s[0:lastElem]

	return res
}

func (s *stringsStack) GetLast() string {
	lastElem := len(s.s) - 1

	if lastElem == -1 {
		return ""
	}

	return s.s[lastElem]
}
