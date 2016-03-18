package localkube

import (
	"time"

	kube2sky "rsprd.com/localkube/k2s"

	"github.com/coreos/go-etcd/etcd"
	backendetcd "github.com/skynetservices/skydns/backends/etcd"
	skydns "github.com/skynetservices/skydns/server"
)

const (
	DNSName = "dns"
)

var (
	DNSEtcdURLs = []string{"http://localhost:9090"}

	DNSEtcdDataDirectory = "/var/dns/data"
)

type DNSServer struct {
	etcd     *EtcdServer
	sky      runner
	kube2sky func() error
}

func NewDNSServer(rootDomain, kubeAPIServer string) (*DNSServer, error) {
	// setup backing etcd store
	peerURLs := []string{"http://localhost:9256"}
	etcdServer, err := NewEtcd(DNSEtcdURLs, peerURLs, DNSName, DNSEtcdDataDirectory)
	if err != nil {
		return nil, err
	}

	// setup skydns
	etcdClient := etcd.NewClient(DNSEtcdURLs)
	skyConfig := &skydns.Config{
		DnsAddr: "0.0.0.0:53",
		Domain:  rootDomain,
	}

	skydns.SetDefaults(skyConfig)

	backend := backendetcd.NewBackend(etcdClient, &backendetcd.Config{
		Ttl:      skyConfig.Ttl,
		Priority: skyConfig.Priority,
	})
	skyServer := skydns.New(backend, skyConfig)

	k2s := kube2sky.NewKube2Sky(rootDomain, DNSEtcdURLs[0], "", kubeAPIServer, 10*time.Second, 8081)

	return &DNSServer{
		etcd:     etcdServer,
		sky:      skyServer,
		kube2sky: k2s,
	}, nil
}

func (*DNSServer) Start() {}

func (*DNSServer) Stop() {
	println("DNS currently can't be stopped.")
}

// Status is currently not support by DNSServer
func (DNSServer) Status() Status {
	return NotImplemented
}

// Name returns the servers unique name
func (DNSServer) Name() string {
	return DNSName
}

// runner starts a server returning an error if it stops.
type runner interface {
	Run() error
}
