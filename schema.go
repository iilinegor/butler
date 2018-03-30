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
	Ms   []Ms   `json:"ms"`
}

// Ms is microservice declaration type
type Ms struct {
	Name string `json:"name"`
	Port string `json:"port"`
}
