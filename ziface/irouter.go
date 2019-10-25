package ziface


// IRouter 服务端应用可以给Zinx框架配置当前链接的处理业务方法
type IRouter interface {
	// 在处理业务链接之前的handle 方法
	PreHandle(request IRequest)
	// 处理业务链接的方法
	Handle(request IRequest)
	// 在处理业务链接后的handle方法
	PostHandle(request IRequest)
}
