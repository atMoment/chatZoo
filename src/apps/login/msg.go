package main

type LoginReq struct {
	ID        string `json:"id"`         // 也可以是name?
	IsVisitor bool   `json:"is_visitor"` // 是否游客
}

type RegisterReq struct {
	Name string `json:"name"`
}

type LoginResp struct {
	Code int32  `json:"code"`
	Err  string `json:"err"` // 错误信息
}
