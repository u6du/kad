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

	id2 := [32]byte{}
	rand.Read(id2[:])
	secret := [32]byte{}
	rand.Read(secret[:])

	kad.Add(id2, secret, &net.UDPAddr{IP: []byte("1234"), Port: 3232})
	t.Logf("kad\n%s", kad)
	kad.Add(id2, secret, &net.UDPAddr{IP: []byte("5431"), Port: 3232})
	t.Logf("kad\n%s", kad)
	kad.Add(secret, id, &net.UDPAddr{IP: []byte("1234"), Port: 3232})
	t.Logf("kad\n%s", kad)

}
