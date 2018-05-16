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

type elementsStack struct {
	s []*element
}

func (s *elementsStack) Push(t *element) {
	s.s = append(s.s, t)
}

func (s *elementsStack) Pop() *element {
	lastElem := len(s.s) - 1

	if lastElem == -1 {
		return nil
	}

	res := s.s[lastElem]
	s.s = s.s[0:lastElem]

	return res
}

func (s *elementsStack) GetLast() *element {
	lastElem := len(s.s) - 1

	if lastElem == -1 {
		return nil
	}

	return s.s[lastElem]
}

func (s *elementsStack) Deep() int {
	return len(s.s)
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
