package env

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Getenv can be replaced for testing purposes
var Getenv = os.Getenv

func Collect() *Collector {
	return &Collector{}
}

type Collector struct {
	errors []string
}

func (c *Collector) String(key string) string {
	value := Getenv(key)
	if value == "" {
		c.errors = append(c.errors, fmt.Sprintf("missing %s", key))
	}
	return value
}

func (c *Collector) StringOr(key, defaultValue string) string {
	value := Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func (c *Collector) Error() error {
	if len(c.errors) == 0 {
		return nil
	}
	return errors.New(strings.Join(c.errors, ". "))
}
