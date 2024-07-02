package gateway

import (
	"fmt"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud/environment-pulumi-blueprint/pkg/gcpgke/endpointdomains/cert"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/controller"
	ingressnamespace "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/namespace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	//JenkinsGatewayIdentifierHttps is used as prefix for naming the gateway resource
	JenkinsGatewayIdentifierHttps = "jenkins-https"
	JenkinsGatewayIdentifierHttp  = "jenkins-http"
)

func Resources(ctx *pulumi.Context) error {
	i := extractInput(ctx)

	gatewayObject := buildGatewayObject(i)
	resourceName := fmt.Sprintf("gateway-%s", gatewayObject.Name)
	manifestPath := filepath.Join(i.workspaceDir, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, gatewayObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}

	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName,
		&pulumik8syaml.ConfigFileArgs{File: manifestPath}, pulumi.Provider(i.kubeProvider))
	if err != nil {
		return errors.Wrap(err, "failed to add gateway manifest")
	}
	return nil
}

func buildGatewayObject(i *input) *v1beta1.Gateway {
	return &v1beta1.Gateway{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "Gateway",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      i.resourceId,
			Namespace: ingressnamespace.Name,
			Labels:    i.labels,
		},
		Spec: networkingv1beta1.Gateway{
			Selector: controller.SelectorLabels,
			Servers: []*networkingv1beta1.Server{
				{
					Name: JenkinsGatewayIdentifierHttps,
					Port: &networkingv1beta1.Port{
						Number:   443,
						Protocol: "HTTPS",
						Name:     JenkinsGatewayIdentifierHttps,
					},
					Hosts: []string{i.externalHostname},
					Tls: &networkingv1beta1.ServerTLSSettings{
						Mode:           networkingv1beta1.ServerTLSSettings_SIMPLE,
						CredentialName: cert.GetCertSecretName(i.envDomainName),
					},
				},
				{
					Name: JenkinsGatewayIdentifierHttp,
					Port: &networkingv1beta1.Port{
						Number:   80,
						Protocol: "HTTP",
						Name:     JenkinsGatewayIdentifierHttp,
					},
					Hosts: []string{i.externalHostname},
					Tls: &networkingv1beta1.ServerTLSSettings{
						HttpsRedirect: true,
					},
				},
			},
		},
	}
}
