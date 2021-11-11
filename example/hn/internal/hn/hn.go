// Package hn is a simple HTTP client for Hacker News.
//
// Algolia graciously provided an API for working with Hacker News over at:
// https://hn.algolia.com/api
package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const baseURL = `http://hn.algolia.com/api/v1`

// New HackerNews Client with defaults
func New() *Client {
	return &Client{http.DefaultClient}
}

// Client for HackerNews. The HTTP Client can be overriden with your own.
type Client struct {
	*http.Client
}

// FrontPage is a convenience function for getting the results on
// https://hackernews.com
func (c *Client) FrontPage(ctx context.Context) (*News, error) {
	return c.Search(ctx, &Search{
		Tags: "front_page",
	})
}

// Newest is a convenience function for getting the results on
// https://news.ycombinator.com/newest
func (c *Client) Newest(ctx context.Context) (*News, error) {
	return c.SearchRecent(ctx, &Search{
		Tags: "story",
	})
}

// AskHN is a convenience function for getting the results on
// https://news.ycombinator.com/ask
func (c *Client) AskHN(ctx context.Context) (*News, error) {
	return c.SearchRecent(ctx, &Search{
		Tags: "ask_hn",
	})
}

// ShowHN is a convenience function for getting the results on
// https://news.ycombinator.com/show
func (c *Client) ShowHN(ctx context.Context) (*News, error) {
	return c.SearchRecent(ctx, &Search{
		Tags: "show_hn",
	})
}

// Story is an individual entry on HackerNews.
type Story struct {
	ID          int        `json:"id,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	CreatedAtI  int        `json:"created_at_i,omitempty"`
	Type        string     `json:"type,omitempty"`
	Author      string     `json:"author,omitempty"`
	Title       string     `json:"title,omitempty"`
	URL         string     `json:"url,omitempty"`
	Text        *string    `json:"text,omitempty"`
	NumComments int        `json:"num_comments,omitempty"`
	Points      int        `json:"points,omitempty"`
	ParentID    *int       `json:"parent_id,omitempty"`
	StoryID     *int       `json:"story_id,omitempty"`
	Children    []Children `json:"children"`
}

// Children contain comments.
type Children struct {
	ID         int        `json:"id,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
	CreatedAtI int        `json:"created_at_i,omitempty"`
	Type       string     `json:"type,omitempty"`
	Author     string     `json:"author,omitempty"`
	Title      *string    `json:"title,omitempty"`
	URL        *string    `json:"url,omitempty"`
	Text       string     `json:"text,omitempty"`
	Points     *int       `json:"points,omitempty"`
	ParentID   int        `json:"parent_id,omitempty"`
	StoryID    int        `json:"story_id,omitempty"`
	Children   []Children `json:"children"`
}

// Find a Story by ID.
func (c *Client) Find(ctx context.Context, id string) (*Story, error) {
	url := baseURL + "/items/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	story := new(Story)
	if err := json.Unmarshal(body, story); err != nil {
		return nil, err
	}
	recursivelySort(story.Children)
	return story, nil
}

func recursivelySort(children []Children) {
	sort.Slice(children, func(a, b int) bool {
		return children[a].CreatedAtI < children[b].CreatedAtI
	})
	for _, child := range children {
		recursivelySort(child.Children)
	}
}

// Search query and filters
type Search struct {
	// Full-text query to search for (e.g. Duo)
	Query string

	// Tags filters the search on a specific tag.
	//
	// The available tags are:
	//   - story
	//   - comment
	//   - poll
	//   - pollopt
	//   - show_hn
	//   - ask_hn,
	//   - front_page
	//   - author_:USERNAME
	//   - story_:ID
	//
	// Tags are ANDed by default, can be ORed if between parenthesis. For example,
	// `author_pg,(story,poll)` filters on `author=pg AND (type=story OR type=poll)`
	Tags string

	// Filter by points. Points is a conditional query, so you can request stories
	// that have more than 500 points with "points > 500".
	Points string

	// Filter by date. CreatedAt is a conditional query, so you can request
	// stories between a time period wtih "created_at_i>X,created_at_i<Y", where
	// X and Y are timestamps in seconds.
	CreatedAt string

	// Filter by the number of comments. Comments is a conditional query, so you
	// can request stories that have more than 10 comments with "comments > 10".
	NumComments string

	// The page number
	Page int

	// ResultsPerPage is the number of results. Defaults to 34.
	ResultsPerPage int
}

// Turns the search input into a query string.
func (s *Search) querystring() string {
	query := url.Values{}
	if s.Query != "" {
		query.Set("query", s.Query)
	}
	if s.Tags != "" {
		query.Set("tags", s.Tags)
	}
	var nfs []string
	if s.Points != "" {
		nfs = append(nfs, injectKey(s.Points, "points"))
	}
	if s.CreatedAt != "" {
		nfs = append(nfs, injectKey(s.CreatedAt, "created_at_i"))
	}
	if s.NumComments != "" {
		nfs = append(nfs, injectKey(s.NumComments, "num_comments"))
	}
	if len(nfs) > 0 {
		query.Set("numericFilters", strings.Join(nfs, ","))
	}
	// Set the number of results per page
	if s.ResultsPerPage == 0 {
		// For some reason the number of results returned by default is 34 results
		s.ResultsPerPage = 34
	}
	query.Set("hitsPerPage", strconv.Itoa(s.ResultsPerPage))
	return query.Encode()
}

// Sugar on top to allow bot "points > 500" and "> 500", to reduce repetition
// with the key (e.g. Points: "points > 500")
func injectKey(query, key string) string {
	parts := strings.Split(query, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
		if !strings.HasPrefix(parts[i], key) {
			parts[i] = key + parts[i]
		}
	}
	return strings.Join(parts, ",")
}

// result of a search
type result struct {
	Stories              []*resultStory `json:"hits,omitempty"`
	NumResults           int            `json:"nbHits,omitempty"`
	Page                 int            `json:"page,omitempty"`
	NumPages             int            `json:"nbPages,omitempty"`
	ResultsPerPage       int            `json:"hitsPerPage,omitempty"`
	ExhaustiveNumResults bool           `json:"exhaustiveNbHits,omitempty"`
	Query                string         `json:"query,omitempty"`
	Params               string         `json:"params,omitempty"`
	ProcessingTimeMS     int            `json:"processingTimeMS,omitempty"`
}

type News struct {
	Stories        []*Story `json:"stories,omitempty"`
	NumResults     int      `json:"num_results,omitempty"`
	Page           int      `json:"page,omitempty"`
	NumPages       int      `json:"num_pages,omitempty"`
	ResultsPerPage int      `json:"stories_per_page,omitempty"`
}

// resultStory is an individual Story from a search result
type resultStory struct {
	ID             string    `json:"objectID,omitempty"`
	Title          string    `json:"title,omitempty"`
	URL            string    `json:"url,omitempty"`
	Author         string    `json:"author,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	Points         int       `json:"points,omitempty"`
	StoryText      *string   `json:"story_text,omitempty"`
	CommentText    *string   `json:"comment_text,omitempty"`
	NumComments    int       `json:"num_comments,omitempty"`
	StoryID        *int      `json:"story_id,omitempty"`
	StoryTitle     *string   `json:"story_title,omitempty"`
	StoryURL       *string   `json:"story_url,omitempty"`
	ParentID       *int      `json:"parent_id,omitempty"`
	CreatedAtI     int       `json:"created_at_i,omitempty"`
	RelevancyScore *int      `json:"relevancy_score,omitempty"`
	Tags           []string  `json:"_tags,omitempty"`
	Highlights     struct {
		Title     Highlight `json:"title,omitempty"`
		URL       Highlight `json:"url,omitempty"`
		Author    Highlight `json:"author,omitempty"`
		StoryText Highlight `json:"story_text,omitempty"`
	} `json:"_highlightResult,omitempty"`
	Children []Children `json:"children"`
}

// Highlight indicates the words that matched the search query
type Highlight struct {
	Value        string   `json:"value,omitempty"`
	MatchLevel   string   `json:"matchLevel,omitempty"`
	MatchedWords []string `json:"matchedWords,omitempty"`
}

// Search for Stories. Sorted by relevance, then points, then number of comments.
func (c *Client) Search(ctx context.Context, search *Search) (*News, error) {
	url := baseURL + "/search?" + search.querystring()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	result := new(result)
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return toNews(result)
}

// Search for Stories. Sorted by date, more recent first.
func (c *Client) SearchRecent(ctx context.Context, search *Search) (*News, error) {
	url := baseURL + "/search_by_date?" + search.querystring()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	result := new(result)
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return toNews(result)
}

func toNews(result *result) (*News, error) {
	news := &News{
		NumResults:     result.NumResults,
		Page:           result.Page,
		NumPages:       result.NumPages,
		ResultsPerPage: result.ResultsPerPage,
	}
	for _, story := range result.Stories {
		id, err := strconv.Atoi(story.ID)
		if err != nil {
			return nil, err
		}
		news.Stories = append(news.Stories, &Story{
			Author:      story.Author,
			Children:    []Children{},
			CreatedAt:   story.CreatedAt,
			CreatedAtI:  story.CreatedAtI,
			ID:          id,
			NumComments: story.NumComments,
			ParentID:    story.ParentID,
			Points:      story.Points,
			StoryID:     story.StoryID,
			Title:       story.Title,
			Text:        nil,
			URL:         story.URL,
		})
	}
	return news, nil
}
