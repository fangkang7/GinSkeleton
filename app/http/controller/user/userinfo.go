package user

import (
	"github.com/gin-gonic/gin"
	userlogic "goskeleton/app/logic/user_logic"
)

type User struct {
}

func (u *User) UserInfo(c *gin.Context) {

	(&userlogic.User{}).GetUserInfo(c)
}
