package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type (
	shutdownFunc func() // 退出处理函数
)

// @Description 系统信号捕获

func SignalHandler(appletName string, shutdown shutdownFunc) {
	var (
		ch         = make(chan os.Signal, 10)
		shutdownCh = make(chan struct{})
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

WAIT_SIGNAL:
	for {
		sig := <-ch
		fmt.Printf("Get a signal: %s, will stop the program: %s\n", sig.String(), appletName)
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			// 安全关闭服务
			go func() {
				shutdown()
				close(shutdownCh)
			}()

			break WAIT_SIGNAL

		default:
			return
		}
	}

	fmt.Println("shutting down server ...")

	<-shutdownCh
	fmt.Println("Gracefully shutting down server ...")
}
