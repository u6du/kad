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


type Kad struct {
	id     uint32
	bucket [][]*addr.Addr
	Ip radixmapaddr.Tree
}

func hash(id [32]byte) uint32{
	return murmur3.Sum32(id[:])
}

func New(id [32]byte) *Kad {
	return &Kad{
		id: hash(id),
		Ip:*radixmapaddr.New(),
	}
}

// Tree通过addr*的byte映射到 secret，通过id映射到addr*

func (k *Kad) Add(id,secret [32]byte, udp *net.UDPAddr) bool {
	addrByte := udpaddr.Byte(udp)
	addrExist, exist := k.Ip.Get(addrByte)
	if exist {
		addrExist.Secret = secret
		if bytes.Compare(addrExist.Id[:], id[:]) != 0 {
			old := k.Distance(addrExist.Id)
			addrExist.Id = id
			now := k.Distance(id)
			if old != now {
				bucketOld := k.bucket[old]
				for i := range bucketOld {
					if addrExist == bucketOld[i] {
						k.bucket[old] = append(bucketOld[:i], bucketOld[i+1:]...)
					}
				}
				k.bucket[now] = append(k.bucket[now], addrExist)
			}
		}
		return false
	} else {
		now := k.Distance(id)
		p := &addr.Addr{
			Secret: secret,
			Id:     id,
			Udp:    udp,
		}
		k.bucket[now] = append(k.bucket[now], p)
		k.Ip.Add(addrByte, p)
		return true
	}

}

func (k *Kad) Distance(id [32]byte) uint16 {
	return uint16(bits.OnesCount32(k.id ^ hash(id)))
}

func (k *Kad) String() string {
	out := strings.Builder{}
	for i := uint16(0); i < 32; i++ {
		b := strings.Builder{}

		for _, node := range k.bucket[i] {
			if node == nil {
				break
			}
			b.WriteString(" ")
			b.WriteString(node.Udp.String())
		}
		s := b.String()
		if len(s) > 0{
			out.WriteString(fmt.Sprintf("%d :", i))
			out.WriteString(s)
			out.WriteString("\n")
		}
	}
	return out.String()
}


func (k *Kad) Len() uint{
	return k.Ip.Len()
}