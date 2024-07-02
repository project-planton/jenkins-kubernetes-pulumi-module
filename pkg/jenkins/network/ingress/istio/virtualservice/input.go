package virtualservice

import (
	jenkinskubernetescontextconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	resourceId      string
	namespace       *kubernetescorev1.Namespace
	hostNames       []string
	workspaceDir    string
	namespaceName   string
	kubeEndpoint    string
	kubeServiceName string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(jenkinskubernetescontextconfig.Key).(jenkinskubernetescontextconfig.ContextConfig)

	return &input{
		resourceId:      ctxConfig.Spec.ResourceId,
		namespace:       ctxConfig.Status.AddedResources.Namespace,
		hostNames:       []string{ctxConfig.Spec.ExternalHostname},
		workspaceDir:    ctxConfig.Spec.WorkspaceDir,
		namespaceName:   ctxConfig.Spec.NamespaceName,
		kubeEndpoint:    ctxConfig.Spec.KubeLocalEndpoint,
		kubeServiceName: ctxConfig.Spec.KubeServiceName,
	}
}
