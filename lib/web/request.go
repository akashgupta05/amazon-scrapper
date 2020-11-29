package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ReadBody(r *http.Request) (map[string]interface{}, error) {
	requestBodyMap := make(map[string]interface{})

	bodyByte, err := ReadBodyBytes(r)
	if err != nil {
		return requestBodyMap, err
	}

	requestBodyMap, err = unmarshalRequestBody(bodyByte)
	if err != nil {
		log.Warn("Error decoding request json", err.Error(), r.Header, r.URL, string(bodyByte))
	}
	return requestBodyMap, err
}

func ReadBodyBytes(r *http.Request) ([]byte, error) {
	if r.ContentLength == 0 {
		return []byte(``), nil
	}

	bodyByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn("Error reading request body", err.Error())
		return []byte(``), err
	}
	return bodyByte, nil
}

func unmarshalRequestBody(body []byte) (map[string]interface{}, error) {
	requestBodyMap := make(map[string]interface{})
	b := bytes.NewBuffer(body)
	decoder := json.NewDecoder(b)
	decoder.UseNumber()
	err := decoder.Decode(&requestBodyMap)
	return requestBodyMap, err
}
