package main

// Response represents the API response.
type Response struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	ErrorObject error  `json:"error"`
}

// Device represents a Host to be woken up.
type Device struct {
	Name        string `json:"name"`
	MAC         string `json:"mac"`
	BroadcastIP string `json:"ip"`
}

// AppData is list of Computer objects defined in JSON config file.
type AppData struct {
	Devices []Device `json:"devices"`
}

// AppConfig represents a configuration object to initialize this application.
type AppConfig struct {
	IP           string `json:"ip" env:"WOLOLOIP" env-default:"0.0.0.0"`
	Port         int    `json:"port" env:"WOLOLOPORT" env-default:"8089"`
	BCastIP      string `json:"bcastip" env:"WOLOLOBCASTIP" env-default:"192.168.1.255:9"`
	LDAPAddr     string `json:"ldapaddr"`
	LDAPBaseDN   string `json:"ldapbasedn"`
	LDAPBindUser string `json:"ldapbinduser"`
	LDAPBindPass string `json:"ldapbindpass"`
}
