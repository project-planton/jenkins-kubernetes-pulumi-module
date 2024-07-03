package gcp

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/environment-pulumi-blueprint/pkg/gcpgke/endpointdomains/hostnames"
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	jenkinskubernetesnetutilshostname "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/network/ingress/netutils/hostname"
	jenkinskubernetesnetutilsservice "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/network/ingress/netutils/service"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func loadConfig(ctx *pulumi.Context, resourceStack *ResourceStack) (*contextstate.ContextConfig, error) {

	kubernetesProvider, err := pulumikubernetesprovider.GetWithStackCredentials(ctx, resourceStack.Input.CredentialsInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup kubernetes provider")
	}

	var jenkinsKubernetesId = resourceStack.Input.ResourceInput.Metadata.Id
	var jenkinsKubernetesName = resourceStack.Input.ResourceInput.Metadata.Name
	var environmentInfo = resourceStack.Input.ResourceInput.Spec.EnvironmentInfo
	var isIngressEnabled = false

	if resourceStack.Input.ResourceInput.Spec.Ingress != nil {
		isIngressEnabled = resourceStack.Input.ResourceInput.Spec.Ingress.IsEnabled
	}

	var endpointDomainName = ""
	var envDomainName = ""
	var ingressType = kubernetesworkloadingresstype.KubernetesWorkloadIngressType_unspecified
	var internalHostname = ""
	var externalHostname = ""

	if isIngressEnabled {
		endpointDomainName = resourceStack.Input.ResourceInput.Spec.Ingress.EndpointDomainName
		envDomainName = hostnames.GetExternalEnvHostname(environmentInfo.EnvironmentName, endpointDomainName)
		ingressType = resourceStack.Input.ResourceInput.Spec.Ingress.IngressType
		internalHostname = jenkinskubernetesnetutilshostname.GetInternalHostname(jenkinsKubernetesId, environmentInfo.EnvironmentName, endpointDomainName)
		externalHostname = jenkinskubernetesnetutilshostname.GetExternalHostname(jenkinsKubernetesId, environmentInfo.EnvironmentName, endpointDomainName)
	}

	return &contextstate.ContextConfig{
		Spec: &contextstate.Spec{
			KubeProvider:       kubernetesProvider,
			ResourceId:         jenkinsKubernetesId,
			ResourceName:       jenkinsKubernetesName,
			Labels:             resourceStack.KubernetesLabels,
			WorkspaceDir:       resourceStack.WorkspaceDir,
			NamespaceName:      jenkinsKubernetesId,
			EnvironmentInfo:    resourceStack.Input.ResourceInput.Spec.EnvironmentInfo,
			ContainerSpec:      resourceStack.Input.ResourceInput.Spec.Container,
			CustomHelmValues:   resourceStack.Input.ResourceInput.Spec.HelmValues,
			IsIngressEnabled:   isIngressEnabled,
			IngressType:        ingressType,
			EndpointDomainName: endpointDomainName,
			EnvDomainName:      envDomainName,
			InternalHostname:   internalHostname,
			ExternalHostname:   externalHostname,
			KubeServiceName:    jenkinskubernetesnetutilsservice.GetKubeServiceName(jenkinsKubernetesName),
			KubeLocalEndpoint:  jenkinskubernetesnetutilsservice.GetKubeServiceNameFqdn(jenkinsKubernetesName, jenkinsKubernetesId),
		},
	}, nil
}
