package utils

import (
	"container/heap"
	"errors"
	"github.com/pion/rtp"
)

type PriorityQueue struct {
	itemHeap    *packets
	current     *rtp.Packet
	priorityMap map[uint16]bool
	lastPacket  *rtp.Packet
}

func NewPq() *PriorityQueue {
	return &PriorityQueue{
		itemHeap:    &packets{},
		priorityMap: make(map[uint16]bool),
	}
}

func (p *PriorityQueue) Len() int {
	return p.itemHeap.Len()
}

func (p *PriorityQueue) Push(v rtp.Packet) {
	if p.priorityMap[v.SequenceNumber] {
		return
	}
	newItem := &packet{
		value:    v,
		priority: v.SequenceNumber,
	}
	heap.Push(p.itemHeap, newItem)
}

func (p *PriorityQueue) Pop() (rtp.Packet, error) {
	if len(*p.itemHeap) == 0 {
		return rtp.Packet{}, errors.New("empty queue")
	}

	item := heap.Pop(p.itemHeap).(*packet)
	return item.value, nil
}



func (p *PriorityQueue) Empty() {
	old := *p.itemHeap
	*p.itemHeap = old[:0]
}

type packets []*packet

type packet struct {
	value    rtp.Packet
	priority uint16
	index    int
}

func (p *packets) Len() int {
	return len(*p)
}

func (p *packets) Less(i, j int) bool {
	return (*p)[i].priority < (*p)[j].priority
}

func (p *packets) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
	(*p)[i].index = i
	(*p)[j].index = j
}

func (p *packets) Push(x interface{}) {
	it := x.(*packet)
	it.index = len(*p)
	*p = append(*p, it)
}

func (p *packets) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*p = old[0 : n-1]
	return item
}
