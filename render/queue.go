// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
