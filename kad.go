package kad

import (
	"bytes"
	"fmt"
	"math/bits"
	"net"
	"strings"

	"github.com/spaolacci/murmur3"
	"github.com/u6du/udpaddr"

	"github.com/u6du/kad/addr"
	"github.com/u6du/kad/radixmapaddr"
)

const MaxDepth = 32

type Kad struct {
	id     uint32
	bucket [MaxDepth][]*addr.Addr
	Ip radixmapaddr.Tree
}

func New(id [32]byte) *Kad {
	return &Kad{
		id: murmur3.Sum32(id[:]),
		Ip:*radixmapaddr.New(),
	}
}

// Tree通过addr*的byte映射到 secret，通过id映射到addr*

func (k *Kad) AddNode(id,secret [32]byte, udp *net.UDPAddr) bool {
	addrByte := udpaddr.Byte(udp)
	addrExist, exist := k.Ip.Get(addrByte)
	if exist{
		addrExist.Secret = secret
		if bytes.Compare(addrExist.Id[:], id[:])!=0{
			old := k.Distance(addrExist.Id)
			addrExist.Id = id
			now := k.Distance(id)
			if old!=now{
				bucketOld := k.bucket[old]
				for i := range bucketOld{
					if addrExist == bucketOld[i]{
						k.bucket[old]=append(bucketOld[:i], bucketOld[i+1:]...)
					}
				}
				k.bucket[now] = append(k.bucket[now],addrExist)
			}
		}
		return false
	}else {
		now := k.Distance(id)
		p:=&addr.Addr{
			Secret: secret,
			Id:     id,
			Udp:    udp,
		}
		k.bucket[now] = append(k.bucket[now], p)
		k.Ip.Add(addrByte, p)
		return true
	}

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
