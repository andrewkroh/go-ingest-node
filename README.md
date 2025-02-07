# go-ingest-node

`go-ingest-node` provides a [Go data model][godoc] and JSON Schema for 
Elasticsearch [Ingest Node][ingest] pipelines. The model is generated from the
[elasticsearch-specification][spec].

[godoc]: https://pkg.go.dev/github.com/andrewkroh/go-ingest-node
[ingest]: https://www.elastic.co/guide/en/elasticsearch/reference/current/processors.html
[spec]: https://github.com/elastic/elasticsearch-specification

The Go types can be used for marshaling to/from JSON and YAML.

The JSON Schema file can be imported into popular editors and IDEs to provide
validation and autocompletion when developing Elasticsearh ingest node
pipelines.