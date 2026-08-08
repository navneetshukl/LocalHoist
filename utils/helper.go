package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"
)

func GenerateClientID() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GetClientIdFromRoute(route string) (error,string) {
	strings.ReplaceAll(route," ","")
	if route==""{
		return errors.New("client id is empty"),""

	}
	routesArray := strings.Split(route, "/")
	log.Println("Routes Array ", routesArray)
	log.Println("Length Of Array is ",len(routesArray))
	if len(routesArray)!=4{
		return errors.New("client id is invalid"),""
	}
	return nil,routesArray[2]
}
