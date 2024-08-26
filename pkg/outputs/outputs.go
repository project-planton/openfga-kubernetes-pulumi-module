package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/openfgakubernetes"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	Namespace               = "namespace"
	Service                 = "service"
	PortForwardCommand      = "port-forward-command"
	KubeEndpoint            = "kube-endpoint"
	IngressExternalHostname = "ingress-external-hostname"
	IngressInternalHostname = "ingress-internal-hostname"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *openfgakubernetes.OpenfgaKubernetesStackInput) *openfgakubernetes.OpenfgaKubernetesStackOutputs {
	return &openfgakubernetes.OpenfgaKubernetesStackOutputs{
		Namespace:          autoapistackoutput.GetVal(pulumiOutputs, Namespace),
		Service:            autoapistackoutput.GetVal(pulumiOutputs, Service),
		PortForwardCommand: autoapistackoutput.GetVal(pulumiOutputs, PortForwardCommand),
		KubeEndpoint:       autoapistackoutput.GetVal(pulumiOutputs, KubeEndpoint),
		ExternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, IngressExternalHostname),
		InternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, IngressInternalHostname),
	}
}
