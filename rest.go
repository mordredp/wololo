// Rest API Implementations

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// restWakeUpWithDeviceName - REST Handler for Processing URLS /virtualdirectory/apipath/<deviceName>
func wakeUpWithDeviceName(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	//vars := mux.Vars(r)
	//deviceName := vars["deviceName"]
	deviceName := chi.URLParam(r, "deviceName")

	var result HTTPResponseObject
	result.Success = false

	// Ensure deviceName is not empty
	if deviceName == "" {
		result.Message = "device name not specified"
		result.ErrorObject = nil
		w.WriteHeader(http.StatusBadRequest)
		// Devicename is empty
	} else {

		// Get Device from List
		for _, c := range appData.Devices {
			if c.Name == deviceName {

				// We found the Devicename
				if err := SendMagicPacket(c.Mac, c.BroadcastIP, ""); err != nil {
					// We got an internal Error on SendMagicPacket
					w.WriteHeader(http.StatusInternalServerError)
					result.Success = false
					result.Message = "internal error on sending the the magic packet"
					result.ErrorObject = err
				} else {
					// Horray we send the WOL Packet succesfully
					result.Success = true
					result.Message = fmt.Sprintf("sent magic packet to device %s with MAC %s on broadcast IP %s", c.Name, c.Mac, c.BroadcastIP)
					result.ErrorObject = nil
				}
			}
		}

		if !result.Success && result.ErrorObject == nil {
			// We could not find the Devicename
			w.WriteHeader(http.StatusNotFound)
			result.Message = fmt.Sprintf("device name %s could not be found", deviceName)
		}
	}
	json.NewEncoder(w).Encode(result)
}
