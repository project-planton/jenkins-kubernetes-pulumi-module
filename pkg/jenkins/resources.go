package gcp

import (
	"github.com/pkg/errors"
	ctxconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	jenkinshelmchart "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/helmchart"
	jenkinsnamespace "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/namespace"
	jenkinsnetwork "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/network"
	jenkinsoutputs "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/outputs"
	jenkinsadminpasswordgenerate "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/password"
	jenkinsadminpasswordsecret "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/secret"
	jenkinskubernetesstackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/jenkinskubernetes/stack/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	WorkspaceDir     string
	Input            *jenkinskubernetesstackmodel.JenkinsKubernetesStackInput
	KubernetesLabels map[string]string
}

func (resourceStack *ResourceStack) Resources(ctx *pulumi.Context) error {

	var ctxConfig, err = loadConfig(ctx, resourceStack)
	if err != nil {
		return errors.Wrap(err, "failed to initiate context context state")
	}
	ctx = ctx.WithValue(ctxconfig.Key, *ctxConfig)

	// Create the namespace resource
	ctx, err = jenkinsnamespace.Resources(ctx)
	if err != nil {
		return err
	}

	// Create the random password resource
	ctx, err = jenkinsadminpasswordgenerate.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create random password resource")
	}

	// Create the secret resource for mongo db root password
	err = jenkinsadminpasswordsecret.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create secret resource")
	}

	// Deploying a Locust Helm chart from the Helm repository.
	err = jenkinshelmchart.Resources(ctx)
	if err != nil {
		return err
	}

	ctx, err = jenkinsnetwork.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add jenkins-kubernetes network resources")
	}

	err = jenkinsoutputs.Export(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to export mongodb cluster outputs")
	}
	return nil
}
