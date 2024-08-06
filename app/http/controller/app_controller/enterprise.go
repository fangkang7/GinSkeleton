package app_controller

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/logic/app_logic"
	"goskeleton/app/utils/response"
)

type Enterprise struct {

}

func (e *Enterprise) EnterpriseInfo(c *gin.Context)  {
	enterpriseInfo, err := (&app_logic.Enterprise{}).EnterpriseInfo(c)

	if err != nil {
		response.RequestFail(c,err)
	}

	response.RequestSuccess(c,enterpriseInfo)
}