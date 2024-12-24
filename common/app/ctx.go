package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Err web服务抛出异常
func Err(errors ...interface{}) {
	errStr := ""
	for i := 0; i < len(errors); i++ {
		switch e := errors[i].(type) {
		case nil:

		case string:
			errStr += e
		case error:
			if errStr != "" {
				errStr += ":"
			}
			errStr += e.Error()
		default:
			fmt.Println(e)
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
