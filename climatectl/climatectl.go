package main

import (
	"fmt"
	"net/http"
	"log"
	"io"
	"time"
	"strconv"
	"os"
)
// set variabels
var temp int
var status int
var tlSensorAddress string


func main() {
	// Get the hostname of the tl-sensor container
	// from an enviromental variable
	tlSensorAddress = os.Getenv("TL_SENSOR_ADDRESS")
	if tlSensorAddress == "" {
		log.Fatal("Enviromental Variabel: TL_SENSOR_ADDRESS not set")
	}
	// Print it out, for troubleshooting
	fmt.Println(tlSensorAddress)

	// Start the function that gets the current temperature
	go getTemperature()
	
	// Create a handler to serve the status 
	http.HandleFunc("/cc-status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", status)
	})
	
	// Start the server
	http.ListenAndServe(":8080", nil)
}

// This function gets the current temperature and calls the climatectl function
func getTemperature() {
	for{
		// Request the current temperature status using the hostname and port
		// from the environmental variabel 
		respTemp, err := http.Get(fmt.Sprintf("http://%s/temp", tlSensorAddress))
		if err != nil {
			log.Print(err)
		}
		
		// Make sure the body is closed, when the function returns
		// TODO: Currently the function never returns
		defer respTemp.Body.Close()
		
        	// Read the body of the response to our request
		bodyTemp, err := io.ReadAll(respTemp.Body)
		if err != nil {
			log.Print(err)
		}

		// set status to a string, that is the converted
		// []byte from our response body 
		temp, err = strconv.Atoi(string(bodyTemp))		
		if err != nil {
			log.Print(err)
		}
		
		// call the climatectl function
		climatectl()

		// set a delay
		time.Sleep(3 * time.Second)
	}
}

// This function sets the status of climatectl
func climatectl () {
	// Print the current temperature the function is working with
	fmt.Println("Current used temperature is:", temp)
	
	// If the temperature below 17 degrees Celcius set status to 1 (= heating)
	// Else if the temperature is above 23 set status to -1 (cooling)
	// Else set the status to 0 (off
	if temp < 17 {
		status = 1
	} else if temp > 23 {
		status = -1
	} else {
		status = 0
	}
	fmt.Println("Current climatectl-status:", status)
}

