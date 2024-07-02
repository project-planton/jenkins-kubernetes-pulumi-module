package outputs

import (
	jenkinskubernetescontextconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	pulumicommonsloadbalancerservice "github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/loadbalancer/service"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	adminUsername = "admin"
	JenkinsPort   = 8080
)

type input struct {
	resourceId                    string
	resourceName                  string
	environmentName               string
	endpointDomainName            string
	namespaceName                 string
	externalLoadBalancerIpAddress string
	internalLoadBalancerIpAddress string
	internalHostname              string
	externalHostname              string
	kubeServiceName               string
	kubeLocalEndpoint             string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(jenkinskubernetescontextconfig.Key).(jenkinskubernetescontextconfig.ContextConfig)
	var externalLoadBalancerIpAddress = ""
	var internalLoadBalancerIpAddress = ""

	if ctxConfig.Status.AddedResources.LoadBalancerExternalService != nil {
		externalLoadBalancerIpAddress = pulumicommonsloadbalancerservice.GetIpAddress(ctxConfig.Status.AddedResources.LoadBalancerExternalService)
	}

	if ctxConfig.Status.AddedResources.LoadBalancerInternalService != nil {
		internalLoadBalancerIpAddress = pulumicommonsloadbalancerservice.GetIpAddress(ctxConfig.Status.AddedResources.LoadBalancerExternalService)
	}

	return &input{
		resourceId:                    ctxConfig.Spec.ResourceId,
		resourceName:                  ctxConfig.Spec.ResourceName,
		environmentName:               ctxConfig.Spec.EnvironmentInfo.EnvironmentName,
		endpointDomainName:            ctxConfig.Spec.EndpointDomainName,
		namespaceName:                 ctxConfig.Spec.NamespaceName,
		externalLoadBalancerIpAddress: externalLoadBalancerIpAddress,
		internalLoadBalancerIpAddress: internalLoadBalancerIpAddress,
		internalHostname:              ctxConfig.Spec.InternalHostname,
		externalHostname:              ctxConfig.Spec.ExternalHostname,
		kubeServiceName:               ctxConfig.Spec.KubeServiceName,
		kubeLocalEndpoint:             ctxConfig.Spec.KubeLocalEndpoint,
	}
}
