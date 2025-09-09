package domain

import (
	"fmt"
	"net/http"
)

type Resp struct {
	HttpStatus int
	Status     string
	Code       int
	Msg        string
	Err        error
}

func errorNew(code int, msg string, httpStatus int) Resp {
	return Resp{
		HttpStatus: httpStatus,
		Status:     ERROR,
		Code:       code,
		Msg:        msg,
		Err:        fmt.Errorf("%s", msg),
	}
}

func successNew(code int, msg string, httpStatus int) Resp {
	return Resp{
		HttpStatus: httpStatus,
		Status:     SUCCESS,
		Code:       code,
		Msg:        msg,
		Err:        nil,
	}
}

var (
	//////////////////////// Success Code ////////////////////////
	Success             = successNew(0, "Success", http.StatusOK)
	RegisterSuccess     = successNew(0, "Success", http.StatusCreated)
	LoginSuccess        = successNew(0, "Success", http.StatusOK)
	AuthSuccess         = successNew(0, "Success", http.StatusOK)
	UpdateMemberSuccess = successNew(0, "Success", http.StatusOK)
	RedirectSuccess     = successNew(0, "Success", http.StatusTemporaryRedirect)

	//////////////////////// Error Code ////////////////////////
	// Config
	ErrLoadConfig = errorNew(1000, "Failed to load config", http.StatusInternalServerError) // ไม่สามารถโหลดค่า config ได้

	// Request
	ErrInvalidInput       = errorNew(1001, "Invalid input", http.StatusBadRequest)        // ข้อมูลไม่ถูกต้อง
	ErrPlatformNotSupport = errorNew(1002, "Platform not support", http.StatusBadRequest) // ระบบไม่รองรับ Platform นี้

	// Member
	ErrMemberRegisterFailDuplicateEmail = errorNew(2000, "Register member failed: duplicate email", http.StatusBadRequest)
	ErrMemberRegisterFailDuplicateHash  = errorNew(2001, "Register member failed: duplicate hash", http.StatusBadRequest)
	ErrMemberEmailNotFound              = errorNew(2002, "Email not found", http.StatusBadRequest)
	ErrCreateMemberFail                 = errorNew(2003, "Failed to create member", http.StatusBadRequest)
	ErrLoginHashPassword                = errorNew(2004, "Failed to hash password", http.StatusBadRequest)
	ErrGetMember                        = errorNew(2005, "Failed to get member", http.StatusBadRequest)
	ErrUpdateMemberFail                 = errorNew(2006, "Failed to update member", http.StatusBadRequest)

	// OAuth
	ErrInvalidOAuthState    = errorNew(3000, "Invalid OAuth state", http.StatusBadRequest)
	ErrSetGoogleState       = errorNew(3001, "Failed to set Google OAuth state", http.StatusBadRequest)
	ErrInvalidOAuthExchange = errorNew(3002, "Invalid OAuth exchange", http.StatusBadRequest)
	ErrInvalidOAuthClient   = errorNew(3003, "Invalid OAuth client", http.StatusBadRequest)
	ErrInvalidOAuthDecode   = errorNew(3004, "Invalid OAuth decode", http.StatusBadRequest)

	// Token
	ErrGenerateToken    = errorNew(4000, "Failed to generate token", http.StatusBadRequest)
	ErrValidateToken    = errorNew(4001, "Failed to validate token", http.StatusBadRequest)
	ErrValidateHashAuth = errorNew(4002, "Failed to validate hash auth", http.StatusBadRequest)
	ErrTokenExpired     = errorNew(4003, "Token has expired", http.StatusBadRequest)

	// General
	ErrActionNotSupport = errorNew(5000, "Action not support", http.StatusInternalServerError)
)

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
