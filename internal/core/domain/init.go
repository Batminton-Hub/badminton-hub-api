package domain

const ( // Platform
	NORMAL = "NORMAL"
	GOOGLE = "GOOGLE"
)

const ( // Action
	REGISTER = "REGISTER"
	LOGIN    = "LOGIN"
)

const ( // Status
	PENDING = "PENDING"
	ACTIVE  = "ACTIVE"
	BANNED  = "BANNED"
	DELETED = "DELETED"
	SUCCESS = "SUCCESS"
	ERROR   = "ERROR"
)

const ( // Type Member
	MEMBER = "MEMBER"
)

const ( // Context
	Authorization = "Authorization"
	Platform      = "platform"
	State         = "state"
	Code          = "code"
	API           = "api"
	PlatformData  = "platform_data"
	UserID        = "user_id"
)

const ( // Google Config
	GoogleUserInfoURL     = "https://www.googleapis.com/oauth2/v2/userinfo"
	GoogleUserInfoEmail   = "https://www.googleapis.com/auth/userinfo.email"
	GoogleUserInfoProfile = "https://www.googleapis.com/auth/userinfo.profile"
)
