package pkg

import "fmt"

// getKubePortForwardCommand returns kubectl port-forward command that can be used by developers.
// ex: "kubectl port-forward -n kubernetes_namespace  service/main-jenkins-kubernetes 8080:8080"
func (s *ResourceStack) exportKubePortForwardCommand(namespaceName, kubeServiceName string) string {
	return fmt.Sprintf("kubectl port-forward -n %s service/%s %d:%d",
		namespaceName, kubeServiceName, JenkinsPort, JenkinsPort)
}
