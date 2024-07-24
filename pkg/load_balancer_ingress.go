package pkg

import (
	"github.com/pkg/errors"
	pulumicommonsloadbalancerservice "github.com/plantoncloud/pulumi-module-golang-commons/pkg/kubernetes/loadbalancer/service"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/pulumi/pulumicustomoutput"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	ExternalLoadBalancerServiceName = "ingress-external-lb"
	InternalLoadBalancerServiceName = "ingress-internal-lb"
)

func (s *ResourceStack) loadBalancerIngress(ctx *pulumi.Context, createdNamespace *kubernetescorev1.Namespace) error {
	createdExternalService, err := kubernetescorev1.NewService(ctx,
		ExternalLoadBalancerServiceName,
		&kubernetescorev1.ServiceArgs{
			Metadata: &kubernetesmetav1.ObjectMetaArgs{
				Name:      pulumi.String(ExternalLoadBalancerServiceName),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    createdNamespace.Metadata.Labels(),
				Annotations: pulumi.StringMap{
					"planton.cloud/endpoint-domain-name":        pulumi.String(i.endpointDomainName),
					"external-dns.alpha.kubernetes.io/hostname": pulumi.String(hostname)}},
			Spec: &kubernetescorev1.ServiceSpecArgs{
				Type: pulumi.String("LoadBalancer"), // Service type is LoadBalancer
				Ports: kubernetescorev1.ServicePortArray{
					&kubernetescorev1.ServicePortArgs{
						Name:       pulumi.String("http"),
						Port:       pulumi.Int(80),
						Protocol:   pulumi.String("TCP"),
						TargetPort: pulumi.String("http"), // This assumes your Jenkins pod has a port named 'http'
					},
				},
				Selector: pulumi.StringMap{
					"app.kubernetes.io/component": pulumi.String("jenkins-controller"),
					"app.kubernetes.io/instance":  createdNamespace.Metadata.Name().Elem(),
				},
			},
		}, pulumi.Parent(createdNamespace))
	if err != nil {
		return errors.Wrapf(err, "failed to create external load balancer service")
	}

	createdInternalService, err := kubernetescorev1.NewService(ctx,
		InternalLoadBalancerServiceName,
		&kubernetescorev1.ServiceArgs{
			Metadata: &kubernetesmetav1.ObjectMetaArgs{
				Name:      pulumi.String(InternalLoadBalancerServiceName),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    createdNamespace.Metadata.Labels(),
				Annotations: pulumi.StringMap{
					"cloud.google.com/load-balancer-type":       pulumi.String("Internal"),
					"planton.cloud/endpoint-domain-name":        pulumi.String(i.endpointDomainName),
					"external-dns.alpha.kubernetes.io/hostname": pulumi.String(hostname),
				},
			},
			Spec: &kubernetescorev1.ServiceSpecArgs{
				Type: pulumi.String("LoadBalancer"), // Service type is LoadBalancer
				Ports: kubernetescorev1.ServicePortArray{
					&kubernetescorev1.ServicePortArgs{
						Name:       pulumi.String("http"),
						Port:       pulumi.Int(80),
						Protocol:   pulumi.String("TCP"),
						TargetPort: pulumi.String("http"), // This assumes your Jenkins pod has a port named 'http'
					},
				},
				Selector: pulumi.StringMap{
					"app.kubernetes.io/component": pulumi.String("jenkins-controller"),
					"app.kubernetes.io/instance":  createdNamespace.Metadata.Name().Elem(),
				},
			},
		}, pulumi.Parent(createdNamespace))
	if err != nil {
		return errors.Wrapf(err, "failed to create external load balancer service")
	}

	ctx.Export(GetExternalLoadBalancerIp(),
		pulumi.String(pulumicommonsloadbalancerservice.GetIpAddress(createdExternalService)))
	ctx.Export(GetInternalLoadBalancerIp(),
		pulumi.String(pulumicommonsloadbalancerservice.GetIpAddress(createdInternalService)))
	return nil
}

func GetExternalLoadBalancerIp() string {
	return pulumicustomoutput.Name("jenkins-ingress-external-lb-ip")
}

func GetInternalLoadBalancerIp() string {
	return pulumicustomoutput.Name("jenkins-ingress-internal-lb-ip")
}
