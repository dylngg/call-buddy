package test_server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ReportRequest struct {
	ResponseStatus int
}

type ReportResponse struct {
	Method string
	Body   string
	// TODO AH: Just make this map[string]string
	Headers map[string]string
	Url     string
}

func TestServer() *http.Server {
	toReturn := &http.Server{}

	toReturn.Handler = TestServerMux()

	return toReturn
}

//TestServeMux returns a servemux that returns a report for ALL paths with an optional
func TestServerMux() *http.ServeMux {
	toReturn := http.NewServeMux()

	toReturn.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		statusToWrite := 200
		var readBytes []byte

		// Read all the body bytes or return an error
		if tempRead, readError := ioutil.ReadAll(request.Body); readError == nil {
			readBytes = tempRead
		} else {
			writer.WriteHeader(509)
			writer.Write([]byte("SERVER FACED ERROR READING BODY: " + readError.Error()))
			return
		}

		//report := ReportRequest{}
		// Unmarshal the body, look for the status code in a json field
		//if request.Header.Get("Content-Type") == "application/json" && json.Unmarshal(readBytes, &report) == nil {
		//	log.Println("Requester's got a status code they want returned!: " + strconv.Itoa(report.ResponseStatus))
		//	statusToWrite = report.ResponseStatus
		//}

		if rStat := request.Header.Get("Requested-Status"); rStat != "" {
			if intVal, intErr := strconv.Atoi(rStat); intErr == nil {
				statusToWrite = intVal
			} else {
				log.Println("Requested-Status was NOT an integer!")
			}
		}

		// Fill out response object
		responseObject := ReportResponse{}

		// Take care of headers
		responseObject.Headers = make(map[string]string, len(request.Header))

		for k, v := range request.Header {
			responseObject.Headers[k] = v[0]
		}

		// Take care of the other things
		responseObject.Body = string(readBytes)
		responseObject.Method = request.Method
		responseObject.Url = request.URL.Path

		// Marshal the body
		var bodyToWrite []byte
		if marshaled, marshalErr := json.MarshalIndent(responseObject, "", "\t"); marshalErr == nil {
			writer.WriteHeader(statusToWrite)
			bodyToWrite = marshaled
		} else {
			writer.WriteHeader(509)
			bodyToWrite = []byte("SERVER FACED ERROR: " + marshalErr.Error())
		}

		if _, writeErr := writer.Write(bodyToWrite); writeErr != nil {
			log.Println("ERROR; Couldn't write the response!: " + writeErr.Error())
		}
	})

	return toReturn
}
