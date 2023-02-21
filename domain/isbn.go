package domain

import "strings"

type ISBN string

func (i ISBN) IsValid() bool {
	return len(string(i)) == 13
}

func (i ISBN) IsEqual(v ISBN) bool {
	return strings.Compare(string(i), string(v)) == 0
}
