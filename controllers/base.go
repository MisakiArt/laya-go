package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya-go/models/dao"
	"github.com/layatips/laya/genv"
	"github.com/layatips/laya/glogs"
	"github.com/layatips/laya/gresp"
	"github.com/layatips/laya/gstore"
	"net/http"
)

// BaseCtrl the controllers with some useful and common function
var Ctrl = &BaseCtrl{}

type BaseCtrl struct {
	gresp.Resp
}

type Memories struct {
	M Condition `json:"M"`
}

type Condition struct {
	Count int                      `json:"count"`
	Item  []map[string]interface{} `json:"item"`
}

// memory status
func (ctrl *BaseCtrl) MemoryStatus(c *gin.Context) {
	var resp Memories
	resp.M.Count = dao.Mem.ItemCount()
	dmi := dao.Mem.Items()
	counter := 0
	for k, v := range dmi {
		var temp = map[string]interface{}{}
		temp[k] = v.Object
		resp.M.Item = append(resp.M.Item, temp)
		counter++
		if counter >= 10 {
			break
		}
	}
	ctrl.Suc(c, resp)
}

// version
func (ctrl *BaseCtrl) Version(c *gin.Context) {
	res := genv.AppName() + " api version: 1.0.0\r\n" + "app_url: " + genv.AppUrl()
	_, _ = c.Writer.Write([]byte(res))
	return
}

// 健康检查
func (ctrl *BaseCtrl) HealthCheck(c *gin.Context) {
	ctrl.Suc(c, "", "success")
}

// 健康检查
func (ctrl *BaseCtrl) ReadyCheck(c *gin.Context) {
	// mysql检测
	err := gstore.DBSurvive(dao.DB)
	if err != nil {
		glogs.ErrorF("探针存活检测失败,mysql凉凉,err=%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// redis检测
	err = gstore.RdbSurvive(dao.Rdb)
	if err != nil {
		glogs.ErrorF("探针存活检测失败,redis凉凉,err=%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// mongo检测
	err = gstore.MdbSurvive(dao.Mdb)
	if err != nil {
		glogs.ErrorF("探针存活检测失败,mongodb凉凉,err=%s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctrl.Suc(c, "pong", "success")
}

// 重新载入配置
func (ctrl *BaseCtrl) Reload(c *gin.Context) {
	ctrl.Suc(c, "pong", "success")
}
