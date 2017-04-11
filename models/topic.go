package models

import (
	"time"

	"github.com/twinj/uuid"
)

// Topic is the struct contains content, votes and creation time
type Topic struct {
	// ID is an auto generated UUID, it will not be displayed on web page. But it distinguishes the the topic from each other
	// even the content are same.
	ID      string `json:"id"`

	// Content is the words created by the user
	Content string `json:"content"`

	// Votes ...
	Votes   int64  `json:"votes"`

	// Ts is the creation time.
	Ts      int64  `json:"ts"`
}

// NewTopic create a new topic object
var NewTopic = func(content string) Topic {
	return Topic{
		ID:      uuid.NewV4().String(),
		Content: content,
		Votes:   0,
		Ts:      time.Now().Unix(),
	}
}

// ByVotes is the type which implement the sort interface.
type ByVotes []Topic

func (a ByVotes) Len() int           { return len(a) }
func (a ByVotes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVotes) Less(i, j int) bool { return a[i].Votes > a[j].Votes }
