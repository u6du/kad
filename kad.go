package kad

import (
	"fmt"
	"math/bits"
	"net"
	"strings"

	"github.com/spaolacci/murmur3"
	"github.com/u6du/radix/radixset"
)

const MaxDepth = 32
const BucketSize = 36

type Bucket [BucketSize]*net.UDPAddr

type Kad struct {
	id     uint32
	bucket [MaxDepth]Bucket
	Tree   *radixset.Tree
}

func New(id [32]byte) *Kad {
	return &Kad{id: murmur3.Sum32(id[:]), Tree: radixset.New()}
}

func (k *Kad) AddNode(id [32]byte, addr *net.UDPAddr) bool {
	if k.Tree.Get(id[:]) {
		return false
	}

	depth := k.Distance(id)

	t := k.bucket[depth]

	for i := range t {
		if t[i] == nil {
			t[i] = addr
			return true
		}
	}

	return false
}

func (k *Kad) depthIsFull(n uint16) bool {
	return k.bucket[n][BucketSize-1] != nil
}

func (k *Kad) Distance(id [32]byte) uint16 {
	hash := murmur3.Sum32(id[:])
	return uint16(bits.OnesCount32(k.id ^ hash))
}

func (k *Kad) String() string {
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
