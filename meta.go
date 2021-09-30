package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

type Metadata struct {
	//The name of the item
	Name string
	// The description of the item
	Desc string `json:"description"`
	// The image URI
	Img        string `json:"image"`
	Attributes []Attribute
}

func ReadMetadata(r io.Reader) (*Metadata, error) {
	data, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	var metadata Metadata

	if err = json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}
	return &metadata, err
}
