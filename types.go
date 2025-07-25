package main

type Target struct {
	Path    string `json:"path"`
	Exclude string `json:"exclude,omitempty"`
}

type Dependency struct {
	Name string `json:"name"`
	Ver  string `json:"ver"`
}

type PacketConfig struct {
	Name    string        `json:"name"`
	Version string        `json:"ver"`
	Targets []interface{} `json:"targets"`
	Packets []Dependency  `json:"packets"`
}

type UpdateList struct {
	Packages []Dependency `json:"packages"`
}
