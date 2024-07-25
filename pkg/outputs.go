package pkg

import (
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-module/pkg/opname"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func PulumiOutputToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *model.JenkinsKubernetesStackInput) *model.JenkinsKubernetesStackOutputs {
	return &model.JenkinsKubernetesStackOutputs{
		Namespace:               autoapistackoutput.GetVal(pulumiOutputs, opname.NamespaceOutputName),
		AdminUsername:           autoapistackoutput.GetVal(pulumiOutputs, opname.AdminUsernameOutputName),
		AdminPasswordSecretName: autoapistackoutput.GetVal(pulumiOutputs, opname.AdminPasswordSecretOutputName),
		Service:                 autoapistackoutput.GetVal(pulumiOutputs, opname.ServiceOutputName),
		PortForwardCommand:      autoapistackoutput.GetVal(pulumiOutputs, opname.PortForwardCommandOutputName),
		KubeEndpoint:            autoapistackoutput.GetVal(pulumiOutputs, opname.KubeEndpointOutputName),
		ExternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, opname.IngressExternalHostnameOutputName),
		InternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, opname.IngressInternalHostnameOutputName),
	}
}
