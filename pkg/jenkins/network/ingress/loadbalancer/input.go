package loadbalancer

import (
	jenkinskubernetescontextconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	ExternalLoadBalancerServiceName             = "ingress-external-lb"
	InternalLoadBalancerServiceName             = "ingress-internal-lb"
	ExternalLoadBalancerExternalNameServiceName = "ingress-external-external-dns"
	InternalLoadBalancerExternalNameServiceName = "ingress-internal-external-dns"
)

type input struct {
	resourceId         string
	resourceName       string
	namespace          *kubernetescorev1.Namespace
	externalHostname   string
	internalHostname   string
	endpointDomainName string
	serviceName        string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(jenkinskubernetescontextconfig.Key).(jenkinskubernetescontextconfig.ContextConfig)

	return &input{
		resourceId:         ctxConfig.Spec.ResourceId,
		resourceName:       ctxConfig.Spec.ResourceName,
		namespace:          ctxConfig.Status.AddedResources.Namespace,
		externalHostname:   ctxConfig.Spec.ExternalHostname,
		internalHostname:   ctxConfig.Spec.InternalHostname,
		endpointDomainName: ctxConfig.Spec.EndpointDomainName,
		serviceName:        ctxConfig.Spec.KubeServiceName,
	}
}
