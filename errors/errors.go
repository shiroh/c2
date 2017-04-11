package errors

import "errors"

var (
	// ErrTopicNotFound ...
	ErrTopicNotFound    = errors.New("Topic not found")

	// ErrUnrecognizedType ...
	ErrUnrecognizedType = errors.New("Type in storage is not recognized")
)
