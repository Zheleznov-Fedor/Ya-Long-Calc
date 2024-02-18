package utils

type Queue []string

func (q *Queue) Put(n string) {
	*q = append(*q, n)
}

func (q *Queue) Get() string {
	if len(*q) == 0 {
		return ""
	}
	element := (*q)[0]
	*q = (*q)[1:]
	return element
}

func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}

type Stack []string

func (s *Stack) Push(n string) {
	*s = append(*s, n)
}

func (s *Stack) Pop() string {
	if len(*s) == 0 {
		return ""
	}
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Head() string {
	if len(*s) == 0 {
		return ""
	}
	index := len(*s) - 1
	element := (*s)[index]
	return element
}
