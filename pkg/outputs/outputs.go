package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	Namespace               = "namespace"
	AdminUsername           = "admin-username"
	AdminPasswordSecret     = "admin-password-secret"
	Service                 = "service"
	PortForwardCommand      = "port-forward-command"
	KubeEndpoint            = "kube-endpoint"
	IngressExternalHostname = "ingress-external-hostname"
	IngressInternalHostname = "ingress-internal-hostname"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *jenkinskubernetes.JenkinsKubernetesStackInput) *jenkinskubernetes.JenkinsKubernetesStackOutputs {
	return &jenkinskubernetes.JenkinsKubernetesStackOutputs{
		Namespace:               autoapistackoutput.GetVal(pulumiOutputs, Namespace),
		AdminUsername:           autoapistackoutput.GetVal(pulumiOutputs, AdminUsername),
		AdminPasswordSecretName: autoapistackoutput.GetVal(pulumiOutputs, AdminPasswordSecret),
		Service:                 autoapistackoutput.GetVal(pulumiOutputs, Service),
		PortForwardCommand:      autoapistackoutput.GetVal(pulumiOutputs, PortForwardCommand),
		KubeEndpoint:            autoapistackoutput.GetVal(pulumiOutputs, KubeEndpoint),
		ExternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, IngressExternalHostname),
		InternalHostname:        autoapistackoutput.GetVal(pulumiOutputs, IngressInternalHostname),
	}
}
