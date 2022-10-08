package queue

type bytesLinkedListItem struct {
	msg []byte
	ll  *bytesLinkedListItem
}

type Queue struct {
	first  *bytesLinkedListItem
	last   *bytesLinkedListItem
	length int
}

func (q *Queue) Produce(msg []byte) *Queue {
	if q == nil {
		blk := &bytesLinkedListItem{
			msg: msg,
			ll:  nil,
		}
		return &Queue{
			first:  blk,
			last:   blk,
			length: 1,
		}
	}

	blk := &bytesLinkedListItem{
		msg: msg,
		ll:  nil,
	}
	q.last.ll = blk
	q.last = blk
	q.length++
	return q
}

func (q *Queue) Consume() ([]byte, *Queue) {
	if q == nil || q.first == nil {
		return nil, nil
	}
	blk := q.first
	q.first = blk.ll
	q.length--
	return blk.msg, q
}

func (q *Queue) Length() int {
	if q == nil {
		return 0
	}
	return q.length
}
