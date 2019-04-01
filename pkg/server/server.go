package server

type BaremetalIntrospectionData struct {
	something     string
	somethingElse string
	anotherThing  string
}

type BaremetalConfig struct {
	Url string `json:"url"`
}

type BaremetalServer struct {
}

func (bms *BaremetalServer) GetConfigUrl(data BaremetalIntrospectionData) (*BaremetalConfig, error) {
	return &BaremetalConfig{Url: "TODO"}, nil
}
