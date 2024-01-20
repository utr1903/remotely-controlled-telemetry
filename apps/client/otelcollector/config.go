package otelcollector

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type otelCollectorConfig struct {
	Receivers struct {
		Otlp struct {
			Protocols struct {
				Grpc struct {
					Endpoint string `yaml:"endpoint"`
				} `yaml:"grpc"`
			} `yaml:"protocols"`
		} `yaml:"otlp"`
	} `yaml:"receivers"`
	Exporters struct {
		File struct {
			Path string `yaml:"path"`
		} `yaml:"file"`
	} `yaml:"exporters"`
	Service struct {
		Pipelines struct {
			Metrics struct {
				Receivers []string `yaml:"receivers"`
				Exporters []string `yaml:"exporters"`
			} `yaml:"metrics"`
		} `yaml:"pipelines"`
	} `yaml:"service"`
}

type otelCollectorConfigGenerator struct {
}

func Test() {
	occg := newOtelCollectorConfigGenerator()
	occg.generate()
}

func newOtelCollectorConfigGenerator() *otelCollectorConfigGenerator {
	return &otelCollectorConfigGenerator{}
}

func (o *otelCollectorConfigGenerator) generate() {

	cfg := &otelCollectorConfig{
		Receivers: struct {
			Otlp struct {
				Protocols struct {
					Grpc struct {
						Endpoint string "yaml:\"endpoint\""
					} "yaml:\"grpc\""
				} "yaml:\"protocols\""
			} "yaml:\"otlp\""
		}{
			Otlp: struct {
				Protocols struct {
					Grpc struct {
						Endpoint string "yaml:\"endpoint\""
					} "yaml:\"grpc\""
				} "yaml:\"protocols\""
			}{
				Protocols: struct {
					Grpc struct {
						Endpoint string "yaml:\"endpoint\""
					} "yaml:\"grpc\""
				}{
					Grpc: struct {
						Endpoint string "yaml:\"endpoint\""
					}{
						Endpoint: "localhost:4317",
					},
				},
			},
		},
		Exporters: struct {
			File struct {
				Path string "yaml:\"path\""
			} "yaml:\"file\""
		}{
			File: struct {
				Path string "yaml:\"path\""
			}{
				Path: "./bin/log",
			},
		},
		Service: struct {
			Pipelines struct {
				Metrics struct {
					Receivers []string "yaml:\"receivers\""
					Exporters []string "yaml:\"exporters\""
				} "yaml:\"metrics\""
			} "yaml:\"pipelines\""
		}{
			Pipelines: struct {
				Metrics struct {
					Receivers []string "yaml:\"receivers\""
					Exporters []string "yaml:\"exporters\""
				} "yaml:\"metrics\""
			}{
				Metrics: struct {
					Receivers []string "yaml:\"receivers\""
					Exporters []string "yaml:\"exporters\""
				}{
					Receivers: []string{
						"otlp",
					},
					Exporters: []string{
						"file",
					},
				},
			},
		},
	}

	// Marshal the struct into YAML format
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling YAML: %v\n", err)
		return
	}

	fmt.Println(string(yamlData))

	// Create the YAML file
	file, err := os.Create("./bin/otel-config-test.yaml")
	if err != nil {
		fmt.Printf("Error creating YAML file: %v\n", err)
		return
	}
	defer file.Close()

	// Write the YAML data to the file
	_, err = file.Write(yamlData)
	if err != nil {
		fmt.Printf("Error writing YAML data to file: %v\n", err)
		return
	}
}
