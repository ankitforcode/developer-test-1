package externalservice

import (
	"errors"
	"sync"
)

// Post is the data structure representing the data sent and received from the
// external service
type Post struct {
	ID int `json:"id"` // the primary key

	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// Client represents the client interface to the external service
type Client interface {
	GET(id int) (*Post, error)
	POST(id int, post *Post) (*Post, error)
}

// ClientImpl implements Client interfacte
type ClientImpl struct {
	Posts map[int]*Post
	mu    sync.RWMutex
}

func (c *ClientImpl) GET(id int) (*Post, error) {
	if c.Posts[id] == nil {
		return nil, errors.New("Post not found")
	}
	return c.Posts[id], nil
}

func (c *ClientImpl) POST(id int, post *Post) (*Post, error) {
	if c.Posts[id] != nil {
		return nil, errors.New("Post id already exists")
	}
	c.mu.Lock()
	c.Posts[id] = post
	c.mu.Unlock()
	return post, nil
}
