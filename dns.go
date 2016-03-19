package localkube

import (
	"fmt"
	"os"
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
	done     chan struct{}
}

func NewDNSServer(rootDomain, serverAddress, kubeAPIServer string) (*DNSServer, error) {
	// setup backing etcd store
	peerURLs := []string{"http://localhost:9256"}
	etcdServer, err := NewEtcd(DNSEtcdURLs, peerURLs, DNSName, DNSEtcdDataDirectory)
	if err != nil {
		return nil, err
	}

	// setup skydns
	etcdClient := etcd.NewClient(DNSEtcdURLs)
	skyConfig := &skydns.Config{
		DnsAddr: serverAddress,
		Domain:  rootDomain,
	}

	skydns.SetDefaults(skyConfig)

	backend := backendetcd.NewBackend(etcdClient, &backendetcd.Config{
		Ttl:      skyConfig.Ttl,
		Priority: skyConfig.Priority,
	})
	skyServer := skydns.New(backend, skyConfig)

	// setup so prometheus doesn't run into nil
	skydns.Metrics()

	// setup kube2sky
	k2s := kube2sky.NewKube2Sky(rootDomain, DNSEtcdURLs[0], "", kubeAPIServer, 10*time.Second, 8081)

	return &DNSServer{
		etcd:     etcdServer,
		sky:      skyServer,
		kube2sky: k2s,
	}, nil
}

func (dns *DNSServer) Start() {
	if dns.done != nil {
		fmt.Fprint(os.Stderr, pad("DNS server already started"))
		return
	}

	dns.done = make(chan struct{})

	dns.etcd.Start()
	go until(dns.kube2sky, os.Stderr, "kube2sky", 2*time.Second, dns.done)
	go until(dns.sky.Run, os.Stderr, "skydns", 1*time.Second, dns.done)
}

func (dns *DNSServer) Stop() {
	// closing chan will prevent servers from restarting but will not kill running server
	close(dns.done)

	dns.etcd.Stop()
}

// Status is currently not support by DNSServer
func (dns *DNSServer) Status() Status {
	if dns.done == nil {
		return Stopped
	}
	return Started
}

// Name returns the servers unique name
func (DNSServer) Name() string {
	return DNSName
}

// runner starts a server returning an error if it stops.
type runner interface {
	Run() error
}
