package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/DwGoing/MarketBrain/internal/service"
	"github.com/DwGoing/MarketBrain/pkg/enum"
	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mkideal/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// 应用名
	AppName = "Market Brain"
	// 版本号
	Version = "0.0.1"
	// 配置路径
	ConfigPath = "./config"
)

// 定义参数
type rootArgv struct {
	Help bool `cli:"!h,help" usage:"Help information"`
}

var rootCommand = &cli.Command{
	Desc: "Market Brain",
	Argv: func() interface{} { return new(rootArgv) },
	Fn: func(ctx *cli.Context) error {
		if len(ctx.NativeArgs()) < 1 {
			ctx.WriteUsage()
			return nil
		}
		return nil
	},
}

var versionCommand = &cli.Command{
	Name: "version",
	Desc: "version infomation",
	Fn: func(ctx *cli.Context) error {
		fmt.Printf("Version: %s", Version)
		return nil
	},
}

type serviceArgv struct {
	Help bool   `cli:"!h,help" usage:"Help information"`
	Type string `cli:"t,type" usage:"Service type. Options: FUNDS|DATA"`
}

var serviceCommand = &cli.Command{
	Name: "service",
	Desc: "Market brain service",
	Argv: func() interface{} { return new(serviceArgv) },
	Fn: func(ctx *cli.Context) error {
		if len(ctx.NativeArgs()) < 1 {
			ctx.WriteUsage()
			return nil
		}
		argv := ctx.Argv().(*serviceArgv)
		serviceType, err := new(enum.ServiceType).Parse(argv.Type)
		if err != nil {
			return err
		}
		switch serviceType {
		case enum.ServiceType_FUNDS:
			err = startFundsService()
			if err != nil {
				return err
			}
		default:
			return errors.New("unsupported service")
		}
		zap.S().Infoln("process running")
		//监听指定信号 ctrl+c kill
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		zap.S().Infof("process finished")
		return nil
	},
}

func main() {
	// 日志配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 获取执行路径
	executablePath, err := os.Executable()
	if err != nil {
		zap.S().Errorf("%s", err)
		os.Exit(1)
	}
	executablePath = filepath.Dir(executablePath)
	hook, err := rotatelogs.New(
		path.Join(executablePath, "./log", "%Y%m%d.log"),
		rotatelogs.WithRotationCount(5),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		zap.S().Errorf("%s", err)
		os.Exit(1)
	}
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)
	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(hook),
		zap.InfoLevel,
	)
	logger := zap.New(
		zapcore.NewTee(consoleCore, fileCore),
		zap.AddCaller(),
		zap.Development(),
	)
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	if err := cli.Root(
		rootCommand,
		cli.Tree(versionCommand),
		cli.Tree(serviceCommand),
	).Run(os.Args[1:]); err != nil {
		zap.S().Errorf("%s", err)
		os.Exit(1)
	}
}

func startFundsService() error {
	// 加载配置
	err := ioc.Load(
		config.WithSearchPath(ConfigPath),
		config.WithProfilesActive("funds"),
	)
	if err != nil {
		return err
	}
	// 实例化服务
	_, err = service.GetFundsSingleton()
	if err != nil {
		return err
	}
	return nil
}
