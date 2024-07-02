package secret

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	err := addSecret(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add secret")
	}
	return nil
}

func addSecret(ctx *pulumi.Context) error {
	i := extractInput(ctx)

	// Encode the password in Base64
	base64Password := i.randomPassword.Result.ApplyT(func(p string) (string, error) {
		return base64.StdEncoding.EncodeToString([]byte(p)), nil
	}).(pulumi.StringOutput)

	// Create or update the secret
	_, err := kubernetescorev1.NewSecret(ctx, i.resourceName, &kubernetescorev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(i.resourceName),
			Namespace: pulumi.String(i.namespaceName),
		},
		Data: pulumi.StringMap{
			JenkinsAdminPasswordKey: base64Password,
		},
	}, pulumi.Provider(i.kubeProvider), pulumi.Parent(i.namespace), pulumi.Parent(i.randomPassword))

	if err != nil {
		return fmt.Errorf("failed to create or update kubernetes secret: %w", err)
	}

	return nil
}
