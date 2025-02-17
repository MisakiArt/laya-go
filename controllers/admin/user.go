package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya-go/models/page/admin"
	"github.com/layatips/laya/glogs"
)

type getUserListParam = admin.UserParam

func (ctrl *controller) GetUserList(c *gin.Context) {
	var param getUserListParam
	err := c.ShouldBind(&param)
	if err != nil {
		ctrl.Fail(c, err)
		return
	}

	resp, err := admin.GetUserList(c, &param)
	if err != nil {
		glogs.ErrorFR(c, "title=获取用户列表,err=%s", err.Error())
		ctrl.Fail(c, err)
		return
	}
	ctrl.Suc(c, resp)
}
