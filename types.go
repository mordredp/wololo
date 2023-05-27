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

// Data is list of Computer objects defined in JSON config file.
type Data struct {
	Devices []Device `json:"devices"`
}

// Config represents a configuration object to initialize this application.
type Config struct {
	IP                string `json:"ip" env:"WOLOLOIP" env-default:"0.0.0.0"`
	Port              int    `json:"port" env:"WOLOLOPORT" env-default:"8089"`
	BCastIP           string `json:"bcastip" env:"WOLOLOBCASTIP" env-default:"192.168.1.255:9"`
	MaxSessionSeconds int    `json:"maxsessionseconds" env-default:"600"`
	StaticPass        string `json:"staticpass"`
	LDAPAddr          string `json:"ldapaddr"`
	LDAPBaseDN        string `json:"ldapbasedn"`
	LDAPBindUser      string `json:"ldapbinduser"`
	LDAPBindPass      string `json:"ldapbindpass"`
	LDAPClassValue    string `json:"ldapclassvalue"`
	LDAPIdKey         string `json:"ldapidkey"`
}
