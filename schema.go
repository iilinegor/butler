package main

// Artef is artefactory config schema
type Artef struct {
	Name    string `json:"name"`
	Ver     int    `json:"ver"`
	Tok     string `json:"tok"`
	GitPath string `json:"gitPath" bson:"gitPath"`
	Bin     string `json:"bin"`
	Port    string `json:"port"`
}

// Squad is registry for runners
type Squad struct {
	Name string `json:"name"`
	Ips  Ips    `hson:"ips"`
	Ms   []Ms   `json:"ms"`
}

// Ips is v4 & v6 container
type Ips struct {
	V4 string `json:"v4"`
	V6 string `json:"v6"`
}

// Ms is microservice declaration type
type Ms struct {
	ID    string `json:"name"`
	Bin   string `json:"bin"`
	Param string `json:"param"`
	Port  string `json:"port"`
}
