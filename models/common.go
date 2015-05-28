package models

const (
	ErrDatabase = -1
	ErrSystem   = -2
	ErrDupRows  = -3
	ErrNotFound = -4
)

type CodeInfo struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

func NewErrorInfo(info string) *CodeInfo {
	return &CodeInfo{-1, info}
}

func NewNormalInfo(info string) *CodeInfo {
	return &CodeInfo{0, info}
}
