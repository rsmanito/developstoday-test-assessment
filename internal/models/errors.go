package models

import "net/http"

type Err struct {
	Msg  string
	Code int
}

func NewError(c int, m string) *Err {
	return &Err{
		Code: c,
		Msg:  m,
	}
}

func (e Err) Error() string {
	return e.Msg
}

var (
	ErrTimeoutExceeded = &Err{"timeout exceeded", http.StatusRequestTimeout}
	ErrNotFound        = &Err{"not found", http.StatusNotFound}
)
