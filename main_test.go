package main

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestValidateOkParams(t *testing.T) {
    from := "km"
    to := "mile"
    value := "11"

    testParams := map[string][]string {
	"from":{ from },
	"to":  { to },
	"value": { value },
    }

    act := validateParams(url.Values(testParams))
    if act != nil {
	t.Errorf("expected %v, got %v", nil, act)
    }
}

func TestValidateMissingParams(t *testing.T) {
    from := ""
    to := ""
    value := "5"

    testParams := map[string][]string {
	"from":{ from },
	"to":  { to },
	"value": { value },
    }

    exp := &ErrorResponse {
	Code: http.StatusBadRequest,
	Message: "missing required query params",
    }
    act := validateParams(url.Values(testParams))

    if !reflect.DeepEqual(act, exp) {
	t.Errorf("expected %+v, got %+v", exp, act)
    }
}

func TestValidateNotSupportedFromParam(t *testing.T) {
    from := "mile"
    to := "yard"
    value := "5"

    testParams := map[string][]string {
	"from":{ from },
	"to":  { to },
	"value": { value },
    }

    exp := &ErrorResponse {
	Code: http.StatusBadRequest,
	Message: "not supported convert value 'yard'",
    }
    act := validateParams(url.Values(testParams))

    if !reflect.DeepEqual(act, exp) {
	t.Errorf("expected %+v, got %+v", exp, act)
    }
}

func TestValidateNotSupportedToParam(t *testing.T) {
    from := "mile"
    to := "feet"
    value := "5"

    testParams := map[string][]string {
	"from":{ from },
	"to":  { to },
	"value": { value },
    }

    exp := &ErrorResponse {
	Code: http.StatusBadRequest,
	Message: "not supported convert value 'feet'",
    }
    act := validateParams(url.Values(testParams))

    if !reflect.DeepEqual(act, exp) {
	t.Errorf("expected %+v, got %+v", exp, act)
    }
}

func TestConvertKmToMile(t *testing.T) {
    from := "km"
    to := "mile"
    value := 5.0

    exp := 8.0467 
    act := convert(from, to, value)

    if act != exp {
	t.Errorf("extected %f, got %f", exp, act)
    }
}

func TestConvertMileToKm(t *testing.T) {
    from := "mile"
    to := "km"
    value := 8.0467

    exp := 5.0
    act := convert(from, to, value)

    if act != exp {
	t.Errorf("extected %f, got %f", exp, act)
    }
}

func TestConvertKmToKm(t *testing.T) {
    from := "km"
    to := "km"
    value := 11.0 

    exp := 11.0
    act := convert(from, to, value)

    if act != exp {
	t.Errorf("extected %f, got %f", exp, act)
    }
}

func TestConvertMileToMile(t *testing.T) {
    from := "mile"
    to := "mile"
    value := 11.0 

    exp := 11.0
    act := convert(from, to, value)

    if act != exp {
	t.Errorf("extected %f, got %f", exp, act)
    }
}
