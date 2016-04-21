package batchposter

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	ErrBufferFull = errors.New("buffer full")
)

type BatchPoster struct {
	url        string        // the url to POST to when flushing the buffer
	buffersize int           // the maximum amount of posts that can be batched
	autoflush  time.Duration // max time a buffer can fillup before flushing
	log        *log.Logger   // where errors get reported to

	mu     sync.RWMutex
	buffer []string
	idx    int
}

func New(url string, maxbatch int, maxtime time.Duration, l *log.Logger) *BatchPoster {

	obj := &BatchPoster{url: url, buffersize: maxbatch, autoflush: maxtime, log: l}
	obj.idx = 0
	obj.buffer = make([]string, obj.buffersize)

	go func() {
		for {
			// TODO: Add ability to start/stop autoflushing
			time.Sleep(obj.autoflush)
			obj.flush()
		}
	}()

	return obj
}

func (b *BatchPoster) Post(payload string) error {

	if b.full() {

		// something is wrong, buffer should never be full. Try flushing and ignore errors.
		go b.flush()

		b.log.Println(ErrBufferFull)
		return ErrBufferFull
	}

	err := b.add(payload)
	if err != nil {
		return err
	}

	if b.full() {
		err := b.flush()
		if err != nil {
			return err
		}
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////
// private parts
//////////////////////////////////////////////////////////////////////////

func (b *BatchPoster) full() bool {

	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.idx >= b.buffersize
}

func (b *BatchPoster) add(payload string) error {

	b.mu.Lock()
	defer b.mu.Unlock()

	// add to buffer
	b.buffer[b.idx] = payload
	b.idx++

	return nil
}

func (b *BatchPoster) flush() error {

	b.mu.Lock()
	defer b.mu.Unlock()

	// bail on nothing to post
	if b.idx == 0 {
		return nil
	}

	// post batch
	batch := strings.Join(b.buffer, "\n")
	_, err := http.Post(b.url, "", strings.NewReader(batch))
	if err != nil {
		b.log.Println(err)
		return err
	}

	// reset
	for i := 0; i < b.buffersize; i++ {
		b.buffer[i] = ""
	}
	b.idx = 0

	return nil
}
