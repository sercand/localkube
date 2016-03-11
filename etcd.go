package localkube

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/coreos/etcd/etcdserver"
	"github.com/coreos/etcd/etcdserver/etcdhttp"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/coreos/etcd/pkg/types"
)

const (
	EtcdName = "etcd"
)

var (
	ClientURLs = []string{"http://localhost:2379"}
	PeerURLs   = []string{"http://localhost:2380"}
)

// Etcd is a Server which manages an Etcd cluster
type Etcd struct {
	*etcdserver.EtcdServer
	clientListens []net.Listener
}

// NewEtcd creates a new default etcd Server using 'dataDir' for persistence. Panics if could not be configured.
func NewEtcd(dataDir string) *Etcd {
	name := "default"
	clientURLs := urlsOrPanic(ClientURLs)
	peerURLs := urlsOrPanic(PeerURLs)

	urlsMap := map[string]types.URLs{
		name: peerURLs,
	}

	clientListeners := createListenersOrPanic(clientURLs)

	config := &etcdserver.ServerConfig{
		Name:                name,
		ClientURLs:          clientURLs,
		PeerURLs:            peerURLs,
		DataDir:             dataDir,
		InitialClusterToken: "etcd-cluster",
		InitialPeerURLsMap:  urlsMap,
		Transport:           http.DefaultTransport.(*http.Transport),

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

	// setup client listeners
	ch := etcdhttp.NewClientHandler(server, config.ReqTimeout())
	for _, l := range clientListeners {
		go func(l net.Listener) {
			srv := &http.Server{
				Handler:     ch,
				ReadTimeout: 5*time.Minute,
			}
			panic(srv.Serve(l))
		}(l)
	}

	return &Etcd{
		EtcdServer:    server,
		clientListens: clientListeners,
	}
}

// Stop closes all connections and stops the Etcd server
func (e *Etcd) Stop() {
	e.EtcdServer.Stop()
	for _, l := range e.clientListens {
		l.Close()
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

func createListenersOrPanic(urls types.URLs) (listeners []net.Listener) {
	for _, url := range urls {
		l, err := net.Listen("tcp", url.Host)
		if err != nil {
			panic(err)
		}

		l, err = transport.NewKeepAliveListener(l, url.Scheme, transport.TLSInfo{})
		if err != nil {
			panic(err)
		}

		listeners = append(listeners, l)
	}
	return listeners
}