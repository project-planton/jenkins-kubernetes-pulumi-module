package jenkins

import (
	"encoding/base64"
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	AdminPasswordKey = "jenkins-admin-password"
)

func (s *ResourceStack) adminPassword(ctx *pulumi.Context,
	addedNamespace *kubernetescorev1.Namespace) (*kubernetescorev1.Secret, error) {

	jenkinsKubernetes := s.Input.ApiResource

	addedRandomPassword, err := random.NewRandomPassword(ctx, "admin-password",
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
	base64Password := addedRandomPassword.Result.ApplyT(func(p string) (string, error) {
		return base64.StdEncoding.EncodeToString([]byte(p)), nil
	}).(pulumi.StringOutput)

	// Create or update the secret
	addedAdminPasswordSecret, err := kubernetescorev1.NewSecret(ctx, jenkinsKubernetes.Metadata.Name,
		&kubernetescorev1.SecretArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String(jenkinsKubernetes.Metadata.Name),
				Namespace: pulumi.String(s.namespace()),
			},
			Data: pulumi.StringMap{
				AdminPasswordKey: base64Password,
			},
		}, pulumi.Parent(addedNamespace))

	return addedAdminPasswordSecret, nil
}
