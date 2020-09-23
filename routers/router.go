package routers

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"gonelist/api"
	"gonelist/conf"
	"gonelist/middleware"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Options{
		AllowedHeaders: []string{"pass"}, // 允许 header,
		AllowedOrigins: []string{"*"},
	}))

	r.Use(static.Serve("/", static.LocalFile(conf.GetDistPATH(), false)))

	// 测试接口
	r.GET("/testapi", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	r.GET("/login", api.Login)
	r.GET("/loginmg", api.LoginMG)
	r.GET("/auth", api.GetCode)
	//r.GET("/cancelLogin", api.CancelLogin)

	// 直接下载接口
	root := r.Group("/")
	root.Use(middleware.CheckLogin())
	{
		r.GET("/d/*path", api.Download)
		r.GET("/README", middleware.CheckFolderPass(), api.GetREADME)
		r.GET("/search", api.Search)
	}

	onedrive := r.Group("/onedrive")
	onedrive.Use(middleware.CheckLogin())
	{
		// 主动获取所有文件，返回整个树的目录
		onedrive.GET("/getallfiles", api.MGGetFileTree)
		// 根据路径获取对应数据
		onedrive.GET("/getpath", middleware.CheckFolderPass(), api.CacheGetPath)
	}

	admin := r.Group("/admin")
	admin.Use(middleware.CheckLogin())
	{
		admin.GET("/refresh", middleware.AdminManualRefresh())
	}

	info := r.Group("/info")
	info.Use(middleware.CheckLogin())
	{
		info.GET("/site", middleware.GetSiteInfo())
	}

	return r
}
