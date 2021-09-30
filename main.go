package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/shotis/webp-go"
)

var (
	ConfigPath     = flag.String("config", "config.json", "Path to config.json")
	ImageDirectory = flag.String("images", "images", "Location where we save the downloaded images")
	MetadataPath   = flag.String("metadata", "metadata.json", "Location where we save the scraped metadata")

	WebPConfig = &webp.Config{
		Lossless: true,
		Quality:  80.0,
	}

	config *Config
)

func init() {
	flag.Parse()

	file, err := os.Open(*ConfigPath)

	if err != nil {
		log.Fatalln(err)
	}

	c, err := ReadConfig(file)

	if err != nil {
		log.Fatalln(err)
	}

	config = c

	os.Mkdir(*ImageDirectory, os.ModePerm)
}

func main() {
	// So we basically need to make a bunch of HTTP requests heh.
	var (
		metadataBaseURL = fmt.Sprintf(config.BaseURI, config.MetadataAddress)
		imagesBaseURL   = fmt.Sprintf(config.BaseURI, config.ImageAddress)

		metadataChannel = make(chan *Metadata, 200)

		metadata map[string]*Metadata = map[string]*Metadata{}
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			close(metadataChannel)
		}()
		for i := 0; i < 9200; i++ {

			res, err := http.Get(metadataBaseURL + fmt.Sprintf("%d.json", i))
			if err != nil {
				log.Fatalln(err)
			}

			md, err := ReadMetadata(res.Body)

			if err != nil {
				log.Fatalln(err)
			}

			for _, attr := range md.Attributes {
				if attr.TraitType == "hand" && strings.Contains(attr.Value, "popsicle") {
					log.Println("found a popsicle")
					metadataChannel <- md
					metadata[strconv.Itoa(i)] = md
				}
			}
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			metadata, closed := <-metadataChannel

			fmt.Println(metadata, closed)
			if !closed {
				break
			}

			tokenIdAttribute := metadata.Attributes[len(metadata.Attributes)-1]
			fileName := fmt.Sprintf("%s.jpg", tokenIdAttribute.Value)
			savedFilePath := filepath.Join(*ImageDirectory, fileName)

			res, err := http.Get(imagesBaseURL + fileName)

			if err != nil {
				log.Fatalln(err)
			}

			if res.Header.Get("Content-Type") == "image/jpeg" {
				b, _ := ioutil.ReadAll(res.Body)
				ioutil.WriteFile(savedFilePath, b, os.ModePerm)
			}

		}
	}()

	wg.Wait()

	if j, err := json.Marshal(metadata); err != nil {
		log.Fatalln(err)
	} else {
		ioutil.WriteFile(*MetadataPath, j, os.ModePerm)
	}
}
