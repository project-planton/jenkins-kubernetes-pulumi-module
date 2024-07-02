package network

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/network/ingress"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (*pulumi.Context, error) {
	i := extractInput(ctx)

	if !i.isIngressEnabled || i.endpointDomainName == "" {
		return ctx, nil
	}
	ctx, err := ingress.Resources(ctx)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to add gateway resources")
	}
	return ctx, nil
}
