package resp

import "fmt"

// 基础resp
type Response struct {
	ErrNo  int         `json:"err_no" yaml:"err_no"`
	ErrMsg string      `json:"err_msg" yaml:"err_msg"`
	Data   interface{} `json:"data" yaml:"data"`
}

func EmptyResp() *Response {
	return new(Response)
}

func SuccessResp() *Response {
	return &Response{ErrNo: 0, ErrMsg: "success"}
}

// 覆盖原来的Data
func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

// 返回带着 err的信息
func (r *Response) WithErr(err error) *Response {
	r.ErrNo = 300
	r.ErrMsg = fmt.Sprintf("failed：%s", err)
	return r
}
