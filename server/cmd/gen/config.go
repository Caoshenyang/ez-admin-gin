package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Module    string  `yaml:"module"`
	Group     string  `yaml:"group"`
	Label     string  `yaml:"label"`
	LabelEn   string  `yaml:"label_en"`
	Table     string  `yaml:"table"`
	Fields    []Field `yaml:"fields"`
	HasStatus bool    `yaml:"has_status"`
}

type Field struct {
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	DbType     string `yaml:"db_type"`
	JsonType   string `yaml:"json_type"`
	Label      string `yaml:"label"`
	Form       bool   `yaml:"form"`
	FormType   string `yaml:"form_type"`
	Required   bool   `yaml:"required"`
	Filter     bool   `yaml:"filter"`
	FilterType string `yaml:"filter_type"`
	IsStatus   bool   `yaml:"is_status"`
}

func ParseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Module == "" {
		return fmt.Errorf("module is required")
	}
	if c.Group == "" {
		return fmt.Errorf("group is required")
	}
	if c.Label == "" {
		return fmt.Errorf("label is required")
	}
	if c.LabelEn == "" {
		return fmt.Errorf("label_en is required")
	}
	if c.Table == "" {
		return fmt.Errorf("table is required")
	}
	if len(c.Fields) == 0 {
		return fmt.Errorf("at least one field is required")
	}
	for i, f := range c.Fields {
		if f.Name == "" {
			return fmt.Errorf("fields[%d].name is required", i)
		}
		if f.Type == "" {
			return fmt.Errorf("fields[%d].type is required", i)
		}
	}
	return nil
}

func (c *Config) FormFields() []Field {
	var fields []Field
	for _, f := range c.Fields {
		if f.Form {
			fields = append(fields, f)
		}
	}
	return fields
}

func (c *Config) FilterFields() []Field {
	var fields []Field
	for _, f := range c.Fields {
		if f.Filter {
			fields = append(fields, f)
		}
	}
	return fields
}

func (c *Config) StatusField() *Field {
	for i := range c.Fields {
		if c.Fields[i].IsStatus {
			return &c.Fields[i]
		}
	}
	return nil
}

func (f *Field) GoTypeResolved() string {
	if strings.HasPrefix(f.Type, "model.") {
		return f.Type
	}
	return f.Type
}

func (f *Field) StatusTypeName() string {
	if strings.HasPrefix(f.Type, "model.") {
		return strings.TrimPrefix(f.Type, "model.")
	}
	return "int"
}

func title(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
