package users

import (
	context "context"
)

type Controller struct {
	// Dependencies...
}

// User struct
type User struct {
	// Fields...
}

// Index of users
// GET /users
func (c *Controller) Index(ctx context.Context) (users []*User, err error) {
	return []*User{}, nil
}

// Show user
// GET /users/:id
func (c *Controller) Show(ctx context.Context, id int) (user *User, err error) {
	return &User{}, nil
}

// Delete user
// DELETE /users/:id
func (c *Controller) Delete(ctx context.Context, id int) (user User, err error) {
	return User{}, nil
}

// Update user
// PATCH /users/:id
func (c *Controller) Update(ctx context.Context, id int) (user User, err error) {
	return User{}, nil
}

// Create user
// POST /users
func (c *Controller) Create(ctx context.Context) (user User, err error) {
	return User{}, nil
}

// New user
// GET /users/new
func (c *Controller) New(ctx context.Context) (user User, err error) {
	return User{}, nil
}

// Edit user
// GET /users/:id/edit
func (c *Controller) Edit(ctx context.Context, id int) (user User, err error) {
	return User{}, nil
}
