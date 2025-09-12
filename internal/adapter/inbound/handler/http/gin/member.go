package gin

import (
	"Badminton-Hub/internal/core/domain"
	"Badminton-Hub/internal/core/port"

	"github.com/gin-gonic/gin"
)

type MemberController interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

type Member struct {
	member port.MemberService
}

func (m *Member) GetProfile(c *gin.Context) {
	userID := getUserID(c)
	userInfo := domain.ReqGetProfile{
		UserID: userID,
	}

	httpStatus, response := m.member.GetProfile(userInfo)
	if response.Resp.Status == domain.ERROR {
		Resp(c, httpStatus, response.Resp.Code, response.Resp.Msg, nil)
		return
	}
	Resp(c, httpStatus, response.Resp.Code, response.Resp.Msg, response.Member)
}

func (m *Member) UpdateProfile(c *gin.Context) {
	userInfo := domain.ReqGetProfile{
		UserID: getUserID(c),
	}
	request := domain.ReqUpdateProfile{}
	if err := c.ShouldBindJSON(&request); err != nil {
		Resp(c, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Code, domain.ErrInvalidInput.Msg, nil)
		return
	}

	httpStatus, response := m.member.UpdateProfile(userInfo, request)
	if response.Resp.Status == domain.ERROR {
		Resp(c, httpStatus, response.Resp.Code, response.Resp.Msg, nil)
		return
	}
	Resp(c, httpStatus, response.Resp.Code, response.Resp.Msg, nil)
}
