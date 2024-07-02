package helmchart

import (
	"github.com/pkg/errors"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	err := addHelmChart(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add helm chart")
	}
	return nil
}

func addHelmChart(ctx *pulumi.Context) error {
	i := extractInput(ctx)

	helmValues := getHelmValues(i)
	// Deploying a Locust Helm chart from the Helm repository.
	_, err := helmv3.NewChart(ctx, i.resourceId, helmv3.ChartArgs{
		Chart:     pulumi.String("jenkins"),
		Version:   pulumi.String("5.1.5"), // Use the Helm chart version you want to install
		Namespace: i.namespace.Metadata.Name().Elem(),
		Values:    helmValues,
		//if you need to add the repository, you can specify `repo url`:
		FetchArgs: helmv3.FetchArgs{
			Repo: pulumi.String("https://charts.jenkins.io"), // The URL for the Helm chart repository
		},
	}, pulumi.Parent(i.namespace))
	if err != nil {
		return errors.Wrap(err, "failed to add helm chart")
	}
	return nil
}
