package helmchart

import (
	jenkinskubernetescontextconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	jenkinskubernetesmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/jenkinskubernetes/model"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	resourceId       string
	namespace        *kubernetescorev1.Namespace
	kubeServiceName  string
	containerSpec    *jenkinskubernetesmodel.JenkinsKubernetesSpecContainerSpec
	customHelmValues map[string]string
	labels           map[string]string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(jenkinskubernetescontextconfig.Key).(jenkinskubernetescontextconfig.ContextConfig)

	return &input{
		resourceId:       ctxConfig.Spec.ResourceId,
		namespace:        ctxConfig.Status.AddedResources.Namespace,
		containerSpec:    ctxConfig.Spec.ContainerSpec,
		customHelmValues: ctxConfig.Spec.CustomHelmValues,
		labels:           ctxConfig.Spec.Labels,
		kubeServiceName:  ctxConfig.Spec.KubeServiceName,
	}
}
