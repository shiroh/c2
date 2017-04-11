package cache

import (
	"sort"
	"sync"
	"time"

	"github.com/Carou/models"
	"github.com/Carou/storage"
	"github.com/Carou/utils"
)

// entry is the wrapper contains the topic array and timestamp it loaded
type entry struct {

	// we use a Read/Write mutex here because topics is refreshed periodically
	sync.RWMutex

	topics map[string]*models.Topic
	ts     int64
}

// refresh refresh the array inside and loaded time
func (e *entry) refresh(newTopics map[string]*models.Topic) {
	e.Lock()
	defer e.Unlock()

	e.topics = newTopics
	e.ts = utils.Now()
}

// getAll return all the topics contained, tailored by the page size, sorted by the votes.
func (e *entry) getAll(pageSize int) []models.Topic {
	e.RLock()
	defer e.RUnlock()

	r := make([]models.Topic, 0)
	count := 0
	for _, v := range e.topics {
		if count >= pageSize {
			break
		}

		r = append(r, *v)
		count += 1
	}

	sort.Sort(models.ByVotes(r))

	return r
}

// updateVote change votes of topic by id.
func (e *entry) updateVote(id string, delta int64) {
	e.Lock()
	defer e.Unlock()

	t, ok := e.topics[id]
	if ok {
		t.Votes += delta
	}
}

// isexpired check whether the topic contained inside is expired or not.
func (e *entry) isexpired(now int64, ttl int64) bool {
	e.RLock()
	defer e.RUnlock()

	if now-e.ts > ttl {
		return true
	}
	return false
}

// v1 is the implementation of cache. It returns the cached topics or loaded from the storage is cache missed.
type v1 struct {
	topics          *entry
	ttl             int64
	store           storage.Storage
	defaultPageSize int
}

// NewV1 create a new V1 Cache.
var NewV1 = func(store storage.Storage) Cache {
	return &v1{
		topics: &entry{
			sync.RWMutex{},
			make(map[string]*models.Topic),
			int64(0),
		},
		ttl:             int64(1 * time.Second),
		store:           store,
		defaultPageSize: 20,
	}
}

// UpvoteTopic increase the votes of topic in cache by 1 and the topic in backend storage by 1
func (impl *v1) UpvoteTopic(topic models.Topic) error {
	impl.topics.updateVote(topic.ID, 1)
	return impl.store.UpvoteTopic(topic)
}

// DownvoteTopic decrease the votes of topic in cache by 1 and the topic in backend storage by 1
func (impl *v1) DownvoteTopic(topic models.Topic) error {
	impl.topics.updateVote(topic.ID, -1)
	return impl.store.DownvoteTopic(topic)
}

// GetTopics load the topics from backend storage if they are expired and return.
func (impl *v1) GetTopics(pageSize int) ([]models.Topic, error) {
	err := impl.refresh()
	if err != nil {
		return nil, err
	}
	return impl.topics.getAll(pageSize), nil
}

// refresh load the topics from storage if it is expired.
func (impl *v1) refresh() error {
	if impl.topics.isexpired(utils.Now(), impl.ttl) {
		t, err := impl.store.GetTopicByVotes(impl.defaultPageSize)
		if err != nil {
			return err
		}

		m := make(map[string]*models.Topic)
		for idx, v := range t {
			// we can't use m[v.ID] = &v here.
			// because it assign the address of "v" to m[v.ID], but in the for loop the address of v is always same
			// but fields inside is updated when iterating.
			m[v.ID] = &t[idx]
		}
		impl.topics.refresh(m)
	}
	return nil
}
