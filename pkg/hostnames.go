package pkg

import "fmt"

func (s *ResourceStack) hostnames() {

}

func GetInternalHostname(jenkinsKubernetesId, endpointDomainName string) string {
	return fmt.Sprintf("%s-internal.%s", jenkinsKubernetesId, endpointDomainName)
}

func GetExternalHostname(jenkinsKubernetesId, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s", jenkinsKubernetesId, endpointDomainName)
}
