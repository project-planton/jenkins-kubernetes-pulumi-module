package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *model.JenkinsKubernetesStackInput
	Labels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesClusterCredential(ctx,
		s.Input.KubernetesClusterCredential)
	if err != nil {
		return errors.Wrap(err, "failed to setup gcp provider")
	}

	//create a new descriptive variable for the api-resource in the input.
	jenkinsKubernetes := s.Input.ApiResource

	//decide on the name of the namespace
	namespaceName := jenkinsKubernetes.Metadata.Id

	//create namespace resource
	createdNamespace, err := kubernetescorev1.NewNamespace(ctx, namespaceName, &kubernetescorev1.NamespaceArgs{
		Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
			Name:   pulumi.String(namespaceName),
			Labels: pulumi.ToStringMap(s.Labels),
		}),
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5s", Update: "5s", Delete: "5s"}),
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s namespace", namespaceName)
	}

	//export name of the namespace
	ctx.Export(outputs.Namespace, createdNamespace.Metadata.Name())

	//create admin-password secret
	createdAdminPasswordSecret, err := s.adminPassword(ctx, createdNamespace)
	if err != nil {
		return errors.Wrap(err, "failed to create admin password resources")
	}

	//install the jenkins helm-chart
	if err := s.helmChart(ctx, createdNamespace, createdAdminPasswordSecret); err != nil {
		return errors.Wrap(err, "failed to create helm-chart resources")
	}

	//export kubernetes service name
	ctx.Export(outputs.Service, pulumi.String(jenkinsKubernetes.Metadata.Name))

	//export kubernetes endpoint
	ctx.Export(outputs.KubeEndpoint,
		pulumi.Sprintf("%s.%s.svc.cluster.local.",
			jenkinsKubernetes.Metadata.Name,
			namespaceName))

	//export kube-port-forward command
	ctx.Export(outputs.PortForwardCommand, pulumi.Sprintf(
		"kubectl port-forward -n %s service/%s 8080:8080",
		namespaceName, jenkinsKubernetes.Metadata.Name))

	//no ingress resources required when ingress is not enabled
	if !jenkinsKubernetes.Spec.Ingress.IsEnabled {
		return nil
	}

	if err := s.istioIngress(ctx, createdNamespace); err != nil {
		return errors.Wrap(err, "failed to create istio ingress resources")
	}

	//export ingress hostnames
	ctx.Export(outputs.IngressExternalHostname, pulumi.Sprintf("%s.%s",
		jenkinsKubernetes.Metadata.Id, jenkinsKubernetes.Spec.Ingress.EndpointDomainName))
	ctx.Export(outputs.IngressInternalHostname, pulumi.Sprintf("%s-internal.%s",
		jenkinsKubernetes.Metadata.Id, jenkinsKubernetes.Spec.Ingress.EndpointDomainName))

	return nil
}
