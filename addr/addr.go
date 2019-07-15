package addr

import "net"

type Addr struct {
	Secret [32]byte
	Id [32]byte
	Udp *net.UDPAddr
}
