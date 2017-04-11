package storage

import (
	"sort"
	"sync/atomic"

	"github.com/Carou/errors"
	"github.com/Carou/models"
	"github.com/Carou/utils"
)

// v1 implement the storage interface. In the demo, it is simply a concurrent hash map.
type v1 struct {
	m utils.ConcurrentMap
}

// NewV1 create a new v1 object
var NewV1 = func() Storage {
	return &v1{
		utils.NewConcurrentMap(),
	}
}

// CreateTopic create a new topic into the storage.
func (impl *v1) CreateTopic(topic models.Topic) error {
	impl.m.Set(topic.ID, &topic)
	return nil
}

// UpvoteTopic increase the votes by 1 if id existed.
// In real world it is ideally put the update event to the message queue
func (impl *v1) UpvoteTopic(topic models.Topic) error {
	t, ok := impl.m.Get(topic.ID)
	if !ok {
		return errors.ErrTopicNotFound
	}
	t2, ok := t.(*models.Topic)
	if !ok {
		return errors.ErrUnrecognizedType
	}

	atomic.AddInt64(&t2.Votes, 1)
	return nil
}

// DownvoteTopic decrease the votes by 1 if id existed
// In real world it is ideally put the update event to the message queue
func (impl *v1) DownvoteTopic(topic models.Topic) error {
	t, ok := impl.m.Get(topic.ID)
	if !ok {
		return errors.ErrTopicNotFound
	}
	t2, ok := t.(*models.Topic)
	if !ok {
		return errors.ErrUnrecognizedType
	}

	atomic.AddInt64(&t2.Votes, -1)
	return nil
}

// GetTopicByVotes return the top {pagesize} topics stored in the storage layer sorted by votes.
func (impl *v1) GetTopicByVotes(pagesize int) ([]models.Topic, error) {
	topics := make([]models.Topic, 0)
	for _, v := range impl.m.Iter() {
		if v.Val == nil {
			utils.Warning.Printf("[GetTopicByVotes][key:%s] Value is empty ", v.Key)
			continue
		}
		t, ok := v.Val.(*models.Topic)
		if !ok {
			utils.Warning.Printf("[GetTopicByVotes][key:%s] Value can't be converted to type Topic", v.Key)
			continue
		}
		topics = append(topics, *t)
	}

	sort.Sort(models.ByVotes(topics))
	if pagesize < len(topics) {
		return topics[:pagesize], nil
	}
	return topics, nil
}
