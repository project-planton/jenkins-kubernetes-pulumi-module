package outputs

import (
	"fmt"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/name/output/custom"
	puluminamekubeoutput "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/name/provider/kubernetes/output"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Export(ctx *pulumi.Context) error {
	var i = extractInput(ctx)
	var kubePortForwardCommand = getKubePortForwardCommand(i.namespaceName, i.resourceName)

	ctx.Export(GetExternalHostnameOutputName(), pulumi.String(i.externalHostname))
	ctx.Export(GetInternalHostnameOutputName(), pulumi.String(i.internalHostname))
	ctx.Export(GetKubeServiceNameOutputName(), pulumi.String(i.kubeServiceName))
	ctx.Export(GetKubeEndpointOutputName(), pulumi.String(i.kubeLocalEndpoint))
	ctx.Export(GetKubePortForwardCommandOutputName(), pulumi.String(kubePortForwardCommand))
	ctx.Export(GetExternalLoadBalancerIp(), pulumi.String(i.externalLoadBalancerIpAddress))
	ctx.Export(GetInternalLoadBalancerIp(), pulumi.String(i.internalLoadBalancerIpAddress))
	ctx.Export(GetNamespaceNameOutputName(), pulumi.String(i.namespaceName))
	ctx.Export(GetAdminUsernameOutputName(), pulumi.String(getAdminUsername()))
	ctx.Export(GetAdminPasswordSecretOutputName(), pulumi.String(i.resourceName))
	return nil
}

func GetExternalHostnameOutputName() string {
	return custom.Name("jenkins-kubernetes-external-hostname")
}

func GetInternalHostnameOutputName() string {
	return custom.Name("jenkins-kubernetes-internal-hostname")
}

func GetAdminPasswordSecretOutputName() string {
	return custom.Name("jenkins-kubernetes-admin-password-secret-name")
}

func GetAdminUsernameOutputName() string {
	return custom.Name("jenkins-kubernetes-admin-username")
}

func getAdminUsername() string {
	return adminUsername
}

func GetKubeServiceNameOutputName() string {
	return custom.Name("jenkins-kubernetes-kubernetes-service-name")
}

func GetKubeEndpointOutputName() string {
	return custom.Name("jenkins-kubernetes-kubernetes-endpoint")
}

func GetKubePortForwardCommandOutputName() string {
	return custom.Name("jenkins-cluster-kube-port-forward-command")
}

func GetExternalLoadBalancerIp() string {
	return custom.Name("jenkins-ingress-external-lb-ip")
}

func GetInternalLoadBalancerIp() string {
	return custom.Name("jenkins-ingress-internal-lb-ip")
}

func GetNamespaceNameOutputName() string {
	return puluminamekubeoutput.Name(kubernetescorev1.Namespace{}, englishword.EnglishWord_namespace.String())
}

// getKubePortForwardCommand returns kubectl port-forward command that can be used by developers.
// ex: "kubectl port-forward -n kubernetes_namespace  service/main-jenkins-kubernetes 8080:8080"
func getKubePortForwardCommand(namespaceName, kubeServiceName string) string {
	return fmt.Sprintf("kubectl port-forward -n %s service/%s %d:%d",
		namespaceName, kubeServiceName, JenkinsPort, JenkinsPort)
}
