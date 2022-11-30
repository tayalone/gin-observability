package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	service     = "user"
	environment = "dev"
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func tracerProvider(url, serviceName, envName string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			/* can add other attribure
			attribute.Int64("ID", id),
			*/
		)),
	)
	return tp, nil
}

func main() {
	// // ------ call tp -------------------------------------------
	tp, err := tracerProvider(os.Getenv("JEAGER_ENDPOINT"), service, environment)
	if err != nil {
		log.Fatal(err)
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	otelCtx := context.Background()
	defer func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(otelCtx)

	// // ---------------------------------------------------------
	r := gin.Default()

	r.Use(otelgin.Middleware(service))

	// // ---------- router ------------------------
	r.GET("/ping", func(c *gin.Context) {
		bar(c.Request.Context())
		delay(c.Request.Context())
		bar(c.Request.Context())
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
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
