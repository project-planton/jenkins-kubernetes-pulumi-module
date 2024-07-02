package contextstate

import (
	code2cloudenvironmentmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/environment/model"
	jenkinskubernetesmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/jenkinskubernetes/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
)

const (
	Key = "ctx-config"
)

type ContextConfig struct {
	Spec   *Spec
	Status *Status
}

type Spec struct {
	KubeProvider       *kubernetes.Provider
	ResourceId         string
	ResourceName       string
	Labels             map[string]string
	WorkspaceDir       string
	NamespaceName      string
	EnvironmentInfo    *code2cloudenvironmentmodel.ApiResourceEnvironmentInfo
	IsIngressEnabled   bool
	IngressType        kubernetesworkloadingresstype.KubernetesWorkloadIngressType
	ContainerSpec      *jenkinskubernetesmodel.JenkinsKubernetesSpecContainerSpec
	CustomHelmValues   map[string]string
	EndpointDomainName string
	EnvDomainName      string
	InternalHostname   string
	ExternalHostname   string
	KubeServiceName    string
	KubeLocalEndpoint  string
}

type Status struct {
	AddedResources *AddedResources
}

type AddedResources struct {
	Namespace                   *kubernetescorev1.Namespace
	LoadBalancerExternalService *kubernetescorev1.Service
	LoadBalancerInternalService *kubernetescorev1.Service
	RandomPassword              *random.RandomPassword
}
