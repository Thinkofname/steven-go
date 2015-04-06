package render

const (
	maxQueueSize = 1 << 12
	queueMask    = maxQueueSize - 1
)

type renderQueue struct {
	queue        [maxQueueSize]renderRequest
	startPointer int
	endPointer   int
}

func (rq *renderQueue) Append(r renderRequest) {
	rq.queue[rq.endPointer] = r
	rq.endPointer = (rq.endPointer + 1) & queueMask
}

func (rq *renderQueue) Take() renderRequest {
	val := rq.queue[rq.startPointer]
	rq.startPointer = (rq.startPointer + 1) & queueMask
	return val
}

func (rq *renderQueue) Empty() bool {
	return rq.startPointer == rq.endPointer
}
