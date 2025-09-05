package service

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"
	core_util "Badminton-Hub/internal/util"
	"Badminton-Hub/util"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	REGISTER = "REGISTER"
	LOGIN    = "LOGIN"

	// Platform
	GOOGLE = "GOOGLE"

	// Status
	PENDING = "PENDING"
	ACTIVE  = "ACTIVE"
	BANNED  = "BANNED"
	DELETED = "DELETED"

	// Type Member
	MEMBER = "MEMBER"
)

type MemberUtil struct {
	memberRepo     port.MemberRepo
	middlewareUtil port.MiddlewareUtil
}

func NewMemberUtil(memberRepo port.MemberRepo, middlewareUtil port.MiddlewareUtil) *MemberUtil {
	memberUtil := &MemberUtil{
		memberRepo:     memberRepo,
		middlewareUtil: middlewareUtil,
	}
	return memberUtil
}

func (m *MemberUtil) RegisterMember(registerForm domain.RegisterForm) (int, domain.ResponseRegisterMember) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.ResponseRegisterMember{}
	config := util.LoadConfig()

	createAt := time.Now()
	updateAt := time.Now()
	memberBody := domain.Member{
		UserID:     core_util.GenUserID(registerForm.Email, registerForm.Password),
		Email:      registerForm.Email,
		Password:   core_util.HashPassword(registerForm.Password, config.KeyHashPassword),
		Gender:     registerForm.Gender,
		Hash:       core_util.GenerateHash(config.KeyHashMember),
		Status:     ACTIVE,
		TypeMember: MEMBER,
		CreatedAt:  createAt,
		UpdatedAt:  updateAt,
	}

	if err := m.memberRepo.SaveMember(ctx, memberBody); err != nil {
		switch err {
		case domain.ErrMemberRegisterFailDuplicateEmail.Err:
			response.Code = domain.ErrMemberRegisterFailDuplicateEmail.Code
			response.Message = domain.ErrMemberRegisterFailDuplicateEmail.Msg
		case domain.ErrMemberRegisterFailDuplicateHash.Err:
			response.Code = domain.ErrMemberRegisterFailDuplicateHash.Code
			response.Message = domain.ErrMemberRegisterFailDuplicateHash.Msg
		default:
			response.Code = domain.ErrCreateMemberFail.Code
			response.Message = domain.ErrCreateMemberFail.Msg
		}
		return http.StatusInternalServerError, response // 500 Internal Server Error for other DB errors
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	token, err := core_util.GenBearerToken(hashAuth, m.middlewareUtil.Encryption())
	if err != nil {
		response.Code = domain.ErrGenerateToken.Code
		response.Message = domain.ErrGenerateToken.Msg
		return http.StatusInternalServerError, response
	}

	response.Code = domain.RegisterSuccess.Code
	response.Message = domain.RegisterSuccess.Msg
	response.BearerToken = token
	return http.StatusCreated, response // 201 Created for successful registration
}

func (m *MemberUtil) Login(loginForm domain.LoginForm) (int, domain.ResponseLogin) {
	response := domain.ResponseLogin{}
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	config := util.LoadConfig()

	memberBody, err := m.memberRepo.FindEmailMember(ctx, loginForm.Email)
	if err != nil {
		response.Code = domain.ErrMemberEmailNotFound.Code
		response.Message = domain.ErrMemberEmailNotFound.Msg
		return http.StatusBadRequest, response
	}

	// Check password
	if memberBody.Password != core_util.HashPassword(loginForm.Password, config.KeyHashPassword) {
		response.Code = domain.ErrLoginHashPassword.Code
		response.Message = domain.ErrLoginHashPassword.Msg
		return http.StatusUnauthorized, response
	}

	// create token
	fmt.Println("Login")
	fmt.Println("CreateAt : ", memberBody.CreatedAt)
	fmt.Println("UserID : ", memberBody.UserID)
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}

	token, err := core_util.GenBearerToken(hashAuth, m.middlewareUtil.Encryption())
	if err != nil {
		response.Code = domain.ErrGenerateToken.Code
		response.Message = domain.ErrGenerateToken.Msg
		return http.StatusInternalServerError, response
	}

	response.Code = domain.LoginSuccess.Code
	response.Message = domain.LoginSuccess.Msg
	response.BearerToken = token
	return http.StatusOK, response
}

func (m *MemberUtil) GoogleRegister(responseGoogle any) (int, domain.ResponseRegisterMember) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.ResponseRegisterMember{}
	info, ok := responseGoogle.(domain.ResponseGoogleRegisterCallback)
	if !ok {
		response.Code = domain.ErrInvalidOAuthDecode.Code
		response.Message = domain.ErrInvalidOAuthDecode.Msg
		return http.StatusUnauthorized, response
	}

	config := util.LoadConfig()

	createAt := time.Now()
	updateAt := time.Now()
	password := util.Sha256(info.UserInfo.ID + GOOGLE)
	memberBody := domain.Member{
		UserID:       core_util.GenUserID(info.UserInfo.Email, password),
		Email:        info.UserInfo.Email,
		Hash:         core_util.GenerateHash(config.KeyHashMember),
		DisplayName:  info.UserInfo.Name,
		ProfileImage: info.UserInfo.Picture,
		Status:       PENDING,
		TypeMember:   MEMBER,
		GoogleID:     info.UserInfo.ID,
		CreatedAt:    createAt,
		UpdatedAt:    updateAt,
	}
	if err := m.memberRepo.SaveMember(ctx, memberBody); err != nil {
		response.Code = domain.ErrCreateMemberFail.Code
		response.Message = domain.ErrCreateMemberFail.Msg
		return http.StatusBadRequest, response
	}

	// create token
	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	token, err := core_util.GenBearerToken(hashAuth, m.middlewareUtil.Encryption())
	if err != nil {
		response.Code = domain.ErrGenerateToken.Code
		response.Message = domain.ErrGenerateToken.Msg
		return http.StatusInternalServerError, response
	}

	response.Code = domain.RegisterSuccess.Code
	response.Message = domain.RegisterSuccess.Msg
	response.BearerToken = token
	return http.StatusCreated, response
}

func (m *MemberUtil) GoogleLogin(responseGoogle any) (int, domain.ResponseLogin) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.ResponseLogin{}
	info, ok := responseGoogle.(domain.ResponseGoogleLoginCallback)
	if !ok {
		response.Code = domain.ErrInvalidOAuthDecode.Code
		response.Message = domain.ErrInvalidOAuthDecode.Msg
		return http.StatusBadRequest, response
	}

	loginForm := domain.LoginForm{
		Email: info.UserInfo.Email,
	}
	memberBody, err := m.memberRepo.FindEmailMember(ctx, loginForm.Email)
	if err != nil {
		response.Code = domain.ErrMemberEmailNotFound.Code
		response.Message = domain.ErrMemberEmailNotFound.Msg
		return http.StatusBadRequest, response
	}

	// create token

	hashAuth := domain.HashAuth{
		CreateAt: memberBody.CreatedAt,
		UserID:   memberBody.UserID,
	}
	token, err := core_util.GenBearerToken(hashAuth, m.middlewareUtil.Encryption())
	if err != nil {
		response.Code = domain.ErrGenerateToken.Code
		response.Message = domain.ErrGenerateToken.Msg
		return http.StatusInternalServerError, response
	}

	response.Code = domain.LoginSuccess.Code
	response.Message = domain.LoginSuccess.Msg
	response.BearerToken = token
	return http.StatusOK, response
}

func (m *MemberUtil) GetProfile(userID string) (int, domain.ResponseGetProfile) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.ResponseGetProfile{}
	member, err := m.memberRepo.GetMemberByUserID(ctx, userID)
	if err != nil {
		response.Code = domain.ErrGetMember.Code
		response.Message = domain.ErrGetMember.Msg
		return http.StatusBadRequest, response
	}

	response.Code = domain.Success.Code
	response.Message = domain.Success.Msg
	response.Member = member
	return http.StatusOK, response
}

type RedirectUtil struct {
	cache port.Cache
}

func NewRedirectUtil(cache port.Cache) *RedirectUtil {
	redirecUtil := &RedirectUtil{
		cache: cache,
	}
	return redirecUtil
}

func (m *RedirectUtil) GoogleLogin() (int, domain.ResponseRedirectGoogleLogin) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.ResponseRedirectGoogleLogin{}
	googleConfig, err := util.GoogleConfig(LOGIN)
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	if googleConfig.State, err = core_util.RandomGoogleState(); err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	ltState := time.Duration(5 * time.Minute)
	if err := m.cache.SetGoogleState(ctx, googleConfig.State, ltState); err != nil {
		response.Code = domain.ErrSetGoogleState.Code
		response.Message = domain.ErrSetGoogleState.Msg
		return http.StatusInternalServerError, response
	}

	response.Code = domain.Success.Code
	response.Message = domain.Success.Msg
	response.URL = googleConfig.Config.AuthCodeURL(
		googleConfig.State,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	return http.StatusTemporaryRedirect, response
}

func (m *RedirectUtil) GoogleRegister() (int, domain.ResponseRedirectGoogleRegister) {
	ctx, cancel := util.InitConText(2 * time.Second)
	defer cancel()

	response := domain.ResponseRedirectGoogleRegister{}
	googleConfig, err := util.GoogleConfig(REGISTER)
	if err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	if googleConfig.State, err = core_util.RandomGoogleState(); err != nil {
		response.Code = domain.ErrLoadConfig.Code
		response.Message = domain.ErrLoadConfig.Msg
		return http.StatusInternalServerError, response
	}

	ltState := time.Duration(5 * time.Minute)
	if err := m.cache.SetGoogleState(ctx, googleConfig.State, ltState); err != nil {
		response.Code = domain.ErrSetGoogleState.Code
		response.Message = domain.ErrSetGoogleState.Msg
		return http.StatusInternalServerError, response
	}

	response.Code = domain.Success.Code
	response.Message = domain.Success.Msg
	response.URL = googleConfig.Config.AuthCodeURL(
		googleConfig.State,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	return http.StatusTemporaryRedirect, response
}
