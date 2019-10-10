package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"github.com/nicolasdeveloper/easymile-tracker/common/broaker"
	"github.com/nicolasdeveloper/easymile-tracker/common/models"
)

func handleUDPConnection(conn *net.UDPConn, massage *broaker.AmqpClient) {

	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)

	strBuffer := string(buffer[:n])

	fmt.Println("UDP client : ", addr)
	fmt.Println("Received from UDP client :  ", strBuffer)

	if err != nil {
		log.Println(err)
	}

	trip, err := extractDataFromStringBuffer(strBuffer)

	if err != nil {
		log.Println(err)
	}

	if trip != nil {
		message, err := json.Marshal(trip)

		go func() {
			massage.Publish(message, trip.VINNumber, trip.VINNumber, "tracker", "direct")
		}()

		if err != nil {
			log.Println(err)
		}
	}

	callbackMessage := []byte("received!")
	_, err = conn.WriteToUDP(callbackMessage, addr)

	if err != nil {
		log.Println(err)
	}
}

func initializeMessaging() *broaker.AmqpClient {
	mc := &broaker.AmqpClient{}
	mc.Connect("amqp://easymile:easymile@localhost:5672/")
	return mc
}

func extractDataFromStringBuffer(command string) (*models.Trip, error) {
	re := regexp.MustCompile(`[,]`)
	result := re.Split(command, -1)

	if len(result) < 24 {
		return nil, errors.New("Parse: erro ao tentar carregar as informações da trip")
	}

	Prefix := result[0]
	CRC := result[1]
	Length, err := strconv.ParseUint(result[2], 32, 32)
	SeqID, err := strconv.ParseInt(result[3], 32, 32)
	UnitID, err := strconv.ParseFloat(result[4], 64)
	GPSDateTime, err := strconv.ParseFloat(result[5], 64)
	RTCDateTime, err := strconv.ParseFloat(result[6], 64)
	PositionSendingDateTime, err := strconv.ParseFloat(result[7], 64)
	Longitude, err := strconv.ParseFloat(result[8], 64)
	Latitude, err := strconv.ParseFloat(result[9], 64)
	Heading, err := strconv.ParseInt(result[10], 16, 16)
	ReportID, err := strconv.ParseInt(result[11], 16, 16)
	Odometer, err := strconv.ParseInt(result[12], 16, 16)
	GPSHDOP, err := strconv.ParseInt(result[13], 16, 16)
	InputStatus, err := strconv.ParseInt(result[14], 16, 16)
	GPSVSSVehicleSpeed, err := strconv.ParseFloat(result[15], 64)
	OutputStatus, err := strconv.ParseFloat(result[16], 64)
	AnalogInputValue, err := strconv.ParseFloat(result[17], 64)
	DriverID, err := strconv.ParseFloat(result[18], 64)
	FirstTemperatureSensorValue, err := strconv.ParseFloat(result[19], 64)
	SecondTemperatureSensorValue, err := strconv.ParseFloat(result[20], 64)
	TextMessage := result[21]
	VINNumber := result[22]
	FuelLevelInPercentage, err := strconv.ParseFloat(result[23], 64)
	FuelUsedInZeroDotOneliter, err := strconv.ParseFloat(result[24], 64)

	if err != nil {
		fmt.Println(err)
	}

	trip := &models.Trip{
		Header: models.Header{
			Prefix: Prefix,
			CRC:    CRC,
			Length: uint8(Length),
			SeqID:  int(SeqID),
			UnitID: float64(UnitID),
		},
		GPSDateTime:                  float64(GPSDateTime),
		RTCDateTime:                  float64(RTCDateTime),
		PositionSendingDateTime:      float64(PositionSendingDateTime),
		Longitude:                    float64(Longitude),
		Latitude:                     float64(Latitude),
		Heading:                      int(Heading),
		ReportID:                     int(ReportID),
		Odometer:                     int(Odometer),
		GPSHDOP:                      int(GPSHDOP),
		InputStatus:                  int(InputStatus),
		GPSVSSVehicleSpeed:           int(GPSVSSVehicleSpeed),
		OutputStatus:                 int(OutputStatus),
		AnalogInputValue:             int(AnalogInputValue),
		DriverID:                     float64(DriverID),
		FirstTemperatureSensorValue:  float64(FirstTemperatureSensorValue),
		SecondTemperatureSensorValue: float64(SecondTemperatureSensorValue),
		TextMessage:                  TextMessage,
		VINNumber:                    VINNumber,
		FuelLevelInPercentage:        float64(FuelLevelInPercentage),
		FuelUsedInZeroDotOneliter:    float64(FuelUsedInZeroDotOneliter),
	}

	if err != nil {
		return nil, errors.New("Parse: erro ao fazer parse das insformações")
	}

	return trip, nil
}

func main() {
	hostName := "localhost"
	portNum := "8080"
	service := hostName + ":" + portNum

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	broakerConn := initializeMessaging()

	fmt.Println("UDP server up and listening on port: " + portNum)

	defer ln.Close()

	for {
		handleUDPConnection(ln, broakerConn)
	}

	handleSigterm(func() {
		broakerConn.Close()
	})
}

func handleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}
