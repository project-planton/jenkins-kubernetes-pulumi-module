package gateway

import (
	jenkinskubernetescontextconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	workspaceDir     string
	kubeProvider     *pulumikubernetes.Provider
	resourceId       string
	labels           map[string]string
	envDomainName    string
	externalHostname string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(jenkinskubernetescontextconfig.Key).(jenkinskubernetescontextconfig.ContextConfig)

	return &input{
		workspaceDir:     ctxConfig.Spec.WorkspaceDir,
		kubeProvider:     ctxConfig.Spec.KubeProvider,
		resourceId:       ctxConfig.Spec.ResourceId,
		labels:           ctxConfig.Spec.Labels,
		envDomainName:    ctxConfig.Spec.EnvDomainName,
		externalHostname: ctxConfig.Spec.ExternalHostname,
	}
}
