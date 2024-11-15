package main

/*
http 服务器, 职责如下
接收客户端login请求, 交换公钥, 存储数据库, 发送 gate地址
问题：login时存数据库会不会被爆冲？即一堆http请求来注册 (正常游戏有sdk验证,但是我没有)
// todo 是否需要有流量监控？流量太大情况下登录不上  排队服？
*/

func main() {
	srv := NewSrv()
	srv.run()
}
