package pkg

var vars = struct {
	JenkinsAdminUsername          string
	JenkinsAdminPasswordSecretKey string
	JenkinsDockerImageTag         string
	HelmChartVersion              string
	HelmChartName                 string
	HelmChartRepoUrl              string
	IstioIngressNamespace         string
}{
	JenkinsAdminUsername:          "admin",
	JenkinsAdminPasswordSecretKey: "jenkins-admin-password",
	JenkinsDockerImageTag:         "2.454-jdk17",
	HelmChartVersion:              "5.1.5",
	HelmChartName:                 "jenkins",
	HelmChartRepoUrl:              "https://charts.jenkins.io",
	IstioIngressNamespace:         "istio-ingress",
}
