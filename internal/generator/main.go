// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"flag"
	"generator/internal/codegen"
	"generator/internal/spec"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fileName := flag.String("f", "", "JSON file containing the elasticsearch-specification. Download from https://github.com/elastic/elasticsearch-specification/blob/main/output/schema/schema.json")
	outputFile := flag.String("o", "-", "Output file name. Defaults to stdout.")
	flag.Parse()

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	model, err := loadModel(*fileName)
	if err != nil {
		log.Fatalf("Failed to load schema model: %v", err)
	}

	var o io.Writer = os.Stdout
	if *outputFile != "-" {
		if err = os.MkdirAll(filepath.Dir(*outputFile), 0o700); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		outFile, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal()
		}
		defer outFile.Close()

		o = outFile
	}

	b := codegen.New(model)

	err = b.BuildCode(o, func(name, inherits spec.TypeName) bool {
		// This is what selects the types to include in generated file.
		// Any dependencies of these types will be included automatically.
		return name.Namespace == "ingest._types" && name.Name == "Pipeline" || inherits.Name == "ProcessorBase"
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done. Generated to", *outputFile)
}

func loadModel(path string) (*spec.Model, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.UseNumber()

	var model spec.Model
	if err = dec.Decode(&model); err != nil {
		return nil, err
	}
	return &model, nil
}
