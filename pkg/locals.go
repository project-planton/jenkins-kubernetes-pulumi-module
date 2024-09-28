package pkg

import (
	"fmt"
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/project-planton/apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes"
	"github.com/plantoncloud/project-planton/apis/zzgo/cloud/planton/apis/commons/apiresource/enums/apiresourcekind"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type Locals struct {
	JenkinsKubernetes            *jenkinskubernetes.JenkinsKubernetes
	Namespace                    string
	IngressCertClusterIssuerName string
	IngressCertSecretName        string
	IngressInternalHostname      string
	IngressExternalHostname      string
	IngressHostnames             []string
	KubeServiceFqdn              string
	KubeServiceName              string
	KubePortForwardCommand       string
	Labels                       map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *jenkinskubernetes.JenkinsKubernetesStackInput) *Locals {
	locals := &Locals{}
	//assign value for the local variable to make it available across the project
	locals.JenkinsKubernetes = stackInput.Target

	jenkinsKubernetes := stackInput.Target

	locals.Labels = map[string]string{
		kuberneteslabelkeys.Environment:  stackInput.Target.Spec.EnvironmentInfo.EnvId,
		kuberneteslabelkeys.Organization: stackInput.Target.Spec.EnvironmentInfo.OrgId,
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceId:   stackInput.Target.Metadata.Id,
		kuberneteslabelkeys.ResourceKind: apiresourcekind.ApiResourceKind_jenkins_kubernetes.String(),
	}

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

	locals.IngressExternalHostname = fmt.Sprintf("%s.%s", jenkinsKubernetes.Metadata.Id,
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressInternalHostname = fmt.Sprintf("%s-internal.%s", jenkinsKubernetes.Metadata.Id,
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressHostnames = []string{
		locals.IngressExternalHostname,
		locals.IngressInternalHostname,
	}

	//export ingress hostnames
	ctx.Export(outputs.IngressExternalHostname, pulumi.String(locals.IngressExternalHostname))
	ctx.Export(outputs.IngressInternalHostname, pulumi.String(locals.IngressInternalHostname))

	//note: a ClusterIssuer resource should have already exist on the kubernetes-cluster.
	//this is typically taken care of by the kubernetes cluster administrator.
	//if the kubernetes-cluster is created using Planton Cloud, then the cluster-issuer name will be
	//same as the ingress-domain-name as long as the same ingress-domain-name is added to the list of
	//ingress-domain-names for the GkeCluster/EksCluster/AksCluster spec.
	locals.IngressCertClusterIssuerName = jenkinsKubernetes.Spec.Ingress.EndpointDomainName

	locals.IngressCertSecretName = jenkinsKubernetes.Metadata.Id

	return locals
}
