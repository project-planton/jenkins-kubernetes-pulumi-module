package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/kubernetes/enums/kubernetesworkloadingresstype"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/kubernetes/pulumikubernetesprovider"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/pulumi/pulumicustomoutput"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	WorkspaceDir     string
	Input            *model.JenkinsKubernetesStackInput
	KubernetesLabels map[string]string
}

const (
	JenkinsPort = 8080
)

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

	ctx.Export(GetNamespaceNameOutputName(), addedNamespace.Metadata.Name())

	addedAdminPasswordSecret, err := s.adminPassword(ctx, addedNamespace)
	if err != nil {
		return errors.Wrap(err, "failed to add admin password resources")
	}

	if err := s.helmChart(ctx, addedNamespace, addedAdminPasswordSecret); err != nil {
		return errors.Wrap(err, "failed to add helm-chart resources")
	}

	ctx.Export(GetKubePortForwardCommandOutputName(), pulumi.String(getKubePortForwardCommand()))

	ctx.Export(GetExternalHostnameOutputName(), pulumi.String(i.externalHostname))
	ctx.Export(GetInternalHostnameOutputName(), pulumi.String(i.internalHostname))

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

func GetKubeServiceNameFqdn(jenkinsKubernetesName, namespace string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", jenkinsKubernetesName, namespace)
}

func GetNamespaceNameOutputName() string {
	return pulumikubernetesprovider.PulumiOutputName(kubernetescorev1.Namespace{},
		englishword.EnglishWord_namespace.String())
}

func GetKubePortForwardCommandOutputName() string {
	return pulumicustomoutput.Name("jenkins-cluster-kube-port-forward-command")
}

func GetExternalHostnameOutputName() string {
	return pulumicustomoutput.Name("jenkins-kubernetes-external-hostname")
}

func GetInternalHostnameOutputName() string {
	return pulumicustomoutput.Name("jenkins-kubernetes-internal-hostname")
}

func GetKubeServiceNameOutputName() string {
	return pulumicustomoutput.Name("jenkins-kubernetes-kubernetes-service-name")
}

func GetKubeEndpointOutputName() string {
	return pulumicustomoutput.Name("jenkins-kubernetes-kubernetes-endpoint")
}
