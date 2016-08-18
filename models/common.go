package models

// Predefined model error codes.
const (
	ErrDatabase = -1
	ErrSystem   = -2
	ErrDupRows  = -3
	ErrNotFound = -4
)

// CodeInfo definiton.
type CodeInfo struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

// NewErrorInfo return a CodeInfo represents error.
func NewErrorInfo(info string) *CodeInfo {
	return &CodeInfo{-1, info}
}

// NewNormalInfo return a CodeInfo represents OK.
func NewNormalInfo(info string) *CodeInfo {
	return &CodeInfo{0, info}
}
