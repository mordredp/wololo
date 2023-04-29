package main

//HTTPResponseObject Data structure for sending the API call status
type HTTPResponseObject struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ErrorObject error  `json:"error"`
}

// Device represents a Computer Object
type Device struct {
	Name        string `json:"name"`
	Mac         string `json:"mac"`
	BroadcastIP string `json:"ip"`
}

// AppData is list of Computer objects defined in JSON config file
type AppData struct {
	Devices []Device `json:"devices"`
}

// AppConfig represents a configuration object to initialize this application
type AppConfig struct {
	IP      string `json:"ip" env:"WOLOLOIP" env-default:"0.0.0.0"`
	Port    int    `json:"port" env:"WOLOLOPORT" env-default:"8089"`
	VDir    string `json:"vdir" env:"WOLOLOVDIR" env-default:"/wolweb"`
	BCastIP string `json:"bcastip" env:"WOLOLOBCASTIP" env-default:"192.168.1.255:9"`
}
