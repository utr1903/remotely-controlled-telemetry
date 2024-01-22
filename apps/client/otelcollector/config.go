package otelcollector

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
	"gopkg.in/yaml.v3"
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
		Filelog struct {
			Include   []string `yaml:"include"`
			Operators []struct {
				Type string `yaml:"type"`
			} `yaml:"operators"`
		} `yaml:"filelog"`
	} `yaml:"receivers"`
	Processors struct {
		Filter struct {
			Logs struct {
				LogRecord []string `yaml:"log_record"`
			} `yaml:"logs"`
		} `yaml:"filter"`
	} `yaml:"processors"`
	Exporters struct {
		File struct {
			Path string `yaml:"path"`
		} `yaml:"file"`
		Otlp struct {
			Endpoint string `yaml:"endpoint"`
			Tls      struct {
				Insecure bool `yaml:"insecure"`
			} `yaml:"tls"`
			Headers struct {
				ApiKey string `yaml:"api-key"`
			} `yaml:"headers"`
		} `yaml:"otlp"`
	} `yaml:"exporters"`
	Service struct {
		Pipelines struct {
			Metrics struct {
				Receivers []string `yaml:"receivers"`
				Exporters []string `yaml:"exporters"`
			} `yaml:"metrics"`
			Logs struct {
				Receivers  []string `yaml:"receivers"`
				Processors []string `yaml:"processors"`
				Exporters  []string `yaml:"exporters"`
			} `yaml:"logs"`
		} `yaml:"pipelines"`
	} `yaml:"service"`
}

type otelCollectorConfigGenerator struct {
	logger *logger.Logger
}

func newOtelCollectorConfigGenerator(
	logger *logger.Logger,
) *otelCollectorConfigGenerator {
	return &otelCollectorConfigGenerator{
		logger: logger,
	}
}

func (o *otelCollectorConfigGenerator) generate() error {

	cfg := &otelCollectorConfig{
		Receivers: struct {
			Otlp struct {
				Protocols struct {
					Grpc struct {
						Endpoint string "yaml:\"endpoint\""
					} "yaml:\"grpc\""
				} "yaml:\"protocols\""
			} "yaml:\"otlp\""
			Filelog struct {
				Include   []string "yaml:\"include\""
				Operators []struct {
					Type string "yaml:\"type\""
				} "yaml:\"operators\""
			} "yaml:\"filelog\""
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
			Filelog: struct {
				Include   []string "yaml:\"include\""
				Operators []struct {
					Type string "yaml:\"type\""
				} "yaml:\"operators\""
			}{
				Include: []string{
					"./logs/log",
				},
				Operators: []struct {
					Type string "yaml:\"type\""
				}{
					{
						Type: "json_parser",
					},
				},
			},
		},
		Processors: struct {
			Filter struct {
				Logs struct {
					LogRecord []string "yaml:\"log_record\""
				} "yaml:\"logs\""
			} "yaml:\"filter\""
		}{
			Filter: struct {
				Logs struct {
					LogRecord []string "yaml:\"log_record\""
				} "yaml:\"logs\""
			}{
				Logs: struct {
					LogRecord []string "yaml:\"log_record\""
				}{
					LogRecord: []string{
						`IsMatch(attributes["level"], "debug")`,
						// "''IsMatch(attributes["level"], "debug")''",
					},
				},
			},
		},
		Exporters: struct {
			File struct {
				Path string "yaml:\"path\""
			} "yaml:\"file\""
			Otlp struct {
				Endpoint string "yaml:\"endpoint\""
				Tls      struct {
					Insecure bool "yaml:\"insecure\""
				} "yaml:\"tls\""
				Headers struct {
					ApiKey string "yaml:\"api-key\""
				} "yaml:\"headers\""
			} "yaml:\"otlp\""
		}{
			File: struct {
				Path string "yaml:\"path\""
			}{
				Path: "./bin/log",
			},
			Otlp: struct {
				Endpoint string "yaml:\"endpoint\""
				Tls      struct {
					Insecure bool "yaml:\"insecure\""
				} "yaml:\"tls\""
				Headers struct {
					ApiKey string "yaml:\"api-key\""
				} "yaml:\"headers\""
			}{
				Endpoint: "otlp.eu01.nr-data.net:4317",
				Tls: struct {
					Insecure bool "yaml:\"insecure\""
				}{
					Insecure: false,
				},
				Headers: struct {
					ApiKey string "yaml:\"api-key\""
				}{
					ApiKey: os.Getenv("NEWRELIC_LICENSE_KEY"),
				},
			},
		},
		Service: struct {
			Pipelines struct {
				Metrics struct {
					Receivers []string "yaml:\"receivers\""
					Exporters []string "yaml:\"exporters\""
				} "yaml:\"metrics\""
				Logs struct {
					Receivers  []string "yaml:\"receivers\""
					Processors []string "yaml:\"processors\""
					Exporters  []string "yaml:\"exporters\""
				} "yaml:\"logs\""
			} "yaml:\"pipelines\""
		}{
			Pipelines: struct {
				Metrics struct {
					Receivers []string "yaml:\"receivers\""
					Exporters []string "yaml:\"exporters\""
				} "yaml:\"metrics\""
				Logs struct {
					Receivers  []string "yaml:\"receivers\""
					Processors []string "yaml:\"processors\""
					Exporters  []string "yaml:\"exporters\""
				} "yaml:\"logs\""
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
						"otlp",
					},
				},
				Logs: struct {
					Receivers  []string "yaml:\"receivers\""
					Processors []string "yaml:\"processors\""
					Exporters  []string "yaml:\"exporters\""
				}{
					Receivers: []string{
						"filelog",
					},
					Processors: []string{
						"filter",
					},
					Exporters: []string{
						"otlp",
					},
				},
			},
		},
	}

	// Marshal the struct into YAML format
	o.logger.LogWithFields(
		logrus.InfoLevel,
		"Generating OTel config file...",
		map[string]string{
			"component.name": "otelconfiggenerator",
		})
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		o.logger.LogWithFields(
			logrus.ErrorLevel,
			"Generating OTel config failed.",
			map[string]string{
				"component.name": "otelconfiggenerator",
				"error.message":  err.Error(),
			})
		return err
	}

	o.logger.LogWithFields(
		logrus.InfoLevel,
		"Generating OTel config succeeded. Creating file...",
		map[string]string{
			"component.name": "otelconfiggenerator",
		})

	// Create the YAML file
	file, err := os.Create("./bin/otel-config.yaml")
	if err != nil {
		o.logger.LogWithFields(
			logrus.InfoLevel,
			"Creating OTel config file failed.",
			map[string]string{
				"component.name": "otelconfiggenerator",
				"error.message":  err.Error(),
			})
		return err
	}
	defer file.Close()

	o.logger.LogWithFields(
		logrus.InfoLevel,
		"Creating OTel config succeeded. Writing to file...",
		map[string]string{
			"component.name": "otelconfiggenerator",
		})

	// Write the YAML data to the file
	_, err = file.Write(yamlData)
	if err != nil {
		o.logger.LogWithFields(
			logrus.ErrorLevel,
			"Writing OTel config to file failed.",
			map[string]string{
				"component.name": "otelconfiggenerator",
				"error.message":  err.Error(),
			})
		return err
	}
	o.logger.LogWithFields(
		logrus.InfoLevel,
		"Writing OTel config to file succeeded.",
		map[string]string{
			"component.name": "otelconfiggenerator",
		})

	return nil
}
