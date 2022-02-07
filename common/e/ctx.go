package e

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type M map[string]interface{}
type I []interface{}

func Err(msg string, err error) {
	if err != nil {
		msg += ":" + err.Error()
	}
	panic(msg)
}

func Ok(c *gin.Context, msg string, data interface{}) {
	res := M{}
	res["code"] = 0
	res["msg"] = msg
	res["data"] = data
	c.JSON(200, res)
}

func NewError(msg string, err error) error {
	if err == nil {
		return errors.New(msg)
	}
	return errors.New(msg + ":" + err.Error())
}
