package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

func loadData() {

	devicesFile, err := os.Open("devices.json")
	if err != nil {
		log.Fatalf("error loading definitions: %q", err)
	}
	devicesDecoder := json.NewDecoder(devicesFile)
	err = devicesDecoder.Decode(&appData)
	if err != nil {
		log.Fatalf("error decoding definitions: %q", err)
	}
	log.Printf("%d definitions loaded", len(appData.Devices))
}

func saveData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var result Response

	log.Printf("saving application data ...")
	err := json.NewDecoder(r.Body).Decode(&appData)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		result.Success = false
		result.Message = "could not save the data" + err.Error()
		result.ErrorObject = err
		log.Printf("error decoding/saving data")
	} else {
		file, _ := os.OpenFile("devices.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		encoder.Encode(appData)

		result.Success = true
		result.Message = "data saved, " + strconv.Itoa(len(appData.Devices)) + " definitions"
		log.Printf("data saved")
	}
	json.NewEncoder(w).Encode(result)

}

func getData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appData)
	log.Printf("data sent")

}
