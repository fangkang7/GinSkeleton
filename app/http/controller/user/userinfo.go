package user

import (
	"github.com/gin-gonic/gin"
	userlogic "goskeleton/app/logic/user_logic"
	"goskeleton/app/utils/response"
)

type User struct {
}

func (u *User) UserInfo(c *gin.Context) {

	data, err := (&userlogic.User{}).GetUserInfo(c)
	if err != nil {
		//response.(c,err)
	}

	response.RequestSuccess(c, data)
}
