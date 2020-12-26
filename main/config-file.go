package main

import (
	"encoding/json"
	"io/ioutil"
)

// Example:
// {
//     "port" : 8888,
//     "localGatewayUrl" : "http://127.0.0.1:8080",
//     "mapping" : {
//         "localhost.foilen-lab.com" : "/ipns/cdn.foilen.com",
//         "localhost2.foilen-lab.com" : "/ipns/k51qzi5uqu5dhuj92m8egzrbx6e0apodpebvs7y4fqe1rc6rxb1hsiwgovl94o"
//     }
// }

// RootConfiguration is the json configuration file
type RootConfiguration struct {
	Port            uint16
	LocalGatewayURL string
	Mapping         map[string]string
}

func getRootConfiguration(path string) (*RootConfiguration, error) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rootConfiguration RootConfiguration
	err = json.Unmarshal(jsonBytes, &rootConfiguration)

	return &rootConfiguration, err
}
