package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"web_template/settings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Init initializes the logger with the given configuration.
func Init(cfg *settings.LogConfig) (err error) {
	writeSyncer := getLogWriter(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxAge,
		cfg.MaxBackups,
	)
		
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)
	lg := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // replace global logger with new logger
	return nil
}

// getEncoder returns a zapcore.Encoder with the default configuration.
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()             // the default encoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         // use ISO8601 format for timestamps
	encoderConfig.TimeKey = "time"                                // use "time" as the key for timestamps
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       // use capital level names for levels
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder // use seconds for durations
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       // use short caller encoder
	return zapcore.NewJSONEncoder(encoderConfig)                  // use JSON format for logs
}
// getLogWriter returns a zapcore.WriteSyncer that writes to a lumberjack logger.
func getLogWriter(filename string, maxSize, maxAge, maxBackups int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{ // set up lumberjack logger
		Filename:   filename,   // log file path
		MaxSize:    maxSize,    // maximum size in megabytes
		MaxAge:     maxAge,     // maximum days to retain old log files
		MaxBackups: maxBackups, // maximum number of old log files to retain
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger returns a gin.HandlerFunc that logs requests using zap.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery returns a gin.HandlerFunc that recovers from any panics and logs them using zap.
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a panic, but a normal response.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error("http request",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				}else{
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}


					