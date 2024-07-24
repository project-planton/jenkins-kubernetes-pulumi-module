package pkg

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/pulumi/pulumicustomoutput"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	adminPasswordKey = "jenkins-admin-password"
	adminUsername    = "admin"
)

func (s *ResourceStack) adminPassword(ctx *pulumi.Context,
	createdNamespace *kubernetescorev1.Namespace) (*kubernetescorev1.Secret, error) {

	jenkinsKubernetes := s.Input.ApiResource

	createdRandomPassword, err := random.NewRandomPassword(ctx, "admin-password",
		&random.RandomPasswordArgs{
			Length:     pulumi.Int(12),
			Special:    pulumi.Bool(true),
			Numeric:    pulumi.Bool(true),
			Upper:      pulumi.Bool(true),
			Lower:      pulumi.Bool(true),
			MinSpecial: pulumi.Int(3),
			MinNumeric: pulumi.Int(2),
			MinUpper:   pulumi.Int(2),
			MinLower:   pulumi.Int(2),
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get random password value")
	}

	// Encode the password in Base64
	base64Password := createdRandomPassword.Result.ApplyT(func(p string) (string, error) {
		return base64.StdEncoding.EncodeToString([]byte(p)), nil
	}).(pulumi.StringOutput)

	// Create or update the secret
	createdAdminPasswordSecret, err := kubernetescorev1.NewSecret(ctx, jenkinsKubernetes.Metadata.Name,
		&kubernetescorev1.SecretArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(jenkinsKubernetes.Metadata.Name),
				Namespace: pulumi.String(s.namespace()),
			},
			Data: pulumi.StringMap{
				adminPasswordKey: base64Password,
			},
		}, pulumi.Parent(createdNamespace))

	ctx.Export(GetAdminUsernameOutputName(), pulumi.String(adminUsername))
	ctx.Export(GetAdminPasswordSecretOutputName(), createdAdminPasswordSecret.Metadata.Name())

	return createdAdminPasswordSecret, nil
}

func GetAdminUsernameOutputName() string {
	return pulumicustomoutput.Name("jenkins-kubernetes-admin-username")
}

func GetAdminPasswordSecretOutputName() string {
	return pulumicustomoutput.Name("jenkins-kubernetes-admin-password-secret-name")
}