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

	for i := 0; i < 9999; i++ {
		id2 := [32]byte{}
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

	rand.Read(id[:])

	near = kad.LookUp(id)
	for i := range near {
		addr := near[i]
		t.Logf("%d %s", i, addr.Udp.String())
	}

	total := 0
	for i := range kad.bucket {
		total += len(kad.bucket[i])
	}
	if uint(total) != kad.Len() {
		t.Error("kad len != sum len(kad.bucket)")
	}

}
