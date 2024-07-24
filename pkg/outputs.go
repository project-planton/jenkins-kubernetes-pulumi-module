package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	NamespaceOutputName               = "namespace"
	AdminUsernameOutputName           = "admin-username"
	AdminPasswordSecretOutputName     = "admin-password-secret"
	ServiceOutputName                 = "service"
	PortForwardCommandOutputName      = "port-forward-command"
	KubeEndpointOutputName            = "kube-endpoint"
	IngressExternalHostnameOutputName = "ingress-external-hostname"
	IngressInternalHostnameOutputName = "ingress-internal-hostname"
)

func PulumiOutputToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *model.JenkinsKubernetesStackInput) *model.JenkinsKubernetesStackOutputs {
	return &model.JenkinsKubernetesStackOutputs{
		Namespace:               autoapistackoutput.GetVal(pulumiOutputs, NamespaceOutputName),
		AdminUsername:           autoapistackoutput.GetVal(pulumiOutputs, AdminUsernameOutputName),
		AdminPasswordSecretName: autoapistackoutput.GetVal(pulumiOutputs, AdminPasswordSecretOutputName),
		Service:                 autoapistackoutput.GetVal(pulumiOutputs, ServiceOutputName),
		PortForwardCommand:      autoapistackoutput.GetVal(pulumiOutputs, PortForwardCommandOutputName),
		KubeEndpoint:            autoapistackoutput.GetVal(pulumiOutputs, KubeEndpointOutputName),
		ExternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, IngressExternalHostnameOutputName),
		InternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, IngressInternalHostnameOutputName),
	}
}
