package kad

import (
	"crypto/rand"
	"net"
	"testing"

	"github.com/u6du/ex"
	"github.com/u6du/udpaddr"
)

func TestKad_Add(t *testing.T) {
	id := [32]byte{}
	_, err := rand.Read(id[:])
	ex.Panic(err)

	t.Logf("id = %x", id)
	kad := New(id)
	t.Logf("kad\n%s", kad)

	var id2 [32]byte
	for i := 0; i < 9999; i++ {
		_, err = rand.Read(id2[:])
		ex.Panic(err)

		secret := [32]byte{}

		_, err = rand.Read(secret[:])
		ex.Panic(err)

		ip := make([]byte, 6)
		_, err = rand.Read(ip)
		ex.Panic(err)

		kad.Add(id2, secret, udpaddr.Addr(ip))
	}
	t.Logf("kad\n%s", kad)

	near := kad.LookUp(id2)
	if len(near) <= 0 {
		t.Error("应该能找到相似的节点")
	}

	near = kad.LookUp(id)
	if len(near) != 0 {
		t.Error("lookup 自己的id应该是返回空列表（只返回和自己同样相似或者更加相似的节点）")
	}

	total := 0
	for i := range kad.bucket {
		total += len(kad.bucket[i])
	}
	if uint(total) != kad.Len() {
		t.Error("kad len != sum len(kad.bucket)")
	}

	var addr *net.UDPAddr

	for i := range kad.bucket {
		addr = kad.bucket[i][0].Udp
		break
	}

	kad.Delete(addr)
	if kad.Len() != uint(total-1) {
		t.Error("Delete 应该删除一个地址")
	}

}
