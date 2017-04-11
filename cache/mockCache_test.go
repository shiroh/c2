package cache

import (
	"testing"

	"github.com/Carou/models"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(m *MockCache) error
	}{
		{
			desc: "DownvoteTopic",
			testFunc: func(m *MockCache) error {
				m.OnDownvoteTopic().Return(nil).Once()
				m.DownvoteTopic(models.Topic{})
				return nil
			},
		},
		{
			desc: "UpvoteTopic",
			testFunc: func(m *MockCache) error {
				m.OnUpvoteTopic().Return(nil).Once()
				m.UpvoteTopic(models.Topic{})
				return nil
			},
		},
		{
			desc: "UpvoteTopic",
			testFunc: func(m *MockCache) error {
				m.OnGetTopics().Return([]models.Topic{}, nil).Once()
				m.GetTopics(1)
				return nil
			},
		},
	}

	for _, c := range testCases {
		m := &MockCache{}
		err := c.testFunc(m)
		assert.Nil(t, err, c.desc)
		assert.True(t, m.AssertExpectations(t), c.desc)
	}

}
