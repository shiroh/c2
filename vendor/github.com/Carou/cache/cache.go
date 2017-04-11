package cache

import "github.com/Carou/models"

// Cache is the interface which define the in memory cache operations.
type Cache interface {
	// UpvoteTopic increase the votes of the topic by 1.
	UpvoteTopic(topic models.Topic) error

	// DownvoteTopic decrease the votes of the topic by 1
	DownvoteTopic(topic models.Topic) error

	// GetTopics return the topics tailored by pagesize.
	GetTopics(pageSize int) ([]models.Topic, error)
}
