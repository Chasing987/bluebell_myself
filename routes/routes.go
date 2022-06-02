package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "bluebell/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	//注册
	v1.POST("/signup", controller.SignUpHandler)
	//登录
	v1.POST("/login", controller.LoginHandler)

	// 根据时间或者分数获取帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)
	v1.GET("/posts", controller.GetPostListHandler)
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/post/:id", controller.GetPostDetailHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT认证中间件

	{
		v1.POST("/post", controller.CreatePostHandler)
		//投票
		v1.POST("/vote", controller.PostVoteController)
	}

	pprof.Register(r) // 注册pprof 相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r

	////注册业务路由
	////注册
	//r.POST("/signup", controller.SignUpHandler)
	////登录
	//r.POST("/login", controller.LoginHandler)
	//
	//r.GET("/version", func(c *gin.Context) {
	//	c.String(http.StatusOK, settings.Conf.Version)
	//})
	//r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//	// 如果是登录的用户，判断请求头中是否有 有效的JWT？
	//	c.String(http.StatusOK, "pong")
	//})
	//return r
}
