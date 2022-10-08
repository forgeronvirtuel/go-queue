package queue

import (
	"bytes"
	"crypto/rand"
	"testing"
)

const dataSize = 2_000_000

func generateBytes() []byte {
	token := make([]byte, 4)
	rand.Read(token)
	return token
}

func TestQueue(t *testing.T) {
	var queue *Queue

	if msg, q := queue.Consume(); msg != nil || q != nil {
		t.Fatalf("msg = %v\nq = %v", msg, q)
	}
	if queue.Length() != 0 {
		t.Fatalf("not equal to zero")
	}

	msgs := [][]byte{
		generateBytes(),
		generateBytes(),
		generateBytes(),
	}

	queue = queue.
		Produce(msgs[0]).
		Produce(msgs[1]).
		Produce(msgs[2])

	if queue.Length() != 3 {
		t.Fatalf("Should be 3 items, got %d", queue.Length())
	}

	if msg, _ := queue.Consume(); bytes.Compare(msg, msgs[0]) != 0 {
		t.Fatalf("#1 Not equal")
	}
	if queue.length != 2 {
		t.Fatalf("Should be 2")
	}
	if msg, _ := queue.Consume(); bytes.Compare(msg, msgs[1]) != 0 {
		t.Fatalf("#2 Not equal")
	}
	if queue.length != 1 {
		t.Fatalf("Should be 1")
	}
	if msg, _ := queue.Consume(); bytes.Compare(msg, msgs[2]) != 0 {
		t.Fatalf("#3 Not equal")
	}
	if queue.length != 0 {
		t.Fatalf("Should be 0")
	}
	if msg, q := queue.Consume(); msg != nil || q != nil {
		t.Fatalf("msg = %v\nq = %v", msg, q)
	}
	if queue.length != 0 {
		t.Fatalf("Should be 0")
	}
}

func BenchmarkQueue_Produce(b *testing.B) {
	var data = make([][]byte, dataSize)
	for i := 0; i < dataSize; i++ {
		data[i] = generateBytes()
	}
	var queue *Queue
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < dataSize; i++ {
		queue = queue.Produce(data[i])
	}
	b.StopTimer()
}

func BenchmarkQueue_Consume(b *testing.B) {
	var data = make([][]byte, dataSize)
	for i := 0; i < dataSize; i++ {
		data[i] = generateBytes()
	}
	var queue *Queue
	for i := 0; i < dataSize; i++ {
		queue = queue.Produce(data[i])
	}

	b.ResetTimer()
	b.StartTimer()
	for queue.Length() != 0 {
		queue.Consume()
	}
	b.StopTimer()
}
