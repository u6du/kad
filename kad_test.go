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

	kad := New(id)

	var id2 [32]byte
	ip := make([]byte, 6)
	for i := 0; i < 9999; i++ {
		_, err = rand.Read(id2[:])
		ex.Panic(err)

		secret := [32]byte{}

		_, err = rand.Read(secret[:])
		ex.Panic(err)

		_, err = rand.Read(ip)
		ex.Panic(err)

		kad.Add(id2, secret, udpaddr.Addr(ip))
	}
	t.Logf("kad\n%s", kad)

	total := kad.Len()

	var id3 [32]byte
	kad.Add(id3, id3, udpaddr.Addr(ip))
	if kad.Len() != total {
		t.Error("重复添加相同的ip端口，应该不会增加总长度")
	}

	near := kad.LookUp(id2)
	if len(near) <= 0 {
		t.Error("应该能找到相似的节点")
	}

	near = kad.LookUp(id)
	if len(near) != 0 {
		t.Error("lookup 自己的id应该是返回空列表（只返回和自己同样相似或者更加相似的节点）")
	}

	total = 0
	for i := range kad.bucket {
		total += uint(len(kad.bucket[i]))
	}
	if total != kad.Len() {
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
