package dataloader

import (
	"fmt"
	"sync"
	"time"
)

// Inspired by:
// - https://github.com/vektah/dataloaden
// - https://github.com/vikstrous/dataloadgen
// - https://github.com/graph-gophers/dataloader

type Config[K comparable, V any] struct {
	// Fetch is a method that provides the data for the loader
	Fetch func(keys []K) ([]V, []error)

	// MaxWait is maximum time to wait before sending a batch
	MaxWait time.Duration

	// MaxSize is maximum number of keys to send in one batch, 0 = not limit
	MaxSize int
}

// New creates a new Dataloader given a fetch function, maxWait, and maxSize
func New[K comparable, V any](cfg Config[K, V]) *Dataloader[K, V] {
	if cfg.MaxWait == 0 {
		cfg.MaxWait = 2 * time.Millisecond
	}
	return &Dataloader[K, V]{
		fetch:   cfg.Fetch,
		maxWait: cfg.MaxWait,
		maxSize: cfg.MaxSize,
		cache:   map[K]func() (V, error){},
	}
}

// Dataloader batches and caches requests
type Dataloader[K comparable, V any] struct {
	// this method provides the data for the loader
	fetch func(keys []K) ([]V, []error)

	// maximum wait time before sending batch
	maxWait time.Duration

	// maximum number of keys in batch
	maxSize int

	// cache
	cache map[K]func() (V, error)

	// the current batch. keys will continue to be collected until timeout is hit,
	// then everything will be sent to the fetch method and out to the listeners
	batch *batch[K, V]

	// mutex to prevent races
	mu sync.Mutex
}

type batch[K comparable, V any] struct {
	keys    []K
	data    []V
	error   []error
	closing bool
	done    chan struct{}
}

// Load a value by key, batching and caching will be applied automatically.
//
// This method will block until max waiting time or max batch size is reached.
func (l *Dataloader[K, V]) Load(key K) (V, error) {
	return l.LoadThunk(key)()
}

// LoadThunk returns a function that when called will block waiting for a value.
//
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *Dataloader[K, V]) LoadThunk(key K) func() (V, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if it, ok := l.cache[key]; ok {
		return it
	}
	if l.batch == nil {
		l.batch = &batch[K, V]{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)

	thunk := func() (V, error) {
		<-batch.done

		var data V

		if len(batch.data) != len(batch.keys) {
			return data, fmt.Errorf("bug in loader: %d values returned for %d keys", len(batch.data), len(batch.keys))
		}

		if pos < len(batch.data) {
			data = batch.data[pos]
		}

		var err error
		// its convenient to be able to return a single error for everything
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		return data, err
	}
	l.cache[key] = thunk
	return thunk
}

// LoadMany fetches many keys at once. It will be broken into appropriate sized
// sub batches depending on how the loader is configured
func (l *Dataloader[K, V]) LoadMany(keys []K) ([]V, []error) {
	return l.LoadManyThunk(keys)()
}

// LoadManyThunk returns a function that when called will block waiting for a value.
//
// This method should be used if you want one goroutine to make requests to many
// different data loaders without blocking until the thunk is called.
func (l *Dataloader[K, V]) LoadManyThunk(keys []K) func() ([]V, []error) {
	thunks := make([]func() (V, error), len(keys))
	for i, key := range keys {
		thunks[i] = l.LoadThunk(key)
	}
	return func() ([]V, []error) {
		values := make([]V, len(keys))
		errors := make([]error, len(keys))
		noErrors := true
		for i, thunk := range thunks {
			values[i], errors[i] = thunk()
			if errors[i] != nil {
				noErrors = false
			}
		}
		if noErrors {
			return values, nil
		}
		return values, errors
	}
}

// keyIndex will return the location of the key in the batch,
// if its not found it will add the key to the batch
func (b *batch[K, V]) keyIndex(l *Dataloader[K, V], key K) int {
	for i, existingKey := range b.keys {
		if key == existingKey {
			return i
		}
	}

	pos := len(b.keys)
	b.keys = append(b.keys, key)
	if pos == 0 {
		go b.startTimer(l)
	}

	if l.maxSize != 0 && pos >= l.maxSize-1 {
		if !b.closing {
			b.closing = true
			l.batch = nil
			go b.end(l)
		}
	}

	return pos
}

func (b *batch[K, V]) startTimer(l *Dataloader[K, V]) {
	time.Sleep(l.maxWait)
	l.mu.Lock()

	// we must have hit a batch limit and are already finalizing this batch
	if b.closing {
		l.mu.Unlock()
		return
	}

	l.batch = nil
	l.mu.Unlock()

	b.end(l)
}

func (b *batch[K, V]) end(l *Dataloader[K, V]) {
	b.data, b.error = l.fetch(b.keys)
	close(b.done)
}
