package init

import (
	"context"
	"github.com/sheginabo/go-quick-gin/init/config"
	"github.com/sheginabo/go-quick-gin/init/gin"
	"github.com/sheginabo/go-quick-gin/init/logger"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type MainInitProcess struct {
	Log       *logger.Module
	Gin       *gin.Module
	OsChannel chan os.Signal
}

var interruptSignals = []os.Signal{
	syscall.SIGTERM,
	syscall.SIGINT,
}

func NewMainInitProcess(configPath string) *MainInitProcess {
	config.NewModule(configPath)
	logModule := logger.NewModule()
	ginModule := gin.NewModule()

	channel := make(chan os.Signal, 1)
	return &MainInitProcess{
		Log:       logModule,
		Gin:       ginModule,
		OsChannel: channel,
	}
}

// Run run gin module
func (m *MainInitProcess) Run() {
	// 使用一種 context 來管理多個 goroutine 的生命週期, 註冊三個取消訊號 (SIGINT, SIGTERM, SIGQUIT)
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	// 當函數返回時，取消訊號會被發送到 ctx，然後取消 ctx，這樣所有使用 ctx 的 goroutine 都會被取消
	defer stop()
	// 使用 errgroup 來管理多個 goroutine 的生命週期
	waitGroup, ctx := errgroup.WithContext(ctx)

	m.Gin.Run(ctx, waitGroup)

	// 等待所有 goroutine 完成
	err := waitGroup.Wait()
	if err != nil {
		m.Log.Logger.Fatal().Err(err).Msg("error from wait group")
	}
}
