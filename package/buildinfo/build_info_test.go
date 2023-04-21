package buildinfo

import (
	"github.com/livebud/bud/internal/is"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	Version = "1"
	BuildTime = time.RFC3339
	Builder = "test builder"
	BuildOS = "Darwin"
	Custom = "key1=value1,key2=value2,key3=value3"
	TestCustom := make(map[string]string)
	TestCustom["key1"] = "value1"
	TestCustom["key2"] = "value2"
	TestCustom["key3"] = "value3"

	flags := Load()

	is := is.New(t)
	is.Equal(5, len(flags.flgs))
	is.Equal(Version, flags.Version())
	is.Equal(BuildTime, flags.BuildTime())
	is.Equal(Builder, flags.Builder())
	for k, v := range flags.Custom() {
		is.Equal(TestCustom[k], v)
	}
}

func TestLoad_SpacesInCustom(t *testing.T) {
	Version = "1"
	BuildTime = time.RFC3339
	Builder = "test builder"
	BuildOS = "Darwin"
	Custom = "key1=value1, key2= value2,key3 =value3 "
	TestCustom := make(map[string]string)
	TestCustom["key1"] = "value1"
	TestCustom["key2"] = "value2"
	TestCustom["key3"] = "value3"

	flags := Load()

	is := is.New(t)
	is.Equal(Version, flags.Version())
	is.Equal(BuildTime, flags.BuildTime())
	is.Equal(Builder, flags.Builder())
	for k, v := range flags.Custom() {
		is.Equal(TestCustom[k], v)
	}
}

func TestAll(t *testing.T) {
	Version = "1"
	BuildTime = time.RFC3339
	Builder = "test builder"
	BuildOS = "Darwin"
	Custom = "key1=value1,key2=value2,key3=value3 "
	TestCustom := make(map[string]string)
	TestCustom["key1"] = "value1"
	TestCustom["key2"] = "value2"
	TestCustom["key3"] = "value3"

	flags := Load()

	expected := "All set build info flags by key/value pairs: Version=1, Build Time=2006-01-02T15:04:05Z07:00, Builder=test builder, Build OS=Darwin, Custom=key1=value1,key2=value2,key3=value3"
	all := flags.All()

	is := is.New(t)

	is.Equal(len(expected), len(all))
}

func TestAllSeparated(t *testing.T) {
	Version = "1"
	BuildTime = time.RFC3339
	Builder = "test builder"
	BuildOS = "Darwin"
	Custom = "key1=value1,key2=value2,key3=value3 "
	TestCustom := make(map[string]string)
	TestCustom["key1"] = "value1"
	TestCustom["key2"] = "value2"
	TestCustom["key3"] = "value3"

	flags := Load()

	expected := "All set build info flags by key/value pairs: Version=1, Build Time=2006-01-02T15:04:05Z07:00, Builder=test builder, Build OS=Darwin, key1=value1, key2=value2, key3=value3"
	all := flags.AllSeparated()

	is := is.New(t)

	is.Equal(len(expected), len(all))
}
