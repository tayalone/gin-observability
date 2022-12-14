package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tayalone/gin-obserbability/user/nested"
	ess "github.com/tayalone/go-ess-package/otel"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	service     = "user"
	environment = "dev"
)

func main() {
	// // ------ call tp -------------------------------------------

	tp, err := ess.JaegertracerProvider(os.Getenv("JEAGER_ENDPOINT"), service, environment)
	if err != nil {
		log.Fatal(err)
	}

	otelCtx := context.Background()
	defer func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(otelCtx)

	// // ---------------------------------------------------------
	logger, _ := zap.NewProduction(zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	// // ---------------------------------------------------------
	// // ---------------------------------------------------------
	otelLogger := otelzap.New(logger, otelzap.WithMinLevel(zap.DebugLevel), otelzap.WithTraceIDField(true))
	undo := otelzap.ReplaceGlobals(otelLogger)
	defer undo()

	otelSugar := otelLogger.Sugar()

	// // ----------------------------------------------------------
	r := gin.Default()

	r.Use(otelgin.Middleware(service))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			otelSugar.Ctx(c.Request.Context()).Errorw(err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// // ---------- router ------------------------
	r.GET("/ping", func(c *gin.Context) {
		type Ping struct {
			Name    string `json:"name"`
			Duraion int    `json:"duraion"`
		}

		sugar.Infow("return pong",
			"input", 1,
			"ping", Ping{Name: "John", Duraion: 100},
		)

		otelSugar.Ctx(c.Request.Context()).Infow("return pong otel",
			"input", 1,
			"ping", Ping{Name: "John", Duraion: 100},
		)

		sugar.Warnw("try warn")
		otelSugar.Ctx(c.Request.Context()).Warnw("return pong otel",
			"input", 1,
			"ping", Ping{Name: "John", Duraion: 100},
		)

		sugar.Errorw("try error")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("Get Panic")
	})

	r.GET("/internal-tracing", func(c *gin.Context) {
		time.Sleep(150 * time.Millisecond)

		nested.Parent(c.Request.Context())

		time.Sleep(100 * time.Millisecond)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	r.GET("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "get User data",
		})
	})

	r.PATCH("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "edit user data",
		})
	})

	r.GET("/:id/todo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "get all user's todo",
		})
	})

	r.POST("/:id/todo", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "create new  todo user data",
		})
	})

	r.PATCH("/:id/todo/:todoId", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "edit  todo user data",
		})
	})

	r.DELETE("/:id/todo/:todoId", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "delete  todo user data",
		})
	})
	// // -------------------------------------------------
	// r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	// // run srv in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// // create channel of os signal for waiting signal
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// // close `quit` channel when app're closed
	// defer func() {
	// 	close(quit)
	// }()
	s := <-quit
	log.Println("signal is: ", s)
	log.Println("Shutting down app...")

	// // The context is used to inform the server it has 5 seconds to finish
	// // the request it is currently handling
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// // create context -> waiting every process in server is done
	srvCtx := context.Background()

	if err := srv.Shutdown(srvCtx); err != nil {
		log.Fatal("App forced to shutdown:", err)
	}

	log.Println("App exiting")
}

func bar(ctx context.Context) {
	// Use the global TrazcerProvider.
	tr := otel.Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	defer span.End()
	// Do bar...
}

func delay(ctx context.Context) {
	// Use the global TrazcerProvider.
	tr := otel.Tracer("component-bar")
	_, span1 := tr.Start(ctx, "delay-1-sec")
	span1.SetAttributes(attribute.Key("desc").String("I delay 1sec"))
	time.Sleep(time.Second * 1)
	span1.End()

	_, span2 := tr.Start(ctx, "delay-1.5-sec")
	span2.SetAttributes(attribute.Key("desc").String("I delay 1.5sec again"))
	time.Sleep(time.Second * 1)
	span2.End()
}
