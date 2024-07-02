package network

import (
	jenkinskubernetescontextconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	isIngressEnabled   bool
	endpointDomainName string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(jenkinskubernetescontextconfig.Key).(jenkinskubernetescontextconfig.ContextConfig)

	return &input{
		isIngressEnabled:   ctxConfig.Spec.IsIngressEnabled,
		endpointDomainName: ctxConfig.Spec.EndpointDomainName,
	}
}
