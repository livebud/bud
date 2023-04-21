package buildinfo

import (
	"fmt"
	"strings"
)

var Version string
var BuildTime string
var Builder string
var BuildOS string
var Custom string

type FlagSet struct {
	flgs map[string]string
}

// Load loads all build flags targeted with -ldflags and
// returns a FlagSet for retrieving flags later
func Load() FlagSet {
	f := FlagSet{flgs: make(map[string]string)}

	f.flgs["v"] = Version
	f.flgs["bt"] = BuildTime
	f.flgs["b"] = Builder
	f.flgs["bOS"] = BuildOS
	f.flgs["c"] = Custom

	return f
}

func buildAll(first bool, all, k, v string) string {
	switch k {
	case "v":
		k = "Version"
	case "b":
		k = "Builder"
	case "bt":
		k = "Build Time"
	case "bOS":
		k = "Build OS"
	case "c":
		k = "Custom"
	}

	if first {
		all = fmt.Sprintf("%s %s=%s", all, strings.TrimSpace(k), strings.TrimSpace(v))
	} else {
		all = fmt.Sprintf("%s, %s=%s", all, strings.TrimSpace(k), strings.TrimSpace(v))
	}

	return all
}

// All returns all the set buildinfo flags in f as a string, following
// the pattern "All set build info flags by key/value pairs: k=v" where
// k is the specific flag set and v is the value assigned to k. This method
// does not comma-separate Custom buildinfo flags. To return a string with
// comma-separated Custom buildinfo flags, use AllSeparated()
func (f FlagSet) All() string {
	all := "All set build info flags by key/value pairs:"

	first := true
	for k, v := range f.flgs {
		if k == "" {
			continue
		}

		all = buildAll(first, all, k, v)

		if first {
			first = false
		}
	}

	return all
}

// AllSeparated returns all the set buildinfo flags in f as a string, following
// the pattern "All set build info flags by key/value pairs: k=v" where
// k is the specific flag set and v is the value assigned to k. This method
// comma-separates Custom buildinfo flags. To return a string without Custom
// buildinfo flags being comma-separated, use All()
func (f FlagSet) AllSeparated() string {
	all := "All set build info flags by key/value pairs:"

	first := true
	for k, v := range f.flgs {
		switch k {
		case "":
			continue
		case "c":
			split := strings.Split(v, ",")
			for _, val := range split {
				kv := strings.Split(val, "=")
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				all = buildAll(first, all, key, value)
			}
		default:
			all = buildAll(first, all, k, v)
		}
		if first {
			first = false
		}
	}

	return all
}

// Version returns the value set for the Version buildinfo flag
func (f FlagSet) Version() string {
	return f.flgs["v"]
}

// BuildTime returns the value set for the BuildTime buildinfo flag
func (f FlagSet) BuildTime() string {
	return f.flgs["bt"]
}

// Builder returns the value set for the Builder buildinfo flag
func (f FlagSet) Builder() string {
	return f.flgs["b"]
}

// BuildOS returns the value set for the BuildOS buildinfo flag
func (f FlagSet) BuildOS() string {
	return f.flgs["bOS"]
}

// Custom returns a map of Custom buildinfo flags with the values
// split on commas, and then each split into a key/value pair on "="
// and stored into the returned map. If no Custom buildinfo flags
// were set, then this method returns nil
func (f FlagSet) Custom() map[string]string {
	if _, ok := f.flgs["c"]; !ok {
		return nil
	}

	cv := make(map[string]string)

	c := strings.Split(f.flgs["c"], ",")

	for _, custom := range c {
		kv := strings.Split(custom, "=")
		cv[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	return cv
}
