package secret

import (
	jenkinskubernetescontextconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	JenkinsAdminPasswordKey = "jenkins-admin-password"
)

type input struct {
	namespaceName  string
	resourceName   string
	labels         map[string]string
	kubeProvider   *pulumikubernetes.Provider
	namespace      *kubernetescorev1.Namespace
	randomPassword *random.RandomPassword
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(jenkinskubernetescontextconfig.Key).(jenkinskubernetescontextconfig.ContextConfig)

	return &input{
		namespaceName:  ctxConfig.Spec.NamespaceName,
		labels:         ctxConfig.Spec.Labels,
		kubeProvider:   ctxConfig.Spec.KubeProvider,
		resourceName:   ctxConfig.Spec.ResourceName,
		namespace:      ctxConfig.Status.AddedResources.Namespace,
		randomPassword: ctxConfig.Status.AddedResources.RandomPassword,
	}
}
