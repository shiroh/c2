package utils

import (
	"encoding/json"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type animal struct {
	name string
}

func TestMapCreation(t *testing.T) {
	m := NewConcurrentMap()
	assert.NotNil(t, m)
	assert.Equal(t, 0, m.Count(), "new map should be empty.")
}

func TestInsert(t *testing.T) {
	m := NewConcurrentMap()
	elephant := animal{"elephant"}
	monkey := animal{"monkey"}

	m.Set("elephant", elephant)
	m.Set("monkey", monkey)

	assert.Equal(t, 2, m.Count(), "map should contain exactly two elements.")
}

func TestGet(t *testing.T) {
	m := NewConcurrentMap()

	// Get a missing element.
	val, ok := m.Get("Money")
	assert.NotEqual(t, true, ok, "ok should be false when item is missing from map.")
	assert.Nil(t, val, "Missing values should return as null.")

	elephant := animal{"elephant"}
	m.Set("elephant", elephant)

	// Retrieve inserted element.
	tmp, ok := m.Get("elephant")
	elephant, _ = tmp.(animal) // Type assertion.
	assert.Equal(t, true, ok, "ok should be true for item stored within the map.")
	assert.NotNil(t, &elephant, "expecting an element, not null.")
	assert.Equal(t, "elephant", elephant.name)
}

func TestGetElseSet(t *testing.T) {
	m := NewConcurrentMap()

	elephant := animal{"elephant"}
	m.Set("elephant", elephant)
	// Retrieve inserted element.
	tmp, ok := m.GetElseSet("elephant", elephant)
	elephant, _ = tmp.(animal) // Type assertion.
	assert.Equal(t, true, ok, "ok should be true for item stored within the map.")
	assert.NotNil(t, &elephant, "expecting an element, not null.")
	assert.Equal(t, "elephant", elephant.name)

	cow := animal{"cow"}
	tmp, ok = m.GetElseSet("cow", cow)
	cow, _ = tmp.(animal) // Type assertion.
	assert.Equal(t, false, ok, "ok should be false for item newly inserted to the map.")
	assert.NotNil(t, &cow, "expecting an element, not null.")
	assert.Equal(t, "cow", cow.name)

	tmp, ok = m.Get("cow")
	cow, _ = tmp.(animal) // Type assertion.
	assert.Equal(t, true, ok, "ok should be true for item stored within the map.")
	assert.NotNil(t, &cow, "expecting an element, not null.")
	assert.Equal(t, "cow", cow.name)
}

func TestHas(t *testing.T) {
	m := NewConcurrentMap()

	assert.False(t, m.Has("Money"), "element shouldn't exists")
	elephant := animal{"elephant"}
	m.Set("elephant", elephant)
	assert.True(t, m.Has("elephant"), "element exists, expecting Has to return True.")
}

func TestRemove(t *testing.T) {
	m := NewConcurrentMap()

	monkey := animal{"monkey"}
	m.Set("monkey", monkey)
	m.Remove("monkey")
	assert.Equal(t, 0, m.Count(), "Expecting count to be zero once item was removed.")

	temp, ok := m.Get("monkey")
	assert.False(t, ok)
	assert.Nil(t, temp)

	// Remove a none existing element.
	m.Remove("noone")
}

func TestRemoveIfValue(t *testing.T) {
	m := NewConcurrentMap()

	// same value type
	name := "John"
	monkey := "monkey"
	m.Set(name, monkey)

	ok, err := m.RemoveIfValue(name, "tiger")
	assert.False(t, ok)
	assert.Nil(t, err)
	temp, ok := m.Get(name)
	assert.True(t, ok)
	assert.Equal(t, monkey, temp)

	ok, err = m.RemoveIfValue("Cath", monkey)
	assert.False(t, ok)
	assert.Nil(t, err)
	temp, ok = m.Get(name)
	assert.True(t, ok)
	assert.Equal(t, monkey, temp)

	ok, err = m.RemoveIfValue(name, monkey)
	assert.True(t, ok)
	assert.Nil(t, err)
	temp, ok = m.Get(name)
	assert.False(t, ok)
	assert.Nil(t, temp)

	// different value types
	circle := make(chan int)
	m.Set("dog", circle)
	dog := "dog"

	ok, err = m.RemoveIfValue("cat", dog)
	assert.False(t, ok)
	assert.Nil(t, err)
	temp, ok = m.Get("dog")
	assert.True(t, ok)
	assert.Equal(t, circle, temp)

	ok, err = m.RemoveIfValue("dog", dog)
	assert.False(t, ok)
	assert.Nil(t, err)
	temp, ok = m.Get("dog")
	assert.True(t, ok)
	assert.Equal(t, circle, temp)

	ok, err = m.RemoveIfValue("dog", circle)
	assert.True(t, ok)
	assert.Nil(t, err)
	temp, ok = m.Get("dog")
	assert.False(t, ok)
	assert.Nil(t, temp)
}

func TestCount(t *testing.T) {
	m := NewConcurrentMap()
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), animal{strconv.Itoa(i)})
	}
	assert.Equal(t, 100, m.Count())
}

func TestIsEmpty(t *testing.T) {
	m := NewConcurrentMap()
	assert.True(t, m.IsEmpty())

	m.Set("elephant", animal{"elephant"})
	assert.False(t, m.IsEmpty())
}

func TestIterator(t *testing.T) {
	m := NewConcurrentMap()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), animal{strconv.Itoa(i)})
	}

	counter := 0
	// Iterate over elements.
	for _, item := range m.Iter() {
		val := item.Val
		assert.NotNil(t, val)
		counter++
	}

	assert.Equal(t, 100, counter)
}

func TestConcurrent(t *testing.T) {
	m := NewConcurrentMap()
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int

	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			m.Set(strconv.Itoa(i), i)

			// Retrieve item from map.
			val, _ := m.Get(strconv.Itoa(i))

			// Write to channel inserted value.
			ch <- val.(int)
		} // Call go routine with current index.
	}()

	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			m.Set(strconv.Itoa(i), i)

			// Retrieve item from map.
			val, _ := m.Get(strconv.Itoa(i))

			// Write to channel inserted value.
			ch <- val.(int)
		} // Call go routine with current index.
	}()

	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}

	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])

	assert.Equal(t, iterations, m.Count())

	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		assert.Equal(t, i, a[i])
	}
}

func TestJsonMarshal(t *testing.T) {
	shardCount = 2
	defer func() { shardCount = 32 }()
	expected := "{\"a\":1,\"b\":2}"
	m := NewConcurrentMap()
	m.Set("a", 1)
	m.Set("b", 2)
	j, err := json.Marshal(m)
	assert.Nil(t, err)
	assert.Equal(t, string(j), expected)
}
