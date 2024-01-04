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
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	"github.com/andrewkroh/go-ingest-node"
)

func main() {
	outputFile := flag.String("o", "-", "Output file name. Defaults to stdout.")
	flag.Parse()
	log.SetPrefix("jsonschema ")

	var o io.Writer = os.Stdout
	if *outputFile != "-" {
		os.MkdirAll(filepath.Dir(*outputFile), 0o700)

		outFile, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal()
		}
		defer outFile.Close()

		o = outFile
	}

	var r jsonschema.Reflector

	// WARNING: This assumes it is always executed from the root of this module's path.
	if err := r.AddGoComments("github.com/andrewkroh/go-ingest-node/internal/generated", "../.."); err != nil {
		log.Fatalf("Failed to load comments: %v", err)
	}

	s := r.Reflect(&ingestnode.Pipeline{})
	s.Description = "Elasticsearch Ingest Node Pipeline Schema"

	enc := json.NewEncoder(o)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	err := enc.Encode(s)
	if err != nil {
		log.Fatalf("Failed to encode schema: %v", err)
	}
	log.Println("Done. Generated to", *outputFile)
}
