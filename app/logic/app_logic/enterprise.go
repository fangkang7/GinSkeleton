package app_logic

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/perror"
	"goskeleton/app/utils/response"
)

type Enterprise struct {

}

func (e *Enterprise) EnterpriseInfo(c *gin.Context) (*response.SuccessData,perror.Error)  {

	returnData := make(map[string]interface{})
	returnData["all_num"] = '1';

	return response.ReturnSuccessData(returnData)
}