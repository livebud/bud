package viewkey_test

import (
	"testing"

	"github.com/livebud/bud/pkg/controller/internal/viewkey"
	"github.com/matthewmueller/diff"
)

func equal(t *testing.T, route, expect string) {
	t.Helper()
	t.Run(route, func(t *testing.T) {
		key, err := viewkey.Infer(route)
		if err != nil {
			if err.Error() == expect {
				return
			}
			t.Fatal(err)
		}
		diff.TestString(t, expect, key)
	})
}

func TestInfer(t *testing.T) {
	equal(t, `/`, "index")
	equal(t, `/{id}`, "show")
	equal(t, `/posts/{id}`, "posts/show")
	equal(t, `/{year}/{month}`, "show")
	equal(t, `/posts/{year}/{month}`, "posts/show")
	equal(t, `/{id}/edit`, "edit")
	equal(t, `/posts`, "posts/index")
	equal(t, `/posts`, "posts/index")
	equal(t, `/about`, "about/index") // This needs to be handled in view finder

	// Extensions
	equal(t, `/.{format}`, "index")
	equal(t, `/.json`, "index")
	equal(t, `/{id}.{format}`, "show")
	equal(t, `/{id}.json`, "show")
	equal(t, `/posts/{post_id}/comments.{format}`, "posts/comments/index")
	equal(t, `/posts/{post_id}/comments.json`, "posts/comments/index")
	equal(t, `/posts/{post_id}/comments/{id}.{format}`, "posts/comments/show")
	equal(t, `/posts/{post_id}/comments/{id}.json`, "posts/comments/show")
	equal(t, `/posts/{post_id}/comments/{id}/edit.{format}`, "posts/comments/edit")
	equal(t, `/posts/{post_id}/comments/{id}/edit.json`, "posts/comments/edit")
	equal(t, `/posts/{post_id}/comments/new.{format}`, "posts/comments/new")
	equal(t, `/posts/{post_id}/comments/new.json`, "posts/comments/new")

	// Layouts, Frames, Errors
	equal(t, `/layout`, "layout")
	equal(t, `/posts/layout`, "posts/layout")
	equal(t, `/posts/{post_id}/comments/{id}/layout`, "posts/comments/layout")
	equal(t, `/frame`, "frame")
	equal(t, `/posts/frame`, "posts/frame")
	equal(t, `/posts/{post_id}/comments/{id}/frame`, "posts/comments/frame")
	equal(t, `/error`, "error")
	equal(t, `/posts/error`, "posts/error")
	equal(t, `/posts/{post_id}/comments/{id}/error`, "posts/comments/error")
}
