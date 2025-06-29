package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type SuccessfullResponse struct {
    From	string	`json:"from"`
    To		string	`json:"to"`
    Input	string	`json:"input"`
    Result	float64	`json:"result"`
}

type ErrorResponse struct {
    Code 	int
    Message	string	`json:"error_message"`
}

type server struct {
    http *http.Server
    Handler http.Handler
}

// allowed units
var units = map[string]float64 {
    "km": 1.0, 
    "mile": 1.60934,
}

func (s *server) Run() { 
    s.http = &http.Server {
	Addr: "8081",
	Handler: s.Handler,
    }

    if err:= s.http.ListenAndServe(); err != nil {
	log.Println("Cannot start HTTP server.", err.Error())
    }
}

func main() {
    handler := createHandler()
    server := &server {
	Handler: handler,
    }
    server.Run()
}

func createHandler() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/convert", handleConvert)   
    mux.HandleFunc("/healthy", func(w http.ResponseWriter, _ *http.Request){
	w.Write([]byte("{\"healthy\": \"true\"}"))
    })   

    return mux
}

func handleConvert(w http.ResponseWriter, req *http.Request) {
    // get query params: from, to, value
    queryParams := req.URL.Query()

    // validate query params
    if err := validateParams(queryParams); err != nil {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	rb, _ := json.Marshal(err)
	w.Write(rb)
    }

    // do actual convert
    inputValue, _ := strconv.ParseFloat(queryParams.Get("value"), 64)
    convertRes := convert(queryParams.Get("from"), queryParams.Get("to"), inputValue)

    // construct json
    response, _ := json.Marshal(SuccessfullResponse {
	From: queryParams.Get("from"),
	To: queryParams.Get("to"),
	Input: queryParams.Get("value"),
	Result: convertRes,
    } )

    // send response
    w.Header().Set("Content-Type", "application/json")
    w.Write(response)
}

func validateParams(queryParams url.Values) *ErrorResponse {
    // check presense
    if queryParams.Get("from") == "" ||
       queryParams.Get("to") == "" ||
       queryParams.Get("value") == "" {

	return &ErrorResponse {
	    Code: http.StatusBadRequest,
	    Message: "missing required query params",
	}
    }

    // check allowed convert values
    if _, ok := units[queryParams.Get("from")]; !ok {
	return &ErrorResponse {
	    Code: http.StatusBadRequest,
	    Message: fmt.Sprintf("not supported convert value '%s'", queryParams.Get("from")),
	}
    }

    if _, ok := units[queryParams.Get("to")]; !ok {
	return &ErrorResponse {
	    Code: http.StatusBadRequest,
	    Message: fmt.Sprintf("not supported convert value '%s'", queryParams.Get("to")),
	}
    }

    // validate input value
    _, err := strconv.ParseFloat(queryParams.Get("value"), 32)
    if err != nil {
	return &ErrorResponse {
	    Code: http.StatusBadRequest,
	    Message: fmt.Sprint("incorrect input value"),
	}
    }

    return nil
}

func convert(from, to string, value float64) float64 {
    var res float64
    switch {
	case from == "km" && to == "mile":
	    res = value * units["mile"]
	case from == "mile" && to == "km": 
	    res = value / units["mile"]
	default:
	    res = value
    }
    return res
}
