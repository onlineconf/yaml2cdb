package main

import (
	"flag"
	"log"
	"os"

	"github.com/colinmarc/cdb"
	"github.com/onlineconf/yaml2cdb/internal/converter"
	"gopkg.in/yaml.v3"
)

var (
	yamlFile = flag.String("in", "", "input yaml file")
	cdbFile  = flag.String("out", "", "result cdb file")
)

func main() {
	flag.Parse()

	if *yamlFile == "" {
		log.Fatal("yaml should be provided")
	}

	if *cdbFile == "" {
		log.Fatal("cdb should be provided")
	}

	content, err := os.ReadFile(*yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]any

	if err = yaml.Unmarshal(content, &data); err != nil {
		log.Fatal(err)
	}

	items, err := converter.CDBItems(data, "")
	if err != nil {
		log.Fatal(err)
	}

	writer, err := cdb.Create(*cdbFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		if err = writer.Put([]byte(item.Key), []byte(item.Value)); err != nil {
			log.Fatal(err)
		}
	}

	writer.Close()
}
