package hn_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/example/hn/internal/hn"
)

func TestSearch(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hnc := hn.New()
	result, err := hnc.SearchRecent(ctx, &hn.Search{
		Points: "> 500",
	})
	is.NoErr(err)
	is.Equal(len(result.Stories), 20) // 20 newest stories over 500 points
}

func TestShowHN(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hnc := hn.New()
	result, err := hnc.ShowHN(ctx)
	is.NoErr(err)
	is.Equal(len(result.Stories), 20) // 20 show stories
}

func TestAskHN(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hnc := hn.New()
	result, err := hnc.AskHN(ctx)
	is.NoErr(err)
	is.Equal(len(result.Stories), 20) // 20 ask stories
}

func TestNewest(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hnc := hn.New()
	result, err := hnc.Newest(ctx)
	is.NoErr(err)
	is.Equal(len(result.Stories), 20) // 20 newest stories
}

func TestFrontPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hnc := hn.New()
	result, err := hnc.FrontPage(ctx)
	is.NoErr(err)
	is.Equal(len(result.Stories), 20) // 20 front page stories
}

func TestFind(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hnc := hn.New()
	post, err := hnc.Find(ctx, "1")
	is.NoErr(err)
	is.Equal(post.Title, "Y Combinator") // title is not Y Combinator
	buf, err := json.MarshalIndent(post, "", "  ")
	is.NoErr(err)
	fmt.Println(string(buf))
}
