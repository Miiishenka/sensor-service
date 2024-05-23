package main

//go:generate swagger generate model -f ../../api/swagger.yaml -t ../../internal/

import (
	"context"
	"errors"
	"homework/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	httpGateway "homework/internal/gateways/http"
	eventRepository "homework/internal/repository/event/inmemory"
	sensorRepository "homework/internal/repository/sensor/inmemory"
	userRepository "homework/internal/repository/user/inmemory"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("can't parse pgxpool config")
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("can't create new pool")
	}
	defer pool.Close()

	er := eventRepository.NewEventRepository()
	sr := sensorRepository.NewSensorRepository()
	ur := userRepository.NewUserRepository()
	sor := userRepository.NewSensorOwnerRepository()

	useCases := httpGateway.UseCases{
		Event:  usecase.NewEvent(er, sr),
		Sensor: usecase.NewSensor(sr),
		User:   usecase.NewUser(ur, sor, sr),
	}

	var functions []func(*httpGateway.Server)

	if httpHost, ok := os.LookupEnv("HTTP_HOST"); ok {
		functions = append(functions, httpGateway.WithHost(httpHost))
	}

	if httpPortEnv, ok := os.LookupEnv("HTTP_PORT"); ok {
		httpPort, err := strconv.ParseUint(httpPortEnv, 10, 16)
		if err == nil {
			functions = append(functions, httpGateway.WithPort(uint16(httpPort)))
		}
	}

	r := httpGateway.NewServer(
		useCases,
		functions...,
	)

	if err := r.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}
