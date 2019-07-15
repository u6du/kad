package kad

import (
	"bytes"
	"fmt"
	"math/bits"
	"net"
	"strings"
	"sync"

	"github.com/spaolacci/murmur3"
	"github.com/u6du/udpaddr"

	"github.com/u6du/kad/addr"
	"github.com/u6du/kad/radixmapaddr"
)

type Kad struct {
	id     uint32
	bucket [][]*addr.Addr
	Ip     radixmapaddr.Tree
	lock   sync.RWMutex
}

func hash(id [32]byte) uint32 {
	return murmur3.Sum32(id[:])
}

func New(id [32]byte) *Kad {
	return &Kad{
		id:     hash(id),
		Ip:     *radixmapaddr.New(),
		bucket: [][]*addr.Addr{{}},
	}
}

const SplitLen = 32

func (k *Kad) add(now int, a *addr.Addr) {

	length := len(k.bucket) - 1
	if now > length {
		b := k.bucket[length]
		bLen := len(b)
		if bLen >= SplitLen {

			next := []*addr.Addr{a}
			var same []*addr.Addr
			for i := 0; i < bLen; i++ {
				if k.Distance(b[i].Id) > length {
					next = append(next, b[i])
				} else {
					same = append(same, b[i])
				}
			}
			k.bucket[length] = same
			k.bucket = append(k.bucket, next)

			return
		} else {
			now = length
		}
	}
	if len(k.bucket[now]) < SplitLen {
		k.bucket[now] = append(k.bucket[now], a)
	}
}

// Tree通过addr*的byte映射到 secret，通过id映射到addr*

func (k *Kad) Add(id, secret [32]byte, udp *net.UDPAddr) bool {
	k.lock.Lock()
	defer k.lock.Unlock()

	addrByte := udpaddr.Byte(udp)
	addrExist, exist := k.Ip.Get(addrByte)
	if exist {
		addrExist.Secret = secret
		if bytes.Compare(addrExist.Id[:], id[:]) != 0 {
			old := k.Distance(addrExist.Id)
			addrExist.Id = id
			now := k.Distance(id)
			if old != now {
				length := len(k.bucket) - 1
				if old > length {
					old = length
				}
				bucketOld := k.bucket[old]
				for i := range bucketOld {
					if addrExist == bucketOld[i] {
						k.bucket[old] = append(bucketOld[:i], bucketOld[i+1:]...)
					}
				}
				k.add(now, addrExist)
			}
		}
		return false
	} else {
		p := &addr.Addr{
			Secret: secret,
			Id:     id,
			Udp:    udp,
		}
		k.add(k.Distance(id), p)
		k.Ip.Add(addrByte, p)
		return true
	}
}

func (k *Kad) Distance(id [32]byte) int {
	return bits.OnesCount32(k.id ^ hash(id))
}

func (k *Kad) String() string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	out := strings.Builder{}

	for i := 0; i < len(k.bucket); i++ {
		b := strings.Builder{}

		for _, node := range k.bucket[i] {
			if node == nil {
				break
			}
			b.WriteString(" ")
			b.WriteString(node.Udp.String())
		}
		s := b.String()
		if len(s) > 0 {
			out.WriteString(fmt.Sprintf("%d :", i))
			out.WriteString(s)
			out.WriteString("\n")
		}
	}
	return out.String()
}

func (k *Kad) Len() uint {
	return k.Ip.Len()
}
