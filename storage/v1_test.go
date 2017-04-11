package storage

import (
	"testing"

	"github.com/Carou/models"
	"github.com/coreos/etcd/gopath/src/github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewV1(t *testing.T) {
	storage := NewV1()
	_, ok := storage.(*v1)
	assert.True(t, ok)
}

func TestV1_UpvoteTopic(t *testing.T) {
	storage := NewV1()
	topic := models.Topic{
		ID:      mock.Anything,
		Content: mock.Anything,
		Votes:   0,
		Ts:      0,
	}

	expectedVotes := int64(1)
	storage.CreateTopic(topic)
	storage.UpvoteTopic(topic)
	updatedTopic, _ := storage.GetTopicByVotes(1)

	assert.Equal(t, expectedVotes, updatedTopic[0].Votes)
}

func TestV1_DownvoteTopic(t *testing.T) {
	storage := NewV1()
	topic := models.Topic{
		ID:      mock.Anything,
		Content: mock.Anything,
		Votes:   0,
		Ts:      0,
	}

	expectedVotes := int64(-1)
	storage.CreateTopic(topic)
	storage.DownvoteTopic(topic)
	updatedTopic, _ := storage.GetTopicByVotes(1)

	assert.Equal(t, expectedVotes, updatedTopic[0].Votes)
}

func TestV1_CreateTopic(t *testing.T) {
	storage := NewV1()
	topic := models.Topic{
		ID:      mock.Anything,
		Content: mock.Anything,
		Votes:   1,
		Ts:      1,
	}

	storage.CreateTopic(topic)
	newTopic, _ := storage.GetTopicByVotes(1)

	assert.Equal(t, topic.ID, newTopic[0].ID)
	assert.Equal(t, topic.Content, newTopic[0].Content)
	assert.Equal(t, topic.Votes, newTopic[0].Votes)
	assert.Equal(t, topic.Ts, newTopic[0].Ts)
}

func TestV1_GetTopicByVotes2(t *testing.T) {
	testCases := []struct {
		topics   []models.Topic
		expected []models.Topic
		pageSize int
		desc     string
	}{
		{
			topics:   []models.Topic{{"Foo", mock.Anything, 1, 1}},
			expected: []models.Topic{{"Foo", mock.Anything, 1, 1}},
			pageSize: 1,
			desc:     "Single Topic",
		},
		{
			topics:   []models.Topic{{"Foo", mock.Anything, 1, 1}, {"Bar", mock.Anything, 2, 1}},
			expected: []models.Topic{{"Bar", mock.Anything, 2, 1}, {"Foo", mock.Anything, 1, 1}},
			pageSize: 2,
			desc:     "Two Topics, pagesize is 2",
		},
		{
			topics:   []models.Topic{{"Foo", mock.Anything, 1, 1}, {"Bar", mock.Anything, 2, 1}},
			expected: []models.Topic{{"Bar", mock.Anything, 2, 1}},
			pageSize: 1,
			desc:     "Two Topics, pagesize is 1",
		},
		{
			topics:   []models.Topic{{"Foo", mock.Anything, 1, 1}, {"Foo", mock.Anything, 2, 1}},
			expected: []models.Topic{{"Foo", mock.Anything, 2, 1}},
			pageSize: 1,
			desc:     "Two Topics with same key",
		},
	}

	for _, c := range testCases {
		storage := NewV1()
		for _, topic := range c.topics {
			storage.CreateTopic(topic)
		}

		newTopics, _ := storage.GetTopicByVotes(c.pageSize)
		assert.Equal(t, len(c.expected), len(newTopics), c.desc)
		for idx, topic := range c.expected {
			assert.Equal(t, topic, newTopics[idx], c.desc)
		}
	}

}
