package domain

import "errors"

type Error struct {
	Code int
	Msg  string
	Err  error
}

func errorNew(code int, msg string) Error {
	return Error{
		Code: code,
		Msg:  msg,
		Err:  errors.New(msg),
	}
}

var (
	// Config
	ErrLoadConfig = errorNew(1000, "Failed to load config")

	// Member
	ErrMemberRegisterFailDuplicateEmail = errorNew(1003, "Register member failed: duplicate email")
	ErrMemberRegisterFailDuplicateHash  = errorNew(1004, "Register member failed: duplicate hash")
	ErrMemberEmailNotFound              = errorNew(1005, "Email not found")
	ErrCreateMemberFail                 = errorNew(1002, "Failed to create member")

	// Token
	ErrGenerateToken    = errorNew(1006, "Failed to generate token")
	ErrValidateToken    = errorNew(1007, "Failed to validate token")
	ErrValidateHashAuth = errorNew(1008, "Failed to validate hash auth")
	ErrTokenExpired     = errorNew(1009, "Token has expired")
)
