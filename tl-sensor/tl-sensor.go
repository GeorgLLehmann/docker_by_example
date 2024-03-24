package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"io"
	"log"
	"strconv"
	"os"
)

// Setting up all variables, that are needed in multiple function
// No mutex is nessacary because only one function is writing to it.
var lux int
var tendency int
var n int
var temp int = 20
var status int
var delay = 3
var climateCtlAddress string 

func main() {
	// Get the hostname of the climatectl-container
	climateCtlAddress = os.Getenv("CLIMATECTL_ADDRESS")
	if climateCtlAddress == "" {
		log.Fatal("CLIMATECTL_ADDRESS environment variable not set")
	}
	// Print it out, for troubleshooting
	fmt.Println(climateCtlAddress)
	
	// Start the function, that simulates the temperature and light values
	go generateValues()
	
	// Start the function, that gets the status of the climatectl in order
	// so it can influence the temperature simulation 
	go getCCStatus()
	
	// Create a handler to serve the lux value
	http.HandleFunc("/lux", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", lux)
	})

	// Create a handler to serve the temp value
	http.HandleFunc("/temp", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", temp)
	})

	// Start the server
	http.ListenAndServe(":8080", nil)
}

// This function is generating simulated light and temperature values.
func generateValues() {
	// Let the generation go on indefinitly 
	for {
		// Generates a random value from 1 (darkness) to 100000 (direct sunlight)
		lux = rand.Intn(100000) 
		// Print the generated value 
		fmt.Println("Generated new value:", lux, "lux")
		
		// creating a tendency so the temperature changes externally in the same
		// for 4 cycles until counter n = 0
		if n != 0 {
			// Adjust the temperature using the status of the climatectl
			// in the most direct manner
			temp = temp + tendency + status
			// Decresing the counter
			n = n - 1
		} else {
			// Setting new tendency which will be -1 || 0 || 1
			tendency = rand.Intn(3) - 1
			// Resetting counter
			n = 3
			// Adjust the temperature using the status of the climatectl
			temp = temp + tendency + status
		}

		// Print the temperature
		fmt.Println("New Generated temperature:", temp, "Grad Celcius")
		// Create a delay, so humans can observe the changes
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

// This function getting the current status of the climatectl
func getCCStatus() {
	for{
		// Request the current climatectl status using the hostname and port
		// from the environmental variabel
		respStatus, err := http.Get(fmt.Sprintf("http://%s/cc-status", climateCtlAddress))
		if err != nil {
			log.Print(err)
		} else {
			// Make sure the body is closed, when the function returns
			// TODO: Currently the function never returns
			defer respStatus.Body.Close()
			
			// Read the body of the respnse to our request
			bodyStatus, err := io.ReadAll(respStatus.Body)
			if err != nil {
				log.Print(err)
			}
			
			// set status to a string, that is the converted
			// []byte from our response body
			status, err = strconv.Atoi(string(bodyStatus))		
			if err != nil {
				log.Print(err)
			}
		}
			// Create a delay
			time.Sleep(time.Duration(delay) * time.Second)
		
	}
}
