package domain

import "fmt"

type ErrorResp struct {
	Code int
	Msg  string
	Err  error
}

type SuccessResp struct {
	Code int
	Msg  string
}

func errorNew(code int, msg string) ErrorResp {
	return ErrorResp{
		Code: code,
		Msg:  msg,
		Err:  fmt.Errorf("%s", msg),
	}
}

func successNew(code int, msg string) SuccessResp {
	return SuccessResp{
		Code: code,
		Msg:  msg,
	}
}

var (
	//////////////////////// Success Code ////////////////////////
	Success         = successNew(0, "Success")
	RegisterSuccess = successNew(0, "Success")
	LoginSuccess    = successNew(0, "Success")
	AuthSuccess     = successNew(0, "Success")

	//////////////////////// Error Code ////////////////////////
	// Config
	ErrLoadConfig = errorNew(1000, "Failed to load config")

	// Request
	ErrInvalidInput = errorNew(1001, "Invalid input")

	// Member
	ErrMemberRegisterFailDuplicateEmail = errorNew(2000, "Register member failed: duplicate email")
	ErrMemberRegisterFailDuplicateHash  = errorNew(2001, "Register member failed: duplicate hash")
	ErrMemberEmailNotFound              = errorNew(2002, "Email not found")
	ErrCreateMemberFail                 = errorNew(2003, "Failed to create member")
	ErrLoginHashPassword                = errorNew(2004, "Failed to hash password")

	// OAuth
	ErrInvalidOAuthState    = errorNew(3000, "Invalid OAuth state")
	ErrSetGoogleState       = errorNew(3001, "Failed to set Google OAuth state")
	ErrInvalidOAuthExchange = errorNew(3002, "Invalid OAuth exchange")
	ErrInvalidOAuthClient   = errorNew(3003, "Invalid OAuth client")
	ErrInvalidOAuthDecode   = errorNew(3004, "Invalid OAuth decode")

	// Token
	ErrGenerateToken    = errorNew(4000, "Failed to generate token")
	ErrValidateToken    = errorNew(4001, "Failed to validate token")
	ErrValidateHashAuth = errorNew(4002, "Failed to validate hash auth")
	ErrTokenExpired     = errorNew(4003, "Token has expired")
)
