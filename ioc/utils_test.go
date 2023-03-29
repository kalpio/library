package ioc_test

import (
	"sync"
)

/**
 * Tests interfaces
 */
type iFirstInterface interface {
	SetText(s string)
	GetText() string
	GetTextFromSecond() string
}

type iSecondInterface interface {
	GetSecondText() string
	SetSecondText(s string)
}

/**
 * Test implementation
 */
type iFirstImpl struct {
	mu     sync.RWMutex
	Second iSecondInterface
	text   string
}

type secondImpl struct {
	mu   sync.Mutex
	text string
}

func newFirstImpl() *iFirstImpl {
	return new(iFirstImpl)
}

func newFirstImplWithSecondInterface(secondInterface iSecondInterface) *iFirstImpl {
	return &iFirstImpl{Second: secondInterface}
}

func (f *iFirstImpl) SetText(s string) {
	f.text = s
}

func (f *iFirstImpl) GetText() string {
	return f.text
}

func (f *iFirstImpl) GetTextFromSecond() string {
	return f.Second.GetSecondText()
}

func newSecondImpl() *secondImpl {
	return new(secondImpl)
}

func newSecondImplNonPointer() secondImpl {
	return secondImpl{}
}

func (second *secondImpl) SetSecondText(s string) {
	second.text = s
}

func (second *secondImpl) GetSecondText() string {
	return second.text
}
