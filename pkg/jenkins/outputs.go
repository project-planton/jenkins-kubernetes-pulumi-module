package jenkins

import (
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-module/pkg/jenkins/outputs"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
)

func OutputMapTransformer(stackOutput map[string]interface{},
	input *model.JenkinsKubernetesStackInput) *model.JenkinsKubernetesStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply ||
		stackOutput == nil {
		return &model.JenkinsKubernetesStackOutputs{}
	}
	return &model.JenkinsKubernetesStackOutputs{
		Namespace:               backend.GetVal(stackOutput, outputs.GetNamespaceNameOutputName()),
		AdminUsername:           backend.GetVal(stackOutput, outputs.GetAdminUsernameOutputName()),
		AdminPasswordSecretName: backend.GetVal(stackOutput, outputs.GetAdminPasswordSecretOutputName()),
		Service:                 backend.GetVal(stackOutput, outputs.GetKubeServiceNameOutputName()),
		PortForwardCommand:      backend.GetVal(stackOutput, outputs.GetKubePortForwardCommandOutputName()),
		KubeEndpoint:            backend.GetVal(stackOutput, outputs.GetKubeEndpointOutputName()),
		ExternalHostname:        backend.GetVal(stackOutput, outputs.GetExternalHostnameOutputName()),
		InternalHostname:        backend.GetVal(stackOutput, outputs.GetInternalHostnameOutputName()),
	}
}
