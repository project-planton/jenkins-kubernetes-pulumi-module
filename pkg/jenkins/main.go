package jenkins

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/network/dns"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/kubernetes/enums/kubernetesworkloadingresstype"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	WorkspaceDir     string
	Input            *model.JenkinsKubernetesStackInput
	KubernetesLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesClusterCredential(ctx,
		s.Input.KubernetesClusterCredential)
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	namespaceName := s.namespace()

	jenkinsKubernetes := s.Input.ApiResource

	addedNamespace, err := kubernetescorev1.NewNamespace(ctx, namespaceName, &kubernetescorev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("namespace"),
		Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
			Name:   pulumi.String(namespaceName),
			Labels: pulumi.ToStringMap(s.KubernetesLabels),
		}),
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5s", Update: "5s", Delete: "5s"}),
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s namespace", namespaceName)
	}

	addedAdminPasswordSecret, err := s.adminPassword(ctx, addedNamespace)
	if err != nil {
		return errors.Wrap(err, "failed to add admin password resources")
	}

	if err := s.helmChart(ctx, addedNamespace, addedAdminPasswordSecret); err != nil {
		return errors.Wrap(err, "failed to add helm-chart resources")
	}

	//no ingress resources required when ingress is not enabled
	if !jenkinsKubernetes.Spec.Ingress.IsEnabled || jenkinsKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return nil
	}

	switch jenkinsKubernetes.Spec.Ingress.IngressType {
	case kubernetesworkloadingresstype.KubernetesWorkloadIngressType_load_balancer:
		if err := s.loadBalancerIngress(ctx, addedNamespace); err != nil {
			return errors.Wrap(err, "failed to add load-balancer ingress resources")
		}
	case kubernetesworkloadingresstype.KubernetesWorkloadIngressType_ingress_controller:
		if err := s.istioIngress(ctx, addedNamespace); err != nil {
			return errors.Wrap(err, "failed to add istio ingress resources")
		}
	}

	return nil
}

func (s *ResourceStack) namespace() string {
	return s.Input.ApiResource.Metadata.Id
}

func GetInternalHostname(jenkinsKubernetesId, environmentName, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s-internal.%s", jenkinsKubernetesId, environmentName, endpointDomainName)
}

func GetExternalHostname(jenkinsKubernetesId, environmentName, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s.%s", jenkinsKubernetesId, environmentName, endpointDomainName)
}

func GetKubeServiceNameFqdn(jenkinsKubernetesName, namespace string) string {
	return fmt.Sprintf("%s.%s.%s", GetKubeServiceName(jenkinsKubernetesName), namespace, dns.DefaultDomain)
}

func GetKubeServiceName(jenkinsKubernetesName string) string {
	return fmt.Sprintf(jenkinsKubernetesName)
}
