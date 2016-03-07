package main

import (
	"fmt"

	"github.com/coreos/etcd/etcdserver"
	"github.com/coreos/etcd/pkg/types"
)

var (
	ClientURLs = []string{"http://localhost:2379", "http://localhost:4001"}
	PeerURLs   = []string{"http://localhost:2380"}
)

func Etcd(dataDir string) (*etcdserver.EtcdServer, error) {
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
		return nil, fmt.Errorf("Etcd config error: %v", err)
	}

	return server, nil
}

func urlsOrPanic(urlStrs []string) types.URLs {
	urls, err := types.NewURLs(urlStrs)
	if err != nil {
		panic(err)
	}
	return urls
}
