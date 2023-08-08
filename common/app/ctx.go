package app

import (
	"github.com/gin-gonic/gin"
)

// Err web服务抛出异常
func Err(params ...interface{}) {
	errStr := ""
	for i := 0; i < len(params); i++ {
		p := params[i]
		switch p.(type) {
		case string:
			errStr += p.(string)
		case error:
			if errStr != "" {
				errStr += ":"
			}
			errStr += p.(error).Error()
		default:
			errStr += "发生错误，无法判断错误类型"
		}
	}
	panic(errStr)
}

// Ok Success正常返回数据
func Ok(c *gin.Context, msg string, data interface{}) {
	res := M{}
	res["code"] = 0
	res["msg"] = msg
	res["data"] = data
	c.JSON(200, res)
}
