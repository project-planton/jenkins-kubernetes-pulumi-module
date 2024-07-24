package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	NamespaceOutputName           = "namespace"
	AdminUsernameOutputName       = "admin-username"
	AdminPasswordSecretOutputName = "admin-password-secret"
	ServiceOutputName             = "service"
	PortForwardCommandOutputName  = "port-forward-command"
	KubeEndpointOutputName        = "kube-endpoint"
	IngressExternalHostname       = "ingress-external-hostname"
	IngressInternalHostname       = "ingress-internal-hostname"
)

func StackOutputsBuilder(stackOutput auto.OutputMap,
	input *model.JenkinsKubernetesStackInput) *model.JenkinsKubernetesStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply ||
		stackOutput == nil {
		return &model.JenkinsKubernetesStackOutputs{}
	}
	return &model.JenkinsKubernetesStackOutputs{
		Namespace:               aut.GetVal(stackOutput, NamespaceOutputName),
		AdminUsername:           backend.GetVal(stackOutput, AdminUsernameOutputName),
		AdminPasswordSecretName: backend.GetVal(stackOutput, AdminPasswordSecretOutputName),
		Service:                 backend.GetVal(stackOutput, ServiceOutputName),
		PortForwardCommand:      backend.GetVal(stackOutput, PortForwardCommandOutputName),
		KubeEndpoint:            backend.GetVal(stackOutput, KubeEndpointOutputName),
		ExternalHostname:        backend.GetVal(stackOutput, IngressExternalHostname),
		InternalHostname:        backend.GetVal(stackOutput, IngressInternalHostname),
	}
}
