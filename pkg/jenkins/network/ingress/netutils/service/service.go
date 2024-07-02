package service

import (
	"fmt"
	"github.com/plantoncloud-inc/go-commons/kubernetes/network/dns"
)

func GetKubeServiceNameFqdn(jenkinsKubernetesName, namespace string) string {
	return fmt.Sprintf("%s.%s.%s", GetKubeServiceName(jenkinsKubernetesName), namespace, dns.DefaultDomain)
}

func GetKubeServiceName(jenkinsKubernetesName string) string {
	return fmt.Sprintf(jenkinsKubernetesName)
}
