package socket

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"
)

func generateBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func TestBufferReader_Next_OneCommand(t *testing.T) {
	reader := bytes.NewBuffer([]byte("add 10;"))
	buf := NewBufferReader(1024, reader)
	res, err := buf.Next(';')
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(res, []byte("add 10;")) != 0 {
		t.Errorf("%s != add 10", res)
	}
}

func TestBufferReader_Next_MultipleCommand(t *testing.T) {
	reader := bytes.NewBuffer([]byte("add 10;close;"))
	buf := NewBufferReader(1024, reader)
	res, err := buf.Next(';')
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(res, []byte("add 10;")) != 0 {
		t.Errorf("%s != add 10", res)
	}
	res, err = buf.Next(';')
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(res, []byte("close;")) != 0 {
		t.Errorf("%s != close;", res)
	}
}

func TestBufferReader_Next_Overflow(t *testing.T) {
	reader := bytes.NewBuffer([]byte("add 10;close;"))
	buf := NewBufferReader(3, reader)
	_, err := buf.Next(';')
	if err == nil {
		t.Fatalf("f should not be nil")
	}
	if err.Error() != err_overflow_cap {
		t.Fatal("not the good message")
	}
}

// TODO test in simulated network
func TestBufferReader_Next_Delayed(t *testing.T) {
	ch := make(chan bool)
	endch := make(chan bool)
	reader := bytes.NewBuffer([]byte("add"))
	buf := NewBufferReader(10, &DelayedBuffer{
		buf: reader,
		ch:  ch,
	})

	go func() {
		res, err := buf.Next(';')
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare([]byte("add;"), res) != 0 {
			t.Errorf("%s != add 10;", res)
		}
		res, err = buf.Next(';')
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare([]byte("close;"), res) != 0 {
			t.Errorf("%s != close;", res)
		}
		res, err = buf.Next(';')
		if err == nil {
			t.Fatal("Should not be nil")
		}
		if err.Error() != err_overflow_cap {
			t.Fatal("Not the same message")
		}
		endch <- true
	}()

	ch <- true
	time.Sleep(100 * time.Millisecond)
	reader.Write([]byte(";cl"))

	ch <- true
	time.Sleep(100 * time.Millisecond)
	reader.Write([]byte("ose;an u"))

	ch <- true
	time.Sleep(100 * time.Millisecond)
	reader.Write([]byte("nexpected"))

	if !(<-endch) {
		t.Fatal("not finished")
	}
}

type DelayedBuffer struct {
	buf *bytes.Buffer
	ch  chan bool
}

func (d *DelayedBuffer) Read(p []byte) (n int, err error) {
	<-d.ch
	return d.buf.Read(p)
}

// TODO test in simulated network
func TestBufferReader_Next_Delayed_other(t *testing.T) {
	ch := make(chan bool)
	endch := make(chan bool)
	reader := bytes.NewBuffer([]byte("add"))
	buf := NewBufferReader(10, &DelayedBuffer{
		buf: reader,
		ch:  ch,
	})

	go func() {
		res, err := buf.Next(';')
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare([]byte("add;"), res) != 0 {
			t.Errorf("%s != add 10;", res)
		}
		endch <- true
	}()

	ch <- true
	time.Sleep(100 * time.Millisecond)
	reader.Write([]byte(";cl"))

	if !(<-endch) {
		t.Fatal("not finished")
	}
}

// TODO test with a lot of Next step to check on indefinitely growing buffer
