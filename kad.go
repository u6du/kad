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

/*
P2P 网络核心技术：Kademlia 协议
https://zhuanlan.zhihu.com/p/40286711

Kademlia：基于异或度量的点对点信息系统
http://t.cn/AipJQv3m

Kademlia2007
http://t.cn/AiWQWxWY
*/

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

const SplitLen = 64

func (k *Kad) add(now int, a *addr.Addr) bool {

	length := len(k.bucket) - 1
	if now > length {
		b := k.bucket[length]
		bLen := len(b)
		if bLen >= SplitLen {

			next := []*addr.Addr{a}
			var same []*addr.Addr
			for i := 0; i < bLen; i++ {
				if k.Similarity(b[i].Id) > length {
					next = append(next, b[i])
				} else {
					same = append(same, b[i])
				}
			}
			k.bucket[length] = same
			k.bucket = append(k.bucket, next)

			return true
		} else {
			now = length
		}
	}
	if len(k.bucket[now]) < SplitLen {
		k.bucket[now] = append(k.bucket[now], a)
		return true
	}
	return false
}

func (k *Kad) Delete(udp *net.UDPAddr) {
	k.lock.Lock()
	defer k.lock.Unlock()

	addrByte := udpaddr.Byte(udp)

	addr, ok := k.Ip.Get(addrByte)
	if ok {
		k.pop(k.bucketN(addr.Id), addr)
		k.Ip.Delete(addrByte)
	}
}

func (k *Kad) pop(n int, addr *addr.Addr) {
	bucket := k.bucket[n]
	for i := range bucket {
		if addr == bucket[i] {
			k.bucket[n] = append(bucket[:i], bucket[i+1:]...)
			break
		}
	}
}

// Tree通过addr*的byte映射到 secret，通过id映射到addr*

func (k *Kad) Add(id, secret [32]byte, udp *net.UDPAddr) bool {
	k.lock.Lock()
	defer k.lock.Unlock()

	addrByte := udpaddr.Byte(udp)
	addrPoint, exist := k.Ip.Get(addrByte)
	if exist {
		addrPoint.Secret = secret
		if bytes.Compare(addrPoint.Id[:], id[:]) != 0 {
			old := k.Similarity(addrPoint.Id)
			addrPoint.Id = id
			now := k.Similarity(id)
			if old != now {
				length := len(k.bucket) - 1
				if old > length {
					old = length
				}
				k.pop(old, addrPoint)
				k.add(now, addrPoint)
			}
		}
		return false
	} else {
		p := &addr.Addr{
			Secret: secret,
			Id:     id,
			Udp:    udp,
		}
		if k.add(k.Similarity(id), p) {
			k.Ip.Add(addrByte, p)
		}
		return true
	}
}

func (k *Kad) bucketN(id [32]byte) int {
	d := k.Similarity(id)
	length := len(k.bucket) - 1
	if d > length {
		d = length
	}
	return d
}

// 只返回和自己同样相似或者更加相似的节点
func (k *Kad) LookUp(id [32]byte) (li []*addr.Addr) {
	d := k.Similarity(id)
	length := len(k.bucket) - 1
	if d <= length {
		return k.bucket[d]
	}

	b := k.bucket[length]
	for i := range b {
		if Similarity(b[i].Id, id) >= d {
			li = append(li, b[i])
		}
	}
	return
}

func Similarity(idA [32]byte, idB [32]byte) int {
	return bits.LeadingZeros32(hash(idA) ^ hash(idB))
}

func (k *Kad) Similarity(id [32]byte) int {
	return bits.LeadingZeros32(k.id ^ hash(id))
}

func (k *Kad) String() string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	out := strings.Builder{}

	for i := 0; i < len(k.bucket); i++ {
		b := strings.Builder{}

		for _, node := range k.bucket[i] {
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
