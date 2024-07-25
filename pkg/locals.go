package pkg

import (
	"fmt"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/jenkinskubernetes/model"
)

var locals = &Locals{}

type Locals struct {
	IngressCertClusterIssuerName string
	IngressCertSecretName        string
	IngressInternalHostname      string
	IngressExternalHostname      string
	IngressHostnames             []string
}

func LocalsInitializer(jenkinsKubernetes *model.JenkinsKubernetes) {
	if jenkinsKubernetes.Spec.Ingress == nil ||
		!jenkinsKubernetes.Spec.Ingress.IsEnabled ||
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return
	}

	locals.IngressExternalHostname = fmt.Sprintf("%s.%s", jenkinsKubernetes.Metadata.Id,
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressInternalHostname = fmt.Sprintf("%s-internal.%s", jenkinsKubernetes.Metadata.Id,
		jenkinsKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressHostnames = []string{
		locals.IngressExternalHostname,
		locals.IngressInternalHostname,
	}

	//note: a ClusterIssuer resource should have already exist on the kubernetes-cluster.
	//this is typically taken care of by the kubernetes cluster administrator.
	//if the kubernetes-cluster is created using Planton Cloud, then the cluster-issuer name will be
	//same as the ingress-domain-name as long as the same ingress-domain-name is added to the list of
	//ingress-domain-names for the GkeCluster/EksCluster/AksCluster spec.
	locals.IngressCertClusterIssuerName = jenkinsKubernetes.Spec.Ingress.EndpointDomainName

	locals.IngressCertSecretName = jenkinsKubernetes.Metadata.Id
}
