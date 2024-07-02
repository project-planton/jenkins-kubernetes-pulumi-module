package hostname

import (
	"fmt"
)

func GetInternalHostname(jenkinsKubernetesId, environmentName, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s-internal.%s", jenkinsKubernetesId, environmentName, endpointDomainName)
}

func GetExternalHostname(jenkinsKubernetesId, environmentName, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s.%s", jenkinsKubernetesId, environmentName, endpointDomainName)
}
