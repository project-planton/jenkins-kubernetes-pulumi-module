package istio

import (
	"github.com/pkg/errors"
	jenkinskubernetesnetworkgateway "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/network/ingress/istio/gateway"
	jenkinskubernetesnetworkvirtualservice "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/network/ingress/istio/virtualservice"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	if err := jenkinskubernetesnetworkgateway.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add gateway resources")
	}
	if err := jenkinskubernetesnetworkvirtualservice.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add gateway resources")
	}
	return nil
}
