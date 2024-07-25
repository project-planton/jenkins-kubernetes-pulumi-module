package pkg

import (
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-module/pkg/vars"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func PulumiOutputToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *model.JenkinsKubernetesStackInput) *model.JenkinsKubernetesStackOutputs {
	return &model.JenkinsKubernetesStackOutputs{
		Namespace:               autoapistackoutput.GetVal(pulumiOutputs, vars.NamespaceOutputName),
		AdminUsername:           autoapistackoutput.GetVal(pulumiOutputs, vars.AdminUsernameOutputName),
		AdminPasswordSecretName: autoapistackoutput.GetVal(pulumiOutputs, vars.AdminPasswordSecretOutputName),
		Service:                 autoapistackoutput.GetVal(pulumiOutputs, vars.ServiceOutputName),
		PortForwardCommand:      autoapistackoutput.GetVal(pulumiOutputs, vars.PortForwardCommandOutputName),
		KubeEndpoint:            autoapistackoutput.GetVal(pulumiOutputs, vars.KubeEndpointOutputName),
		ExternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, vars.IngressExternalHostnameOutputName),
		InternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, vars.IngressInternalHostnameOutputName),
	}
}
