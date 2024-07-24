package pkg

import (
	"github.com/pkg/errors"
	istiov1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/istio/networking/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	v1 "istio.io/api/networking/v1"
)

const (
	IstioIngressNamespace = "istio-ingress"
)

func (s *ResourceStack) istioIngress(ctx *pulumi.Context, createdNamespace *kubernetescorev1.Namespace) error {
	jenkinsKubernetes := s.Input.ApiResource

	_, err := istiov1.NewGateway(ctx, jenkinsKubernetes.Metadata.Id, &istiov1.GatewayArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(jenkinsKubernetes.Metadata.Id),
			Namespace: createdNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(s.KubernetesLabels),
		},
		Spec: istiov1.GatewaySpecArgs{
			//the selector labels map should match the desired istio-ingress deployment.
			Selector: pulumi.StringMap{
				"app":   pulumi.String("istio-ingress"),
				"istio": pulumi.String("ingress"),
			},
			Servers: istiov1.GatewaySpecServersArray{
				&istiov1.GatewaySpecServersArgs{
					Name: pulumi.String("jenkins-https"),
					Port: &istiov1.GatewaySpecServersPortArgs{
						Number:   pulumi.Int(443),
						Name:     pulumi.String("jenkins-https"),
						Protocol: pulumi.String("HTTPS"),
					},
					Hosts: pulumi.StringArray{
						pulumi.Sprintf("%s.%s", jenkinsKubernetes.Metadata.Id,
							jenkinsKubernetes.Spec.Ingress.EndpointDomainName),
						pulumi.Sprintf("%s-internal.%s", jenkinsKubernetes.Metadata.Id,
							jenkinsKubernetes.Spec.Ingress.EndpointDomainName),
					},
					Tls: &istiov1.GatewaySpecServersTlsArgs{
						CredentialName: pulumi.String("coming-soon"), //todo: create a certificate for jenkins-server
						Mode:           pulumi.String(v1.ServerTLSSettings_SIMPLE.String()),
					},
				},
				&istiov1.GatewaySpecServersArgs{
					Name: pulumi.String("jenkins-http"),
					Port: &istiov1.GatewaySpecServersPortArgs{
						Number:   pulumi.Int(80),
						Name:     pulumi.String("jenkins-http"),
						Protocol: pulumi.String("HTTP"),
					},
					Hosts: pulumi.StringArray{
						pulumi.Sprintf("%s.%s", jenkinsKubernetes.Metadata.Id,
							jenkinsKubernetes.Spec.Ingress.EndpointDomainName),
						pulumi.Sprintf("%s-internal.%s", jenkinsKubernetes.Metadata.Id,
							jenkinsKubernetes.Spec.Ingress.EndpointDomainName),
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
			Metadata: metav1.ObjectMetaArgs{
				Name:      pulumi.String(jenkinsKubernetes.Metadata.Id),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    pulumi.ToStringMap(s.KubernetesLabels),
			},
			Spec: istiov1.VirtualServiceSpecArgs{
				Gateways: pulumi.StringArray{
					pulumi.Sprintf("%s/%s", IstioIngressNamespace,
						jenkinsKubernetes.Metadata.Id),
				},
				Hosts: pulumi.StringArray{
					pulumi.Sprintf("%s.%s", jenkinsKubernetes.Metadata.Id,
						jenkinsKubernetes.Spec.Ingress.EndpointDomainName),
					pulumi.Sprintf("%s-internal.%s", jenkinsKubernetes.Metadata.Id,
						jenkinsKubernetes.Spec.Ingress.EndpointDomainName),
				},
				Http: istiov1.VirtualServiceSpecHttpArray{
					&istiov1.VirtualServiceSpecHttpArgs{
						Name: pulumi.String(jenkinsKubernetes.Metadata.Id),
						Route: istiov1.VirtualServiceSpecHttpRouteArray{
							&istiov1.VirtualServiceSpecHttpRouteArgs{
								Destination: istiov1.VirtualServiceSpecHttpRouteDestinationArgs{
									Host: pulumi.String(""),
									Port: istiov1.VirtualServiceSpecHttpRouteDestinationPortArgs{
										Number: pulumi.Int(8080),
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
