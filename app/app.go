package app

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"funds-system/controller"
	"funds-system/docs"
	"funds-system/service/chain_service"
	"funds-system/service/config_service"
	"funds-system/service/funds_service"

	"github.com/alibaba/ioc-golang/extension/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewApp
type App struct {
	GinPort          *config.ConfigInt64           `config:",app.gin.port"`
	GrpcPort         *config.ConfigInt64           `config:",app.grpc.port"`
	ConfigService    *config_service.ConfigService `singleton:""`
	ChainService     *chain_service.ChainService   `singleton:""`
	FundsService     *funds_service.FundsService   `singleton:""`
	ConfigController *controller.ConfigController  `singleton:""`
	FundsController  *controller.FundsController   `singleton:""`

	logger *log.Logger
}

// @title	Funds System
// @version	1.0
// @query.collection.format	multi
/*
@title	构造函数
@param 	app *App 	App实例
@return _ 	*App 	App实例
@return _ 	error 	异常信息
*/
func NewApp(app *App) (*App, error) {
	app.logger = log.New(os.Stderr, "[App]", log.LstdFlags)
	return app, nil
}

/*
@title	初始化
@param 	Self	*App 	App实例
@return _ 		*App 	App实例
@return _ 		error 	异常信息
*/
func (Self *App) Initialize() {
	// gin
	go func() {
		docs.SwaggerInfo.BasePath = "/"
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Self.GinPort.Value()))
		if err != nil {
			Self.logger.Fatalf("Gin初始化失败:%s", err)
		}
		engine := gin.Default()
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		engine.POST("/callback", func(ctx *gin.Context) {
			body, _ := io.ReadAll(ctx.Request.Body)
			log.Printf("%s", body)
			defer ctx.Request.Body.Close()
		})
		v1Router := engine.Group("/v1")
		{
			configRouter := v1Router.Group("config")
			{
				configRouter.GET("/load", Self.ConfigController.Load)
				configRouter.POST("/set", Self.ConfigController.Set)
			}
			fundsRouter := v1Router.Group("/funds")
			{
				fundsRouter.POST("/getRechargeWallet", Self.FundsController.GetRechargeWallet)
				fundsRouter.GET("/getRechargeRecords", Self.FundsController.GetRechargeRecords)
			}
		}
		Self.logger.Printf("Gin正在监听: %s", listener.Addr())
		if err := engine.RunListener(listener); err != nil {
			Self.logger.Fatalf("Gin启动失败: %s", err)
		}
	}()

	// gRPC
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Self.GrpcPort.Value()))
		if err != nil {
			Self.logger.Fatalf("gRPC初始化失败:%s", err)
		}
		server := grpc.NewServer()
		config_service.RegisterConfigServiceServer(server, Self.ConfigService)
		chain_service.RegisterChainServiceServer(server, Self.ChainService)
		funds_service.RegisterFundsServiceServer(server, Self.FundsService)
		Self.logger.Printf("gRPC正在监听: %s", listener.Addr())
		if err = server.Serve(listener); err != nil {
			Self.logger.Fatalf("gRPC启动失败: %s", err)
		}
	}()
}
