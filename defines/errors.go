package defines

import "fmt"

var (
	UnknownError                         = NewCexError("UNKNOWN")
	CexErrorInvalidAPIKey                = NewCexError("INVALID_API_KEY")
	MarketOrderErrorNotEnoughBalance     = NewCexError("NOT_ENOUGH_BALANCE")
	MarketOrderErrorFilterFailureLotSize = NewCexError("FILTER_FAILURE_LOT_SIZE")
)

type CexError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Raw     interface{} `json:"raw"`
}

func (e CexError) Error() string {
	return fmt.Sprintf("(Code %s): %s", e.Code, e.Message)
}

func (e *CexError) SetRawError(err error) {
	e.Message = err.Error()
	e.Raw = err
}

func NewCexError(code string) *CexError {
	return &CexError{
		Code: code,
	}
}

func NewCexErrorWithMsg(code string, msg string) *CexError {
	return &CexError{
		Code:    code,
		Message: msg,
	}
}

func IsCexError(err error) bool {
	_, ok := err.(*CexError)
	return ok
}

func IsError(err error, cexError *CexError) bool {
	if !IsCexError(err) {
		return false
	}

	return cexError.Code == err.(*CexError).Code
}
