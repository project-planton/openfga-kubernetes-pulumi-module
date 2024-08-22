package pkg

import (
	"fmt"
	"github.com/plantoncloud/openfga-kubernetes-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/openfgakubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	OpenfgaKubernetes                 *openfgakubernetes.OpenfgaKubernetes
	Namespace                         string
	IngressCertClusterIssuerName      string
	IngressCertSecretName             string
	IngressHttpInternalHostname       string
	IngressHttpExternalHostname       string
	IngressGrpcInternalHostname       string
	IngressGrpcExternalHostname       string
	IngressPlaygroundInternalHostname string
	IngressPlaygroundExternalHostname string
	IngressMetricsInternalHostname    string
	IngressMetricsExternalHostname    string
	IngressHostnames                  []string
	KubeServiceFqdn                   string
	KubeServiceName                   string
	KubePortForwardCommand            string
}

func initializeLocals(ctx *pulumi.Context, stackInput *openfgakubernetes.OpenfgaKubernetesStackInput) *Locals {
	locals := &Locals{}
	//assign value for the local variable to make it available across the project
	locals.OpenfgaKubernetes = stackInput.ApiResource

	jenkinsKubernetes := stackInput.ApiResource

	//decide on the namespace
	locals.Namespace = jenkinsKubernetes.Metadata.Id

	locals.KubeServiceName = jenkinsKubernetes.Metadata.Name

	//export kubernetes service name
	ctx.Export(outputs.Service, pulumi.String(locals.KubeServiceName))

	locals.KubeServiceFqdn = fmt.Sprintf("%s.%s.svc.cluster.local",
		jenkinsKubernetes.Metadata.Name, locals.Namespace)

	//export kubernetes endpoint
	ctx.Export(outputs.KubeEndpoint, pulumi.String(locals.KubeServiceFqdn))

	locals.KubePortForwardCommand = fmt.Sprintf("kubectl port-forward -n %s service/%s 8080:8080",
		locals.Namespace, jenkinsKubernetes.Metadata.Name)

	//export kube-port-forward command
	ctx.Export(outputs.PortForwardCommand, pulumi.String(locals.KubePortForwardCommand))

	if jenkinsKubernetes.Spec.Ingress == nil ||
		!jenkinsKubernetes.Spec.Ingress.IsEnabled ||
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return locals
	}

	locals.IngressHttpExternalHostname = fmt.Sprintf("%s-http.%s",
		jenkinsKubernetes.Metadata.Id, jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressHttpInternalHostname = fmt.Sprintf("%s-http-internal.%s", jenkinsKubernetes.Metadata.Id,
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressGrpcExternalHostname = fmt.Sprintf("%s-grpc.%s",
		jenkinsKubernetes.Metadata.Id, jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressGrpcInternalHostname = fmt.Sprintf("%s-grpc-internal.%s", jenkinsKubernetes.Metadata.Id,
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressPlaygroundExternalHostname = fmt.Sprintf("%s.%s",
		jenkinsKubernetes.Metadata.Id, jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressPlaygroundInternalHostname = fmt.Sprintf("%s-internal.%s", jenkinsKubernetes.Metadata.Id,
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressMetricsExternalHostname = fmt.Sprintf("%s-metrics.%s",
		jenkinsKubernetes.Metadata.Id, jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressMetricsInternalHostname = fmt.Sprintf("%s-metrics-internal.%s", jenkinsKubernetes.Metadata.Id,
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressHostnames = []string{
		locals.IngressHttpExternalHostname,
		locals.IngressHttpInternalHostname,
		locals.IngressGrpcExternalHostname,
		locals.IngressGrpcInternalHostname,
		locals.IngressPlaygroundExternalHostname,
		locals.IngressPlaygroundInternalHostname,
		locals.IngressMetricsExternalHostname,
		locals.IngressMetricsInternalHostname,
	}

	//export ingress hostnames
	//ctx.Export(outputs.IngressExternalHostname, pulumi.String(locals.IngressExternalHostname))
	//ctx.Export(outputs.IngressInternalHostname, pulumi.String(locals.IngressInternalHostname))

	//note: a ClusterIssuer resource should have already exist on the kubernetes-cluster.
	//this is typically taken care of by the kubernetes cluster administrator.
	//if the kubernetes-cluster is created using Planton Cloud, then the cluster-issuer name will be
	//same as the ingress-domain-name as long as the same ingress-domain-name is added to the list of
	//ingress-domain-names for the GkeCluster/EksCluster/AksCluster spec.
	locals.IngressCertClusterIssuerName = jenkinsKubernetes.Spec.Ingress.EndpointDomainName

	locals.IngressCertSecretName = jenkinsKubernetes.Metadata.Id

	return locals
}
