package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	userlogic "goskeleton/app/logic/user_logic"
	"testing"
)

func TestUserInfo(t *testing.T) {
	c := &gin.Context{}
	userInfo, err := (&userlogic.User{}).GetUserInfo(c)
	if err != nil {
		fmt.Println(err)
	}

	data, _ := json.Marshal(userInfo)
	t.Log(string(data))
}
