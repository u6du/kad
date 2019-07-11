package main

import (
	"math/bits"
	"net"

	"github.com/spaolacci/murmur3"
)

const MaxDepth = 32
const BucketSize = 24

type bucket [BucketSize]*net.UDPAddr

type kad struct {
	id     uint32
	bucket [MaxDepth]bucket
	depth  uint16
	Node   map[*net.UDPAddr][32]byte
}

func New(id [32]byte) *kad {
	return &kad{id: murmur3.Sum32(id[:]), depth: 0}
}

func (k *kad) depthIsFull(n uint16) bool {
	return k.bucket[n][BucketSize-1] != nil
}

func (k *kad) split() bool {
	nowDepth := bucket{}
	nowPos := 0
	next := bucket{}
	nextPos := 0

	kDepth := k.depth

	if kDepth == MaxDepth {
		return false
	}
	bucket := k.bucket[kDepth]
	for _, addr := range bucket {
		if k.Distance(k.Node[addr]) > kDepth {
			next[nextPos] = addr
			nextPos++
		} else {
			nowDepth[nowPos] = addr
			nowPos++
		}
	}
	if nextPos > 0 {
		k.bucket[kDepth] = nowDepth
		k.depth++
		k.bucket[k.depth] = next
		return true
	} else {
		return false
	}
}

func (k *kad) Distance(id [32]byte) uint16 {
	hash := murmur3.Sum32(id[:])
	return uint16(bits.OnesCount32(k.id ^ hash))
}

func (k *kad) AddNode(id [32]byte, addr *net.UDPAddr) bool {
	hashDepth := k.Distance(id)
	depth := hashDepth

	kDepth := k.depth
	if depth > kDepth {
		if k.depthIsFull(kDepth) {
			if !k.split() {
				return false
			}
		}
		depth = k.depth
	}
	t := k.bucket[depth]
	if t[BucketSize-1] != nil {
		return false
	}
	for i := range t {
		if t[i] == nil {
			t[i] = addr
			return true
		}
	}
	return false
}
