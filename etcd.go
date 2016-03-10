package localkube

import (
	"fmt"

	"github.com/coreos/etcd/etcdserver"
	"github.com/coreos/etcd/pkg/types"
)

const (
	EtcdName = "etcd"
)

var (
	ClientURLs = []string{"http://localhost:2379", "http://localhost:4001"}
	PeerURLs   = []string{"http://localhost:2380"}
)

// Etcd is a Server which manages an Etcd cluster
type Etcd struct {
	*etcdserver.EtcdServer
}

// NewEtcd creates a new default etcd Server using 'dataDir' for persistence. Panics if could not be configured.
func NewEtcd(dataDir string) *Etcd {
	name := "default"
	clientURLs := urlsOrPanic(ClientURLs)
	peerURLs := urlsOrPanic(PeerURLs)

	urlsMap := map[string]types.URLs{
		name: peerURLs,
	}

	config := &etcdserver.ServerConfig{
		Name:                name,
		ClientURLs:          clientURLs,
		PeerURLs:            peerURLs,
		DataDir:             dataDir,
		InitialClusterToken: "etcd-cluster",
		InitialPeerURLsMap:  urlsMap,

		NewCluster: true,

		SnapCount:     etcdserver.DefaultSnapCount,
		MaxSnapFiles:  5,
		MaxWALFiles:   5,
		TickMs:        100,
		ElectionTicks: 10,
	}

	server, err := etcdserver.NewServer(config)
	if err != nil {
		msg := fmt.Sprintf("Etcd config error: %v", err)
		panic(msg)
	}

	return &Etcd{
		EtcdServer: server,
	}
}

// Status is currently not support by Etcd
func (Etcd) Status() Status {
	return NotImplemented
}

// Name returns the servers unique name
func (Etcd) Name() string {
	return EtcdName
}

func urlsOrPanic(urlStrs []string) types.URLs {
	urls, err := types.NewURLs(urlStrs)
	if err != nil {
		panic(err)
	}
	return urls
}
