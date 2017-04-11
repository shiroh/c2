package storage

import "github.com/Carou/models"

// Storage is interface which define methods interact with storage layer.
type Storage interface {

	// CreateTopic create a new topic in the storage
	CreateTopic(topic models.Topic) error

	// UpvoteTopic increase the votes by 1 if topic id existed
	UpvoteTopic(topic models.Topic) error

	// DownvoteTopic decrease the votes by 1 if topic id existed
	DownvoteTopic(topic models.Topic) error

	// GetTopicByVotes retrieve the top {pagesize} topics from storage sort by votes.
	GetTopicByVotes(pagesize int) ([]models.Topic, error)
}
