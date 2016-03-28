package localkubectl

import (
	"fmt"

	kubecfg "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"
	kubectlcfg "k8s.io/kubernetes/pkg/kubectl/cmd/config"
)

// SetupContext creates a new cluster and context in ~/.kube using the provided API Server. If setCurrent is true, it is made the current context.
func SetupContext(clusterName, contextName, kubeAPIServer string, setCurrent bool) error {
	pathOpts := kubectlcfg.NewDefaultPathOptions()

	config, err := pathOpts.GetStartingConfig()
	if err != nil {
		return fmt.Errorf("could not setup config: %v", err)
	}

	cluster, exists := config.Clusters[clusterName]
	if !exists {
		cluster = kubecfg.NewCluster()
	}

	// configure cluster
	cluster.Server = kubeAPIServer
	cluster.InsecureSkipTLSVerify = true
	config.Clusters[clusterName] = cluster

	context, exists := config.Contexts[contextName]
	if !exists {
		context = kubecfg.NewContext()
	}

	// configure context
	context.Cluster = clusterName
	config.Contexts[contextName] = context

	// set as current if requested or no current context
	if len(config.CurrentContext) == 0 || config.CurrentContext == contextName || setCurrent {
		config.CurrentContext = contextName
	} else {
		fmt.Fprintf(Out, "kubectl is currently setup to use another context\n\n%s", SwitchContextInstructions(contextName))
	}

	return kubectlcfg.ModifyConfig(pathOpts, *config, true)
}

// SwitchContextInstructions
func SwitchContextInstructions(contextName string) string {
	l1 := "To setup kubectl to use localkube run:\n"
	l2 := fmt.Sprintf("kubectl config use-context %s\n", contextName)
	return l1 + l2
}
