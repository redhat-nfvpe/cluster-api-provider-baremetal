package server

type BaremetalIntrospectionData struct {
	something     string
	somethingElse string
	anotherThing  string
}

type BaremetalConfig struct {
	url string
}

type BaremetalServer struct {
}

func (bms *BaremetalServer) GetConfigUrl(data BaremetalIntrospectionData) (*BaremetalConfig, error) {
	return &BaremetalConfig{url: "TODO"}, nil
}
