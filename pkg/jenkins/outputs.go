package gcp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/outputs"
	jenkinskubernetesmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/jenkinskubernetes/model"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	jenkinskubernetesstackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/jenkinskubernetes/stack/model"
)

func Outputs(ctx context.Context, input *jenkinskubernetesstackmodel.JenkinsKubernetesStackInput) (*jenkinskubernetesmodel.JenkinsKubernetesStatusStackOutputs, error) {
	pulumiOrgName, err := org.GetOrgName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi org name")
	}
	stackOutput, err := backend.StackOutput(pulumiOrgName, input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return OutputMapTransformer(stackOutput, input), nil
}

func OutputMapTransformer(stackOutput map[string]interface{}, input *jenkinskubernetesstackmodel.JenkinsKubernetesStackInput) *jenkinskubernetesmodel.JenkinsKubernetesStatusStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &jenkinskubernetesmodel.JenkinsKubernetesStatusStackOutputs{}
	}
	return &jenkinskubernetesmodel.JenkinsKubernetesStatusStackOutputs{
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
