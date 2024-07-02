package helmchart

import (
	"github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/secret"
	"github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/containerresources"
	"github.com/plantoncloud/pulumi-blueprint-commons/pkg/kubernetes/helm/mergemaps"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func getHelmValues(i *input) pulumi.Map {
	// https://github.com/jenkinsci/helm-charts/blob/main/charts/jenkins/values.yaml
	var baseValues = pulumi.Map{
		"fullnameOverride": pulumi.String(i.kubeServiceName),
		"controller": pulumi.Map{
			"image": pulumi.Map{
				"tag": pulumi.String("2.454-jdk17"),
			},
			"resources": containerresources.ConvertToPulumiMap(i.containerSpec.Resources),
			"admin": pulumi.Map{
				"passwordKey":    pulumi.String(secret.JenkinsAdminPasswordKey),
				"existingSecret": pulumi.String(i.kubeServiceName),
			},
		},
	}
	mergemaps.MergeMapToPulumiMap(baseValues, i.customHelmValues)
	return baseValues
}
