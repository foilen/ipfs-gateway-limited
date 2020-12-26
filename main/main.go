package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

var rootConfig *RootConfiguration
var nextRequestID *uint64 = new(uint64)

func main() {

	// Get the configuration
	args := os.Args[1:]

	if len(args) != 1 {
		panic("You need to specify the config file to use")
	}

	var err error
	rootConfig, err = getRootConfiguration(args[0])
	if err != nil {
		panic(err)
	}

	// Start the reverse proxy
	log.Println("Starting local reverse proxy on port", rootConfig.Port)

	http.HandleFunc("/", handler)
	log.Print(http.ListenAndServe(fmt.Sprintf(":%v", rootConfig.Port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	requestID := "[" + strconv.FormatUint(atomic.AddUint64(nextRequestID, 1), 10) + "]"

	// Get the details
	hostname := strings.Split(r.Host, ":")[0]
	path := r.URL.Path

	log.Println(requestID, "Host is", hostname, "; path is", path)
	log.Println(requestID, "Browser request", r)

	// Get the mapping and check if valid
	mappedPath := rootConfig.Mapping[hostname]
	if mappedPath == "" {
		log.Println(requestID, "[ERROR] Unknown host")
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprint(w, "The hostname ", hostname, " is not mapped on this service")
		return
	}
	log.Println(requestID, "Is mapped to", mappedPath)

	// Call the Gateway - Browser -> Gateway
	var err error
	gatewayRequest := r.Clone(r.Context())
	gatewayRequest.RequestURI = ""
	gatewayRequest.Host = strings.Split(rootConfig.LocalGatewayURL, "://")[1]
	gatewayRequest.URL, err = url.Parse(rootConfig.LocalGatewayURL + mappedPath + path)
	if err != nil {
		log.Println(requestID, "[ERROR] Parsing URL", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(requestID, "Gateway request", gatewayRequest)
	gatewayResponse, err := http.DefaultClient.Do(gatewayRequest)
	if err != nil {
		log.Println(requestID, "[ERROR] Calling the Gateway", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer gatewayResponse.Body.Close()
	log.Println(requestID, "Gateway response", gatewayResponse)

	// Copy the headers - Gateway -> Browser
	for k, vv := range gatewayResponse.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	// When HTML, replace all paths
	if w.Header().Get("Content-Type") == "text/html" {
		log.Println(requestID, "Is an HTML page. Replacing all paths")

		buf := new(bytes.Buffer)
		buf.ReadFrom(gatewayResponse.Body)
		bodyStr := buf.String()

		newBody := strings.ReplaceAll(bodyStr, mappedPath, "")
		gatewayResponse.Body = ioutil.NopCloser(strings.NewReader(newBody))
	}

	// Send the response
	w.WriteHeader(gatewayResponse.StatusCode)
	copied, err := io.Copy(w, gatewayResponse.Body)
	if err != nil {
		log.Println(requestID, "[ERROR] Copying the body", err)
		return
	}
	log.Println(requestID, "Browser response", w)

	// Log the outcome
	log.Println(requestID, "[OK] Sent", copied, "bytes in the body")
}
