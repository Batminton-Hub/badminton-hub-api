package domain

import "errors"

var (
	// Member
	ErrMemberRegisterFailDuplicateEmail = errors.New("register member failed: duplicate email")
	ErrMemberRegisterFailDuplicateHash  = errors.New("register member failed: duplicate hash")
	ErrMemberEmailNotFound              = errors.New("email not found")
)
