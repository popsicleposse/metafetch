package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Config struct {
	// the BAFY address for the metadata location
	MetadataAddress string
	// the BAFY address for the image locations
	ImageAddress string
	// The base URI for the IPFS node
	BaseURI string
}

func ReadConfig(r io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var config Config

	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, err
}
