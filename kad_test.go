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

	for i := 0; i < 99999; i++ {
		id2 := [32]byte{}
		rand.Read(id2[:])
		secret := [32]byte{}
		rand.Read(secret[:])
		ip := make([]byte, 4)
		rand.Read(ip)
		kad.Add(id2, secret, &net.UDPAddr{IP: ip, Port: 3232})
	}
	t.Logf("kad\n%s", kad)

}
