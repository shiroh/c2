package storage

import (
	"testing"

	"github.com/Carou/models"
	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(m *MockStorage) error
	}{
		{
			desc: "CreateTopic",
			testFunc: func(m *MockStorage) error {
				m.OnCreateTopic().Return(nil).Once()
				m.CreateTopic(models.Topic{})
				return nil
			},
		},
		{
			desc: "DownvoteTopic",
			testFunc: func(m *MockStorage) error {
				m.OnDownvoteTopic().Return(nil).Once()
				m.DownvoteTopic(models.Topic{})
				return nil
			},
		},
		{
			desc: "UpvoteTopic",
			testFunc: func(m *MockStorage) error {
				m.OnUpvoteTopic().Return(nil).Once()
				m.UpvoteTopic(models.Topic{})
				return nil
			},
		},
		{
			desc: "GetTopicByVotes",
			testFunc: func(m *MockStorage) error {
				m.OnGetTopicByVotes().Return([]models.Topic{}, nil).Once()
				_, _ = m.GetTopicByVotes(1)
				return nil
			},
		},
	}
	for _, c := range testCases {
		m := &MockStorage{}
		err := c.testFunc(m)
		assert.Nil(t, err, c.desc)
		assert.True(t, m.AssertExpectations(t), c.desc)
	}
}
