package main

import (
	"fmt"
	"math/bits"
	"net"
	"strings"

	"github.com/spaolacci/murmur3"
)

const MaxDepth = 32
const BucketSize = 24

type bucket [BucketSize]*net.UDPAddr

type kad struct {
	id     uint32
	bucket [MaxDepth]bucket
	Node   map[*net.UDPAddr][32]byte
}

func New(id [32]byte) *kad {
	return &kad{id: murmur3.Sum32(id[:])}
}

func (k *kad) depthIsFull(n uint16) bool {
	return k.bucket[n][BucketSize-1] != nil
}

func (k *kad) Distance(id [32]byte) uint16 {
	hash := murmur3.Sum32(id[:])
	return uint16(bits.OnesCount32(k.id ^ hash))
}

func (k *kad) String() string {
	b := strings.Builder{}
	for i := uint16(0); i < MaxDepth; i++ {
		b.WriteString(fmt.Sprintf("%d :", i))
		for _, node := range k.bucket[i] {
			if node == nil {
				break
			}
			b.WriteString(" ")
			b.WriteString(node.String())
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (k *kad) AddNode(id [32]byte, addr *net.UDPAddr) bool {
	depth := k.Distance(id)

	t := k.bucket[depth]

	if t[BucketSize-1] == nil {
		for i := range t {
			if t[i] == nil {
				t[i] = addr
				return true
			}
		}
	}

	return false
}
