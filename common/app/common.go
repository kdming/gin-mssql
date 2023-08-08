package app

import "errors"

type M map[string]interface{}
type I []interface{}

// NewError new error
func NewError(params ...interface{}) error {
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
			errStr += "发生错误!"
		}
	}
	return errors.New(errStr)
}
