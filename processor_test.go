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

package ingestnode

import (
	"encoding/json"
	"strings"
	"testing"
)

const pipeline = `
{
  "description": "Parse Common Log Format.",
  "processors": [
    {
      "grok": {
        "description": "Extract fields from 'message'",
        "field": "message",
        "patterns": [
          "%{IPORHOST:source.ip} %{USER:user.id} %{USER:user.name} \\[%{HTTPDATE:@timestamp}\\] \"%{WORD:http.request.method} %{DATA:url.original} HTTP/%{NUMBER:http.version}\" %{NUMBER:http.response.status_code:int} (?:-|%{NUMBER:http.response.body.bytes:int}) %{QS:http.request.referrer} %{QS:user_agent}"
        ]
      }
    },
    {
      "date": {
        "description": "Format '@timestamp' as 'dd/MMM/yyyy:HH:mm:ss Z'",
        "field": "@timestamp",
        "formats": [
          "dd/MMM/yyyy:HH:mm:ss Z"
        ]
      }
    },
    {
      "geoip": {
        "description": "Add 'source.geo' GeoIP data for 'source.ip'",
        "field": "source.ip",
        "target_field": "source.geo"
      }
    },
    {
      "user_agent": {
        "if": "ctx.user_agent != null",
        "description": "Extract fields from 'user_agent'",
        "field": "user_agent"
      }
    }
  ],
  "on_failure": [
    {
      "set": {
        "field": "event.kind",
        "value": "pipeline_error"
      }
    }
  ]
}
`

func TestPipelineJSONUnmarshal(t *testing.T) {
	dec := json.NewDecoder(strings.NewReader(pipeline))
	dec.DisallowUnknownFields()

	var p Pipeline
	if err := dec.Decode(&p); err != nil {
		t.Fatal(err)
	}
	if len(p.Processors) != 4 {
		t.Fatal()
	}
	if p.Processors[0].Grok == nil {
		t.Fatal("expected grok processor")
	}
}

func TestPipelineJSONMarshal(t *testing.T) {
	p := &Pipeline{
		Description: ptrTo("Parse Common Log Format."),
		Processors: []ProcessorContainer{
			{
				Grok: &GrokProcessor{
					ProcessorBase: ProcessorBase{
						Description: ptrTo("Extract fields from 'message'"),
					},
					Field: "message",
					Patterns: []GrokPattern{
						"%{IPORHOST:source.ip}",
					},
				},
			},
		},
	}

	want := `{
  "description": "Parse Common Log Format.",
  "processors": [
    {
      "grok": {
        "description": "Extract fields from 'message'",
        "field": "message",
        "patterns": [
          "%{IPORHOST:source.ip}"
        ]
      }
    }
  ]
}`

	got, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != want {
		t.Fatalf("want:\n%s got:\n%s", want, string(got))
	}
}

func ptrTo[T any](v T) *T { return &v }
