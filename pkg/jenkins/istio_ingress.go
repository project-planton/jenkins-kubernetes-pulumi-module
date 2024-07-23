package jenkins

import (
	"github.com/pkg/errors"
	istiov1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/istio/kubernetes/networking/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	v1 "istio.io/api/networking/v1"
)

const (
	IstioIngressNamespace = "istio-ingress"
	//GatewayIdentifierHttps is used as prefix for naming the gateway resource
	GatewayIdentifierHttps = "jenkins-https"
	GatewayIdentifierHttp  = "jenkins-http"
	Port                   = 8080
)

func (s *ResourceStack) istioIngress(ctx *pulumi.Context, addedNamespace *kubernetescorev1.Namespace) error {
	jenkinsKubernetes := s.Input.ApiResource

	_, err := istiov1.NewGateway(ctx, jenkinsKubernetes.Metadata.Id, &istiov1.GatewayArgs{
		ApiVersion: pulumi.String("networking.istio.io/v1"),
		Kind:       pulumi.String("Gateway"),
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(jenkinsKubernetes.Metadata.Id),
			Namespace: addedNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(s.KubernetesLabels),
		},
		Spec: istiov1.GatewaySpecArgs{
			Selector: nil, //will need to add map to pulumi-module-commons repo and get it from there
			Servers: istiov1.GatewaySpecServersArray{
				&istiov1.GatewaySpecServersArgs{
					Name: pulumi.String(GatewayIdentifierHttps),
					Port: &istiov1.GatewaySpecServersPortArgs{
						Number:   pulumi.Int(443),
						Name:     pulumi.String(GatewayIdentifierHttps),
						Protocol: pulumi.String("HTTPS"),
					},
					Hosts: pulumi.StringArray{
						pulumi.String("coming-soon"), //todo: add external-hostname
					},
					Tls: &istiov1.GatewaySpecServersTlsArgs{
						CredentialName: pulumi.String("coming-soon"), //todo: create a certificate for jenkins-server
						Mode:           pulumi.String(v1.ServerTLSSettings_SIMPLE.String()),
					},
				},
				&istiov1.GatewaySpecServersArgs{
					Name: pulumi.String(GatewayIdentifierHttp),
					Port: &istiov1.GatewaySpecServersPortArgs{
						Number:   pulumi.Int(80),
						Name:     pulumi.String(GatewayIdentifierHttp),
						Protocol: pulumi.String("HTTP"),
					},
					Hosts: pulumi.StringArray{
						pulumi.String("coming-soon"), //todo: add external-hostname
					},
					Tls: &istiov1.GatewaySpecServersTlsArgs{
						HttpsRedirect: pulumi.Bool(true),
					},
				},
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "error creating gateway")
	}

	_, err = istiov1.NewVirtualService(ctx, jenkinsKubernetes.Metadata.Id,
		&istiov1.VirtualServiceArgs{
			ApiVersion: pulumi.String("networking.istio.io/v1"),
			Kind:       pulumi.String("VirtualService"),
			Metadata: metav1.ObjectMetaArgs{
				Name:      pulumi.String(jenkinsKubernetes.Metadata.Id),
				Namespace: addedNamespace.Metadata.Name(),
				Labels:    pulumi.ToStringMap(s.KubernetesLabels),
			},
			Spec: istiov1.VirtualServiceSpecArgs{
				Gateways: pulumi.StringArray{
					pulumi.Sprintf("%s/%s", IstioIngressNamespace,
						jenkinsKubernetes.Metadata.Id),
				},
				Hosts: nil, //todo: add hostnames
				Http: istiov1.VirtualServiceSpecHttpArray{
					&istiov1.VirtualServiceSpecHttpArgs{
						Name: pulumi.String(jenkinsKubernetes.Metadata.Id),
						Route: istiov1.VirtualServiceSpecHttpRouteArray{
							&istiov1.VirtualServiceSpecHttpRouteArgs{
								Destination: istiov1.VirtualServiceSpecHttpRouteDestinationArgs{
									Host: nil, //todo: add kube-endpoint
									Port: istiov1.VirtualServiceSpecHttpRouteDestinationPortArgs{
										Number: pulumi.Int(Port),
									},
								},
							},
						},
					},
				},
			},
			Status: nil,
		})
	if err != nil {
		return errors.Wrap(err, "error creating virtual-service")
	}
	return nil
}
