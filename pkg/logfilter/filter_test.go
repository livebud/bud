package logfilter_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/pkg/log"
	"gitlab.com/mnm/bud/pkg/logfilter"
)

func Test(t *testing.T) {
	t.SkipNow()
	for _, test := range tests {
		t.Run(test.pattern, func(t *testing.T) {
			is := is.New(t)
			matcher, err := logfilter.Parse(test.pattern)
			for _, input := range test.inputs {
				if err != nil {
					is.Equal(err.Error(), input.err)
					return
				}
				is.Equal(matcher.Match(input.entry), input.expect)
			}
		})
	}
}

type input struct {
	entry  log.Entry
	expect bool
	err    string
}

var tests = []struct {
	pattern string
	inputs  []input
}{
	{
		pattern: "debug",
		inputs: []input{
			{
				entry:  log.Entry{Level: log.DebugLevel},
				expect: true,
			},
		},
	},
	{
		pattern: "deBug",
		inputs: []input{
			{
				entry:  log.Entry{Level: log.DebugLevel},
				expect: true,
			},
		},
	},
	{
		pattern: "DEBUG",
		inputs: []input{
			{
				entry:  log.Entry{Level: log.DebugLevel},
				expect: true,
			},
			{ // Log levels include all higher levels
				entry:  log.Entry{Level: log.InfoLevel},
				expect: true,
			},
			{
				entry:  log.Entry{Level: log.NoticeLevel},
				expect: true,
			},
			{
				entry:  log.Entry{Level: log.WarnLevel},
				expect: true,
			},
			{
				entry:  log.Entry{Level: log.ErrorLevel},
				expect: true,
			},
		},
	},
}
