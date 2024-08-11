package utils

import (
	"ComicCollector/main/backend/utils/env"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func ConvertToDateTime(layout string, input time.Time) primitive.DateTime {
	timezone := env.Timezone
	currentTime := input.Format(layout)

	//set local timezone and apply the desired layout
	parsedTime, err := time.ParseInLocation(layout, currentTime, timezone)
	if err != nil {
		log.Println("Error parsing time:", err)
	}

	return primitive.NewDateTimeFromTime(parsedTime)
}

func ConvertToLocalTime(input primitive.DateTime) time.Time {
	//set local timezone
	timezone := env.Timezone
	return input.Time().In(timezone)
}

func GetCurrentLocalTime() time.Time {
	timezone := env.Timezone
	currentTime := time.Now().Format(time.DateTime)

	//set local timezone and apply the desired layout
	parsedTime, err := time.ParseInLocation(time.DateTime, currentTime, timezone)
	if err != nil {
		log.Println("Error parsing time:", err)
	}

	return parsedTime
}

func GetCurrentTimeFormatted() string {
	currentTime := GetCurrentLocalTime()

	return currentTime.Format(time.RFC1123)
}
