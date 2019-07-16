package kad

import (
	"crypto/rand"
	"net"
	"testing"
)

func TestKad_Add(t *testing.T) {
	id := [32]byte{}
	rand.Read(id[:])
	t.Logf("id = %x", id)
	kad := New(id)
	t.Logf("kad\n%s", kad)

	var id2 [32]byte
	for i := 0; i < 9999; i++ {
		rand.Read(id2[:])
		secret := [32]byte{}
		rand.Read(secret[:])
		ip := make([]byte, 4)
		rand.Read(ip)
		kad.Add(id2, secret, &net.UDPAddr{IP: ip, Port: 3232})
	}
	t.Logf("kad\n%s", kad)

	near := kad.LookUp(id)
	if len(near) != 0 {
		t.Error("lookup 自己的id应该是返回空列表（只返回和自己同样相似或者更加相似的节点）")
	}

	near = kad.LookUp(id2)
	if len(near) <= 0 {
		t.Error("应该能找到相似的节点")
	}

	total := 0
	for i := range kad.bucket {
		total += len(kad.bucket[i])
	}
	if uint(total) != kad.Len() {
		t.Error("kad len != sum len(kad.bucket)")
	}

}
