package routes

import (
	api "picture-oss-proxy/api/v1"
	"picture-oss-proxy/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(middleware.Cors())
	r.Use(sessions.Sessions("mysession", store))
	v1 := r.Group("api/v1")
	v1.Use(middleware.Limiter())

	{

		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})

		//获取图片
		v1.GET("picture", api.GetPicture)

		//增加jwt验证
		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{

		}

	}
	return r
}
