package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/kubernetes/containerresources"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/kubernetes/helm/mergemaps"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (s *ResourceStack) helmChart(ctx *pulumi.Context,
	createdNamespace *kubernetescorev1.Namespace,
	createdAdminPasswordSecret *kubernetescorev1.Secret) error {

	jenkinsKubernetes := s.Input.ApiResource

	helmValues := getHelmValues(jenkinsKubernetes, createdAdminPasswordSecret)

	// Deploying a Locust Helm chart from the Helm repository.
	_, err := helmv3.NewChart(ctx,
		jenkinsKubernetes.Metadata.Id,
		helmv3.ChartArgs{
			Chart:     pulumi.String("jenkins"),
			Version:   pulumi.String("5.1.5"), // Use the Helm chart version you want to install
			Namespace: createdNamespace.Metadata.Name().Elem(),
			Values:    helmValues,
			//if you need to add the repository, you can specify `repo url`:
			FetchArgs: helmv3.FetchArgs{
				Repo: pulumi.String("https://charts.jenkins.io"),
			},
		}, pulumi.Parent(createdNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to create helm chart")
	}

	return nil
}

func getHelmValues(jenkinsKubernetes *model.JenkinsKubernetes,
	createdAdminPasswordSecret *kubernetescorev1.Secret) pulumi.Map {
	// https://github.com/jenkinsci/helm-charts/blob/main/charts/jenkins/values.yaml
	var baseValues = pulumi.Map{
		"fullnameOverride": pulumi.String(jenkinsKubernetes.Metadata.Name),
		"controller": pulumi.Map{
			"image": pulumi.Map{
				"tag": pulumi.String("2.454-jdk17"),
			},
			"resources": containerresources.ConvertToPulumiMap(jenkinsKubernetes.Spec.Container.Resources),
			"admin": pulumi.Map{
				"passwordKey":    pulumi.String("jenkins-admin-password"),
				"existingSecret": createdAdminPasswordSecret.Metadata.Name(),
			},
		},
	}
	mergemaps.MergeMapToPulumiMap(baseValues, jenkinsKubernetes.Spec.HelmValues)
	return baseValues
}
