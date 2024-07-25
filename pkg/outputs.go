package pkg

import (
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-module/pkg/variables"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func PulumiOutputToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *model.JenkinsKubernetesStackInput) *model.JenkinsKubernetesStackOutputs {
	return &model.JenkinsKubernetesStackOutputs{
		Namespace:               autoapistackoutput.GetVal(pulumiOutputs, variables.NamespaceOutputName),
		AdminUsername:           autoapistackoutput.GetVal(pulumiOutputs, variables.AdminUsernameOutputName),
		AdminPasswordSecretName: autoapistackoutput.GetVal(pulumiOutputs, variables.AdminPasswordSecretOutputName),
		Service:                 autoapistackoutput.GetVal(pulumiOutputs, variables.ServiceOutputName),
		PortForwardCommand:      autoapistackoutput.GetVal(pulumiOutputs, variables.PortForwardCommandOutputName),
		KubeEndpoint:            autoapistackoutput.GetVal(pulumiOutputs, variables.KubeEndpointOutputName),
		ExternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, variables.IngressExternalHostnameOutputName),
		InternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, variables.IngressInternalHostnameOutputName),
	}
}
