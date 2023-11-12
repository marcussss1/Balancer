//package main
//
//import (
//	"github.com/labstack/echo-contrib/prometheus"
//	"github.com/labstack/echo/v4"
//	"github.com/labstack/gommon/log"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"net/http"
//)
//
//type handler struct{}
//
//func (h handler) A(ctx echo.Context) error {
//	return ctx.NoContent(http.StatusOK)
//}
//
//func main() {
//	e := echo.New()
//	e.GET("/", handler{}.A)
//
//	p := prometheus.NewPrometheus("echo", nil)
//	eProtheus := echo.New()
//
//	e.Use(p.HandlerFunc)
//	eProtheus.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
//	go func() {
//		err := eProtheus.Start(":9000")
//		if err != nil {
//			log.Fatal(err)
//		}
//	}()
//
//	e.Logger.Fatal(e.Start(":8080"))
//}

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	echoPrometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	executionTime *prometheus.HistogramVec
	totalHits     prometheus.Counter
}

func NewPrometheusMetrics(serviceName string) *Metrics {
	m := &Metrics{
		executionTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: serviceName + "_durations",
			Help: "DurationMinutes execution of request",
		}, []string{"status", "path", "method"}),
		totalHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: serviceName + "_total_hits",
		}),
	}

	return m
}

func main() {
	m := NewPrometheusMetrics("api")

	e := echo.New()
	e.Use(echoPrometheus.MetricsMiddleware())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/", func(c echo.Context) error {
		m.totalHits.Inc()
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
