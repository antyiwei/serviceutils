package gindefault

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(addr string, initial func(engine *gin.Engine)) {

	time.LoadLocation("Asia/Shanghai")

	engine := gin.New()

	engine.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "X-Auth-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	initial(engine)

	log.Println("Starting http server at address:", addr)

	engine.Run(addr)
}

func Run2(addr string, initial func(engine *gin.Engine), doSomethingBeforeExit func() error) {

	time.LoadLocation("Asia/Shanghai")

	engine := gin.New()

	engine.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "X-Auth-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	initial(engine)

	log.Println("Starting http server at address:", addr)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// 启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待一个INT或TERM信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// 创建超时上下文，Shutdown可以让未处理的连接在这个时间内关闭
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	// 停止HTTP服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// 执行自定义的退出函数
	if doSomethingBeforeExit != nil {
		if err := doSomethingBeforeExit(); err != nil {
			log.Println("beforeExit error:", err)
		}
	}

	log.Println("Server exiting")
}
