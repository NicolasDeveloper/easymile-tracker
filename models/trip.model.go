package models

// Header : Unit header of trip request
type Header struct {
	Prefix string  `json:"prefix" bson:"prefix"`
	CRC    string  `json:"crc" bson:"crc"`
	Length uint8   `json:"length" bson:"length"`
	SeqID  int     `json:"seqID" bson:"seqID"`
	UnitID float64 `json:"unitID" bson:"unitID"`
}

// Trip : Trip sensor data
type Trip struct {
	Header                       Header  `json:"header" bson:"header"`
	GPSDateTime                  float64 `json:"gpsDateTime" bson:"gpsDateTime"`
	RTCDateTime                  float64 `json:"rtcDateTime" bson:"rtcDateTime"`
	PositionSendingDateTime      float64 `json:"positionSendingDateTime" bson:"positionSendingDateTime"`
	Longitude                    float64 `json:"longitude" bson:"longitude"`
	Latitude                     float64 `json:"latitude" bson:"latitude"`
	Heading                      int     `json:"heading" bson:"heading"`
	ReportID                     int     `json:"reportID" bson:"reportID"`
	Odometer                     int     `json:"odometer" bson:"odometer"`
	GPSHDOP                      int     `json:"gpshdop" bson:"gpshdop"`
	InputStatus                  int     `json:"inputStatus" bson:"inputStatus"`
	GPSVSSVehicleSpeed           int     `json:"gpsvssVehicleSpeed" bson:"gpsvssVehicleSpeed"`
	OutputStatus                 int     `json:"outputStatus" bson:"outputStatus"`
	AnalogInputValue             int     `json:"analogInputValue" bson:"analogInputValue"`
	DriverID                     float64 `json:"driverID" bson:"driverID"`
	FirstTemperatureSensorValue  float64 `json:"firstTemperatureSensorValue" bson:"firstTemperatureSensorValue"`
	SecondTemperatureSensorValue float64 `json:"secondTemperatureSensorValue" bson:"secondTemperatureSensorValue"`
	TextMessage                  string  `json:"textMessage" bson:"textMessage"`
	VINNumber                    string  `json:"vinNumber" bson:"vinNumber"`
	FuelLevelInPercentage        float64 `json:"fuelLevelInPercentage" bson:"fuelLevelInPercentage"`
	FuelUsedInZeroDotOneliter    float64 `json:"fuelUsedInZeroDotOneliter" bson:"fuelUsedInZeroDotOneliter"`
}
