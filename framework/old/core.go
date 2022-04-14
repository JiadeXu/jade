package old

import (
	"log"
	"net/http"
	"strings"
)

// 服务框架核心结构
type Core struct {
	router      map[string]*Tree
	middlewares []ControllerHandler
}

// 初始化框架核心结构
func NewCore() *Core {
	// 初始化路由
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()

	return &Core{
		router: router,
	}
}

// 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// 对应 Method = Get
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	// 将core的middleware 和 handlers结合起来
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRoute(url, allHandlers); err != nil {
		log.Fatalln("add router error: ", err)
	}
}

// 对应 Method = POST
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	// 将core的middleware 和 handlers结合起来
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRoute(url, allHandlers); err != nil {
		log.Fatalln("add router error: ", err)
	}
}

// 对应 Method = PUT
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	// 将core的middleware 和 handlers结合起来
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRoute(url, allHandlers); err != nil {
		log.Fatalln("add router error: ", err)
	}
}

// 对应 Method = DELETE
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	// 将core的middleware 和 handlers结合起来
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRoute(url, allHandlers); err != nil {
		log.Fatalln("add router error: ", err)
	}
}

// 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteByRequest(request *http.Request) *node {
	// uri 和 method 全部转换为大写 保证大小写一致
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		// 查找第二层
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

// 框架核心结构实现Handler 接口
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// 封装对应的context
	ctx := NewContext(request, response)

	// 寻找路由
	node := c.FindRouteByRequest(request)
	if node == nil {
		ctx.Json(404).Json("not found")
		return
	}

	ctx.SetHandlers(node.handlers)

	// 设置路由参数
	params := node.parseParamsFromEndNode(request.URL.Path)
	ctx.SetParams(params)

	// 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
	if err := ctx.Next(); err != nil {
		ctx.SetStatus(500).Json("inner error")
		return
	}
}
