package main

//go:generate swagger generate model -f ../../api/swagger.yaml -t ../../internal/

import (
	"errors"
	"fmt"
	"homework/internal/usecase"
	"log"
	"net/http"
	"os"
	"strconv"

	httpGateway "homework/internal/gateways/http"
	eventRepository "homework/internal/repository/event/inmemory"
	sensorRepository "homework/internal/repository/sensor/inmemory"
	userRepository "homework/internal/repository/user/inmemory"
)

func main() {
	er := eventRepository.NewEventRepository()
	sr := sensorRepository.NewSensorRepository()
	ur := userRepository.NewUserRepository()
	sor := userRepository.NewSensorOwnerRepository()

	useCases := httpGateway.UseCases{
		Event:  usecase.NewEvent(er, sr),
		Sensor: usecase.NewSensor(sr),
		User:   usecase.NewUser(ur, sor, sr),
	}

	httpHost := os.Getenv("HTTP_HOST")
	httpPort, err := strconv.ParseUint(os.Getenv("HTTP_PORT"), 10, 16)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "HTTP_PORT should be unsigned 16-bit integer")
		os.Exit(1)
	}

	r := httpGateway.NewServer(
		useCases,
		httpGateway.WithHost(httpHost),
		httpGateway.WithPort(uint16(httpPort)),
	)

	if err := r.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error during server shutdown: %v", err)
	}
}
