package login

type LoginReq struct {
	ID        string `json:"id"` // 也可以是name?
	Pwd       string `json:"pwd"`
	IsVisitor bool   `json:"is_visitor"` // 是否游客
	PublicKey string `json:"public_key"`
}

type RegisterReq struct {
	Name      string `json:"name"`
	Pwd       string `json:"pwd"`
	PublicKey string `json:"public_key"`
}

type LoginResp struct {
	Code      int32  `json:"code"`
	Err       string `json:"err"` // 错误信息
	PublicKey string `json:"public_key"`
	GateAddr  string `json:"gate_addr"` // 网关地址
}
