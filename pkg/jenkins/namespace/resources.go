package namespace

import (
	"github.com/pkg/errors"
	jenkinskubernetescontextconfig "github.com/plantoncloud/jenkins-kubernetes-pulumi-blueprint/pkg/jenkins/contextstate"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (*pulumi.Context, error) {
	namespace, err := addNamespace(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add namespace")
	}

	var ctxConfig = ctx.Value(jenkinskubernetescontextconfig.Key).(jenkinskubernetescontextconfig.ContextConfig)

	addNamespaceToContext(&ctxConfig, namespace)
	ctx = ctx.WithValue(jenkinskubernetescontextconfig.Key, ctxConfig)
	return ctx, nil
}

func addNamespace(ctx *pulumi.Context) (*kubernetescorev1.Namespace, error) {
	var i = extractInput(ctx)

	ns, err := kubernetescorev1.NewNamespace(ctx, i.namespaceName, &kubernetescorev1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("namespace"),
		Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
			Name:   pulumi.String(i.namespaceName),
			Labels: pulumi.ToStringMap(i.labels),
		}),
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5s", Update: "5s", Delete: "5s"}),
		pulumi.Provider(i.kubeProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s namespace", i.namespaceName)
	}
	return ns, nil
}

func addNamespaceToContext(existingConfig *jenkinskubernetescontextconfig.ContextConfig, namespace *kubernetescorev1.Namespace) {
	if existingConfig.Status.AddedResources == nil {
		existingConfig.Status.AddedResources = &jenkinskubernetescontextconfig.AddedResources{
			Namespace: namespace,
		}
		return
	}
	existingConfig.Status.AddedResources.Namespace = namespace
}
