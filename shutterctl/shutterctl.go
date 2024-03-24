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

// Setting up all variables, that are needed in multiple function
var lux int
var status int
var tlSensorAddress string

func main() {
	// Get the hostname of the tl-sensor container
        // from an enviromental variable
	tlSensorAddress = os.Getenv("TL_SENSOR_ADDRESS")
	if tlSensorAddress == "" {
		log.Print("Enviromental Variabel: TL_SENSOR_ADDRESS not set")
	}
	
        // Start the function that gets the current lux-value
	go getLux()

	// Create a handler to serve the current shutterctl status
	http.HandleFunc("/d-status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", status)
	})
	
	// Start a server
		http.ListenAndServe(":8080", nil)
}

// This function gets the current lux-value
func getLux() {
	for{
                // Request the current lux-value using the hostname and port
                // from the environmental variabel
		respLux, err := http.Get(fmt.Sprintf("http://%s/lux", tlSensorAddress))
		if err != nil {
			log.Print(err)
		}

	// Make sure the body is closed, when the function returns
                // TODO: Currently the function never returns
		defer respLux.Body.Close()
	

	                // Read the body of the response to our request
		bodyLux, err := io.ReadAll(respLux.Body)
		if err != nil {
			log.Print(err)
		}
		
		// set status to a string, that is the converted
                // []byte from our response body
		lux, err = strconv.Atoi(string(bodyLux))		
		if err != nil {
			log.Print(err)
		}
		
		// call the shutterctl function
		shutterctl()

		// set a delay
		time.Sleep(3 * time.Second)
	}
}

// This function controls the shuttercontrol status
func shutterctl() {
	fmt.Println("Last requested lux-value:", lux)
	// If Illuminance greate than 32000 (direct sunlight)
	// set status 1 (shutter closed)
	// Else set status 0 (shutter open)
	if lux < 32000 {
		status = 1
	} else {
		status = 0
	} 
	fmt.Println("Current shutter-status:", status) 
}

