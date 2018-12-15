package pathfinder

import (
	"fmt"
)

type Error interface {
	error
	Timeout() bool
	MaxPayloadSize() bool
	RegexMatchFailed() bool
	DecodeFailed() bool
}

type MaxPayloadSizeError struct {
	cur int
	max int
}

func (e MaxPayloadSizeError) Error() string {
	return fmt.Sprintf("maximum payload size exceeded, is %d, should be <= (%d - %d) [maxPayloadSize - minimumAvailablePayloadSize]", e.cur, e.max, minimumAvailablePayloadSize)
}

func (e MaxPayloadSizeError) Timeout() bool {
	return false
}

func (e MaxPayloadSizeError) MaxPayloadSize() bool {
	return true
}

func (e MaxPayloadSizeError) RegexMatchFailed() bool {
	return false
}

func (e MaxPayloadSizeError) DecodeFailed() bool {
	return false
}

func NewMaxPayloadSizeError(cur, max int) Error {
	var e MaxPayloadSizeError
	e.cur = cur
	e.max = max

	return e
}

type RegexMatchFailed struct {
}

func (e RegexMatchFailed) Error() string {
	return "regex match failed"
}

func (e RegexMatchFailed) Timeout() bool {
	return false
}

func (e RegexMatchFailed) MaxPayloadSize() bool {
	return true
}

func (e RegexMatchFailed) RegexMatchFailed() bool {
	return false
}

func (e RegexMatchFailed) DecodeFailed() bool {
	return false
}

func NewRegexMatchFailed() Error {
	var e RegexMatchFailed

	return e
}

type DecodeFailed struct {
}

func (e DecodeFailed) Error() string {
	return "regex match failed"
}

func (e DecodeFailed) Timeout() bool {
	return false
}

func (e DecodeFailed) MaxPayloadSize() bool {
	return true
}

func (e DecodeFailed) RegexMatchFailed() bool {
	return false
}

func (e DecodeFailed) DecodeFailed() bool {
	return false
}

func NewDecodeFailed() Error {
	var e DecodeFailed

	return e
}
