package cache

import (
	"sync"
	"testing"

	"github.com/Carou/models"
	"github.com/Carou/storage"
	"github.com/Carou/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEntryRefresh(t *testing.T) {
	entry := &entry{
		sync.RWMutex{},
		nil,
		int64(0),
	}
	newTopics := map[string]*models.Topic{
		mock.Anything: {
			mock.Anything,
			mock.Anything,
			1,
			0,
		},
	}
	entry.refresh(newTopics)

	assert.NotNil(t, entry.topics)

	for k, v := range entry.topics {
		r := newTopics[k]
		assert.Equal(t, v.ID, r.ID, "Topic should be found after refreshing")
		assert.Equal(t, v.Content, r.Content, "Topic should be found after refreshing")
		assert.Equal(t, v.Votes, r.Votes, "Topic should be found after refreshing")
		assert.Equal(t, v.Ts, r.Ts, "Topic should be found after refreshing")
	}
}

func TestEntryIsExpired(t *testing.T) {
	testCases := []struct {
		now      int64
		ttl      int64
		ts       int64
		expected bool
		desc     string
	}{
		{
			now:      1000,
			ttl:      10,
			ts:       989,
			expected: true,
			desc:     "Expired",
		},
		{
			now:      1000,
			ttl:      10,
			ts:       999,
			expected: false,
			desc:     "Not Expired",
		},
	}

	for _, c := range testCases {
		entry := &entry{
			sync.RWMutex{},
			nil,
			c.ts,
		}
		assert.Equal(t, c.expected, entry.isexpired(c.now, c.ttl), c.desc)
	}

}

func TestEntryGetall(t *testing.T) {
	testCases := []struct {
		pageSize int
		topics   *entry
		expected []models.Topic
		desc     string
	}{
		{
			pageSize: 1,
			topics: &entry{
				sync.RWMutex{},
				map[string]*models.Topic{
					mock.Anything: {
						mock.Anything,
						mock.Anything,
						1,
						0,
					},
				},
				int64(0),
			},
			expected: []models.Topic{{mock.Anything, mock.Anything, 1, 0}},
			desc:     "Happypath",
		},
		{
			pageSize: 0,
			topics: &entry{
				sync.RWMutex{},
				map[string]*models.Topic{
					mock.Anything: {
						mock.Anything,
						mock.Anything,
						1,
						0,
					},
				},
				int64(0),
			},
			expected: []models.Topic{},
			desc:     "Tailor the size",
		},
		{
			pageSize: 2,
			topics: &entry{
				sync.RWMutex{},
				map[string]*models.Topic{
					"Foo": {
						"Foo",
						mock.Anything,
						1,
						0,
					},
					"Bar": {
						"Bar",
						mock.Anything,
						2,
						0,
					},
				},
				int64(0),
			},
			expected: []models.Topic{{"Bar", mock.Anything, 2, 0}, {"Foo", mock.Anything, 1, 0}},
			desc:     "Topics should be sorted",
		},
	}

	for _, c := range testCases {
		topics := c.topics.getAll(c.pageSize)
		assert.Equal(t, len(c.expected), len(topics), c.desc)
		for idx, v := range c.expected {
			assert.Equal(t, v, topics[idx])
		}
	}
}

func TestEntryUdpateVote(t *testing.T) {
	testCases := []struct {
		topics   *entry
		id       string
		delta    int64
		expected int64
		desc     string
	}{
		{
			topics: &entry{
				sync.RWMutex{},
				map[string]*models.Topic{
					mock.Anything: {
						mock.Anything,
						mock.Anything,
						1,
						0,
					},
				},
				int64(0),
			},
			id:       mock.Anything,
			delta:    1,
			expected: 2,
			desc:     "Happypath Positive",
		},
		{
			topics: &entry{
				sync.RWMutex{},
				map[string]*models.Topic{
					mock.Anything: {
						mock.Anything,
						mock.Anything,
						1,
						0,
					},
				},
				int64(0),
			},
			id:       mock.Anything,
			delta:    -1,
			expected: 0,
			desc:     "Happypath Negative",
		},
		{
			topics: &entry{
				sync.RWMutex{},
				map[string]*models.Topic{
					mock.Anything: {
						mock.Anything,
						mock.Anything,
						1,
						0,
					},
				},
				int64(0),
			},
			id:       "Foo",
			delta:    1,
			expected: 1,
			desc:     "Id not existed",
		},
	}

	for _, c := range testCases {
		c.topics.updateVote(c.id, c.delta)
		r := c.topics.topics[mock.Anything]
		assert.Equal(t, c.expected, r.Votes, c.desc)
	}
}

func TestV1_refresh(t *testing.T) {
	defer func(originNow func() int64) {
		utils.Now = originNow
	}(utils.Now)

	utils.Now = func() int64 {
		return 1
	}

	testCases := []struct {
		cache     *v1
		newTopics []models.Topic
		expected  []models.Topic
		desc      string
	}{
		{
			cache: &v1{
				topics: &entry{
					sync.RWMutex{},
					make(map[string]*models.Topic),
					int64(0),
				},
				ttl:             0,
				store:           storage.NewV1(),
				defaultPageSize: 1,
			},
			newTopics: []models.Topic{{"Foo", mock.Anything, 1, 0}},
			expected:  []models.Topic{{"Foo", mock.Anything, 1, 0}},
			desc:      "Happypath",
		},
		{
			cache: &v1{
				topics: &entry{
					sync.RWMutex{},
					make(map[string]*models.Topic),
					int64(0),
				},
				ttl:             0,
				store:           storage.NewV1(),
				defaultPageSize: 0,
			},
			newTopics: []models.Topic{{"Foo", mock.Anything, 1, 0}},
			expected:  []models.Topic{},
			desc:      "Tailor the size",
		},
		{
			cache: &v1{
				topics: &entry{
					sync.RWMutex{},
					make(map[string]*models.Topic),
					int64(0),
				},
				ttl:             0,
				store:           storage.NewV1(),
				defaultPageSize: 0,
			},
			newTopics: []models.Topic{{"Bar", mock.Anything, 2, 0}, {"Foo", mock.Anything, 1, 0}},
			expected:  []models.Topic{{"Bar", mock.Anything, 2, 0}, {"Foo", mock.Anything, 1, 0}},
			desc:      "Refresh multiple topics",
		},
	}

	for _, c := range testCases {
		v := c.cache
		for _, topic := range c.newTopics {
			v.store.CreateTopic(topic)
		}
		v.refresh()

		assert.Equal(t, int64(1), v.topics.ts, c.desc)
		assert.Equal(t, len(c.expected), len(v.topics.topics), c.desc)
		for _, topic := range c.expected {
			assert.Equal(t, topic, *v.topics.topics[topic.ID], c.desc)
		}
	}
}

func TestV1_GetTopics(t *testing.T) {
	defer func(originNow func() int64) {
		utils.Now = originNow
	}(utils.Now)

	utils.Now = func() int64 {
		return 0
	}

	testCases := []struct {
		cache    *v1
		expected []models.Topic
		pageSize int
		desc     string
	}{
		{
			cache: &v1{
				topics: &entry{
					sync.RWMutex{},
					map[string]*models.Topic{
						"Foo": {
							"Foo",
							mock.Anything,
							1,
							0,
						},
					},
					int64(0),
				},
				ttl:             10,
				store:           storage.NewV1(),
				defaultPageSize: 20,
			},
			expected: []models.Topic{{"Foo", mock.Anything, 1, 0}},
			pageSize: 1,
			desc:     "Happypath",
		},
	}

	for _, c := range testCases {
		v, _ := c.cache.GetTopics(c.pageSize)
		assert.Equal(t, len(c.expected), len(v), c.desc)
		for idx, topic := range c.expected {
			assert.Equal(t, topic, v[idx], c.desc)
		}
	}
}
