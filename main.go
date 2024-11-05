package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_template/routes"
	"web_template/settings"
	"web_template/utils/logger"
	"web_template/utils/mysql"
	"web_template/utils/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	//load settings
	if err := settings.Init(); err != nil {
		fmt.Println("initialization of settings failed", err)
		return
	}
	//initialize the logger
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Println("initialization of logger failed", err)
		return
	}
	defer zap.L().Sync()  // flushes buffer, if any

	//initialize the MYSQL connection
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("initialization of MYSQL connection failed", err)
		return
	}
	defer mysql.Close()  //close the connection when the program exits

	//initialize the Redis connection
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("initialization of Redis connection failed", err)
		return
	}
	defer redis.Close()  //close the connection when the program exits

	//rigister the routers
	r := routes.SetUp()

	//Start the server(shut down gracefully)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	//Start the server in a goroutine so that it doesn't block the main thread
	go func() {
		//create a goroutine to start the server
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	//Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal,1)  // Create channel to receive signal notifications 
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)  // Register signal handlers
	<-quit  // Block until a signal is received.
	zap.L().Info("Shutting down server...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Create a context with timeout=5s
	defer cancel()  // The cancel should be deferred in case of a timeout or other error.

	//Shutdown the server, waiting for it to complete or timeout, and log any errors, if any.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exited ")
}
