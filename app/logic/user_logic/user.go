package userlogic

import (
	"github.com/gin-gonic/gin"
	perror2 "goskeleton/app/utils/perror"
	"goskeleton/app/utils/response"
)

type User struct {
}

func (u User) GetUserInfo(c *gin.Context) (*response.SuccessData, perror2.Error) {
	//myMap := map[string]interface{}{
	//	"apple":  5,
	//	"banana": 10,
	//	"cherry": 15,
	//}

	mySlice := []int{1, 2, 3, 4}
	// 这里随便模拟一条数据返回
	return response.ReturnSuccessData(mySlice)
}
