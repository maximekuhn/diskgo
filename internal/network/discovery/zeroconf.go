package discovery

import (
	"context"
	"fmt"
	"github.com/grandcat/zeroconf"
	"github.com/maximekuhn/diskgo/internal/network"
	"net/netip"
	"strings"
)

const (
	serviceName = "diskgo"
	domain      = "local."
)

type ZeroconfResolver struct {
}

func NewZeroconfResolver() *ZeroconfResolver {
	return &ZeroconfResolver{}
}

func (z *ZeroconfResolver) Resolve(ctx context.Context, peers chan<- network.Peer) error {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for {
			select {
			case entry := <-results:
				if entry.Service != serviceName || entry.Domain != domain {
					continue
				}

				texts := entry.Text

				if len(texts) != 2 {
					continue
				}

				// first text should be nickname
				nicknameText := texts[0]
				if !strings.HasPrefix(nicknameText, "nickname=") {
					continue
				}
				nickname := strings.TrimPrefix(nicknameText, "nickname=")

				// second text should be addr
				addrText := texts[1]
				if !strings.HasPrefix(addrText, "addr=") {
					continue
				}
				addr := strings.TrimPrefix(addrText, "addr=")
				addrPort, err := netip.ParseAddrPort(addr)
				if err != nil {
					continue
				}

				peer := network.Peer{
					Name: nickname,
					Addr: addrPort,
				}

				peers <- peer

			case <-ctx.Done():
				return
			}
		}
	}(entries)

	err = resolver.Browse(ctx, serviceName, domain, entries)
	if err != nil {
		return err
	}

	return nil
}

type ZeroConfAdvertiser struct {
	Nickname string
	Addr     netip.AddrPort
}

func NewZeroConfAdvertiser(nickname string, addr netip.AddrPort) *ZeroConfAdvertiser {
	return &ZeroConfAdvertiser{Nickname: nickname, Addr: addr}
}

func (z *ZeroConfAdvertiser) Advertise(ctx context.Context) error {
	instance := fmt.Sprintf("%s-%d", z.Nickname, z.Addr.Port())
	txtNickname := fmt.Sprintf("nickname=%s", z.Nickname)
	txtListenAddr := fmt.Sprintf("addr=%s", z.Addr.String())
	server, err := zeroconf.Register(instance, serviceName, domain, int(z.Addr.Port()), []string{txtNickname, txtListenAddr}, nil)
	if err != nil {
		return err
	}
	defer server.Shutdown()

	<-ctx.Done()

	return nil
}
