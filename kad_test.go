package kad

import (
	"crypto/rand"
	"testing"
)

func Testkad(t *testing.T) {
	id := [32]byte{}
	rand.Read(id[:])
	t.Logf("id = %x", id)
	kad := New(id)
	t.Logf("kad\n%s", kad)
}