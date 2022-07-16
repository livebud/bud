package controller

import (
	context "context"
)

// Controller for posts
type Controller struct {
}

// Post struct
type Post struct {
	ID int `json:"id"`
}

// Index of posts
// GET
func (c *Controller) Index(ctx context.Context) (posts []*Post, err error) {
	return []*Post{}, nil
}

// New returns a view for creating a new post
// GET new
func (c *Controller) New(ctx context.Context) {
}

// Create post
// POST
func (c *Controller) Create(ctx context.Context) (post *Post, err error) {
	return &Post{
		ID: 0,
	}, nil
}

// Show post
// GET :id
func (c *Controller) Show(ctx context.Context, id int) (post *Post, err error) {
	return &Post{
		ID: id,
	}, nil
}

// Edit returns a view for editing a post
// GET :id/edit
func (c *Controller) Edit(ctx context.Context, id int) (post *Post, err error) {
	return &Post{
		ID: id,
	}, nil
}

// Update post
// PATCH :id
func (c *Controller) Update(ctx context.Context, id int) (post *Post, err error) {
	return &Post{
		ID: id,
	}, nil
}

// Delete post
// DELETE :id
func (c *Controller) Delete(ctx context.Context, id int) error {
	return nil
}
