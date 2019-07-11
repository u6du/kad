package main

import (
	"math/bits"
	"net"

	"github.com/spaolacci/murmur3"
)

const MaxDepth = 32
const BucketSize = 24

type Node struct {
	addr *net.UDPAddr
	id   uint32
}

type bucket [BucketSize]*Node

type kad struct {
	id     uint32
	bucket [MaxDepth]bucket
	depth  int
}

func New(id []byte) *kad {
	return &kad{id: murmur3.Sum32(id), depth: 0}
}

func (k *kad) depthIsFull(n int) bool {
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
	for _, node := range bucket {
		if k.Distance(node.id) > kDepth {
			next[nextPos] = node
			nextPos++
		} else {
			nowDepth[nowPos] = node
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

func (k *kad) Distance(hash uint32) int {
	return bits.OnesCount32(k.id ^ hash)
}

func (k *kad) AddNode(id []byte, addr *net.UDPAddr) bool {
	hash := murmur3.Sum32(id)
	depth := k.Distance(hash)

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
			t[i] = &Node{addr, hash}
			return true
		}
	}
	return false
}
