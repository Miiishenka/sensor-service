package http

import (
	"encoding/json"
	"errors"
	"homework/internal/domain"
	"homework/internal/models"
	"homework/internal/usecase"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/go-openapi/strfmt"

	"github.com/gin-gonic/gin"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
	httpRequestsErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"path", "method"},
	)
)

func setupRouter(r *gin.Engine, uc UseCases, h *WebSocketHandler) {
	r.HandleMethodNotAllowed = true

	r.Use(func(c *gin.Context) {
		path := c.FullPath()
		method := c.Request.Method

		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(path, method))
		defer timer.ObserveDuration()

		c.Next()

		httpRequestsTotal.WithLabelValues(path, method).Inc()
		if c.Writer.Status() >= 400 {
			httpRequestsErrorsTotal.WithLabelValues(path, method).Inc()
		}
	})

	r.POST("/users", func(c *gin.Context) {
		if c.ContentType() != "application/json" {
			c.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}

		userDto := &models.UserToCreate{}
		if err := c.BindJSON(userDto); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := userDto.Validate(nil); err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		user := &domain.User{Name: *userDto.Name}
		user, err := uc.User.RegisterUser(c, user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.OPTIONS("/users", func(c *gin.Context) {
		methods := []string{http.MethodOptions, http.MethodPost}
		c.Header("Allow", strings.Join(methods, ","))
		c.Status(http.StatusNoContent)
	})

	r.GET("/sensors", func(c *gin.Context) {
		if c.GetHeader("Accept") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		sensors, err := uc.Sensor.GetSensors(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, sensors)
	})

	r.HEAD("/sensors", func(c *gin.Context) {
		if c.GetHeader("Accept") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		sensors, err := uc.Sensor.GetSensors(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
		}
		jsonSensor, _ := json.Marshal(sensors)
		c.Header("Content-Length", strconv.Itoa(len(jsonSensor)))
		c.Status(http.StatusOK)
	})

	r.POST("/sensors", func(c *gin.Context) {
		if c.ContentType() != "application/json" {
			c.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}

		sensorDto := &models.SensorToCreate{}
		if err := c.BindJSON(sensorDto); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := sensorDto.Validate(nil); err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		sensor := &domain.Sensor{
			Description:  *sensorDto.Description,
			IsActive:     *sensorDto.IsActive,
			SerialNumber: *sensorDto.SerialNumber,
			Type:         domain.SensorType(*sensorDto.Type),
		}
		sensor, err := uc.Sensor.RegisterSensor(c, sensor)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, sensor)
	})

	r.OPTIONS("/sensors", func(c *gin.Context) {
		methods := []string{
			http.MethodOptions,
			http.MethodPost,
			http.MethodGet,
			http.MethodHead,
		}

		c.Header("Allow", strings.Join(methods, ","))
		c.Status(http.StatusNoContent)
	})

	r.GET("/sensors/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			return
		}

		sensor, err := uc.Sensor.GetSensorByID(c, id)
		if errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		if c.GetHeader("Accept") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		c.JSON(http.StatusOK, sensor)
	})

	r.HEAD("/sensors/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			return
		}

		sensor, err := uc.Sensor.GetSensorByID(c, id)
		if errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		if c.GetHeader("Accept") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		jsonSensor, _ := json.Marshal(sensor)
		c.Header("Content-Length", strconv.Itoa(len(jsonSensor)))
		c.Status(http.StatusOK)
	})

	r.OPTIONS("/sensors/:id", func(c *gin.Context) {
		methods := []string{
			http.MethodOptions,
			http.MethodGet,
			http.MethodHead,
		}

		c.Header("Allow", strings.Join(methods, ","))
		c.Status(http.StatusNoContent)
	})

	r.GET("/users/:id/sensors", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			return
		}

		sensors, err := uc.User.GetUserSensors(c, id)
		if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}

		if c.GetHeader("Accept") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		c.JSON(http.StatusOK, sensors)
	})

	r.HEAD("/users/:id/sensors", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
			return
		}

		sensors, err := uc.User.GetUserSensors(c, id)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if c.GetHeader("Accept") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		jsonSensors, _ := json.Marshal(sensors)
		c.Header("Content-Length", strconv.Itoa(len(jsonSensors)))
		c.Status(http.StatusOK)
	})

	r.POST("/users/:id/sensors", func(c *gin.Context) {
		userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if c.ContentType() != "application/json" {
			c.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}

		sensorDto := &models.SensorToUserBinding{}
		if err := c.BindJSON(sensorDto); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := sensorDto.Validate(nil); err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		err = uc.User.AttachSensorToUser(c, userId, *sensorDto.SensorID)
		if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusCreated)
	})

	r.OPTIONS("/users/:id/sensors", func(c *gin.Context) {
		methods := []string{
			http.MethodOptions,
			http.MethodPost,
			http.MethodHead,
			http.MethodGet,
		}

		c.Header("Allow", strings.Join(methods, ","))
		c.Status(http.StatusNoContent)
	})

	r.POST("/events", func(c *gin.Context) {
		if c.ContentType() != "application/json" {
			c.AbortWithStatus(http.StatusUnsupportedMediaType)
			return
		}

		eventDto := &models.SensorEvent{}
		if err := c.BindJSON(eventDto); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := eventDto.Validate(nil); err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		event := &domain.Event{
			Payload:            *eventDto.Payload,
			SensorSerialNumber: *eventDto.SensorSerialNumber,
			Timestamp:          time.Now(),
		}

		err := uc.Event.ReceiveEvent(c, event)
		if errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusCreated)
	})

	r.OPTIONS("/events", func(c *gin.Context) {
		methods := []string{http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPost}
		c.Header("Allow", strings.Join(methods, ","))
		c.Status(http.StatusNoContent)
	})

	r.GET("/sensors/:id/events", func(c *gin.Context) {
		sensorId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if _, err = uc.Sensor.GetSensorByID(c, sensorId); errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err = h.Handle(c, sensorId); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
	})

	r.OPTIONS("/sensors/:id/history", func(c *gin.Context) {
		methods := []string{http.MethodOptions, http.MethodGet, http.MethodHead}
		c.Header("Allow", strings.Join(methods, ","))
		c.Status(http.StatusNoContent)
	})

	r.GET("/sensors/:id/history", func(c *gin.Context) {
		sensorId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if _, err := uc.Sensor.GetSensorByID(c, sensorId); errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		startDate, err := time.Parse("2006-01-02T15:04:05.999Z", c.Query("start_date"))
		if err != nil {
			startDate = time.Time{}
		}

		endDate, err := time.Parse("2006-01-02T15:04:05.999Z", c.Query("end_date"))
		if err != nil {
			endDate = time.Now()
		}

		events, err := uc.Event.GetSensorHistory(c, sensorId, startDate, endDate)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var sensorStates []models.SensorState
		for _, event := range events {
			timestamp := strfmt.DateTime(event.Timestamp)
			payload := event.Payload
			sensorStates = append(sensorStates, models.SensorState{
				Payload:   &payload,
				Timestamp: &timestamp,
			})
		}

		if c.GetHeader("Accept") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		c.JSON(http.StatusOK, sensorStates)
	})

	r.HEAD("/sensors/:id/history", func(c *gin.Context) {
		if c.GetHeader("Accept") != "application/json" {
			c.AbortWithStatus(http.StatusNotAcceptable)
			return
		}

		sensorId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if _, err := uc.Sensor.GetSensorByID(c, sensorId); errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		startDate, err := time.Parse("2006-01-02T15:04:05.999Z", c.Query("start_date"))
		if err != nil {
			startDate = time.Time{}
		}

		endDate, err := time.Parse("2006-01-02T15:04:05.999Z", c.Query("end_date"))
		if err != nil {
			endDate = time.Now()
		}

		events, err := uc.Event.GetSensorHistory(c, sensorId, startDate, endDate)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var sensorStates []models.SensorState
		for _, event := range events {
			payload := event.Payload
			timestamp := strfmt.DateTime(event.Timestamp)
			sensorStates = append(sensorStates, models.SensorState{
				Payload:   &payload,
				Timestamp: &timestamp,
			})
		}

		jsonSensor, _ := json.Marshal(sensorStates)
		c.Header("Content-Length", strconv.Itoa(len(jsonSensor)))
		c.Status(http.StatusOK)
	})
}
