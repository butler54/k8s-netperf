package k8s

import (
	"testing"
)

func TestCreateDeploymentWorkloadOptions(t *testing.T) {
	testCases := []struct {
		name           string
		annotations    map[string]string
		workloadLabels map[string]string
	}{
		{name: "defaults preserved"},
		{
			name:           "custom annotations and labels",
			annotations:    map[string]string{"example.com/isolation": "enabled"},
			workloadLabels: map[string]string{"example.com/team": "networking"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := DeploymentParams{
				Name:           "workload",
				Namespace:      "netperf",
				Replicas:       1,
				Image:          "example.com/image",
				Labels:         map[string]string{"app": "workload"},
				Commands:       [][]string{{"sleep", "3600"}},
				Annotations:    tc.annotations,
				WorkloadLabels: tc.workloadLabels,
			}
			deployment := newDeployment(params)
			if got := deployment.Spec.Template.Annotations["sidecar.istio.io/inject"]; got != "true" {
				t.Errorf("Istio annotation = %q, want true", got)
			}
			if got := deployment.Spec.Template.Annotations["example.com/isolation"]; got != tc.annotations["example.com/isolation"] {
				t.Errorf("custom annotation = %q, want %q", got, tc.annotations["example.com/isolation"])
			}
			if got := deployment.Spec.Template.Labels["example.com/team"]; got != tc.workloadLabels["example.com/team"] {
				t.Errorf("custom label = %q, want %q", got, tc.workloadLabels["example.com/team"])
			}
			if got := deployment.Spec.Template.Labels["app"]; got != "workload" {
				t.Errorf("managed label = %q, want workload", got)
			}
		})
	}
}

func TestNewDeploymentManagedLabelsOverrideCustomLabels(t *testing.T) {
	params := DeploymentParams{
		Name:           "workload",
		Namespace:      "netperf",
		Replicas:       1,
		Image:          "example.com/image",
		Labels:         map[string]string{"role": "server"},
		WorkloadLabels: map[string]string{"role": "client", "example.com/team": "networking"},
		Commands:       [][]string{{"sleep", "3600"}},
	}

	deployment := newDeployment(params)
	if got := deployment.Spec.Selector.MatchLabels["role"]; got != "server" {
		t.Errorf("selector role = %q, want server", got)
	}
	if got := deployment.Spec.Template.Labels["role"]; got != "server" {
		t.Errorf("template role = %q, want server", got)
	}
	if got := deployment.Spec.Template.Labels["example.com/team"]; got != "networking" {
		t.Errorf("custom label = %q, want networking", got)
	}
}

func TestNewDeploymentPreservesGeneratedNetworkAnnotationsWithCustomAnnotations(t *testing.T) {
	params := DeploymentParams{
		Name:               "workload",
		Namespace:          "netperf",
		Replicas:           1,
		Image:              "example.com/image",
		Labels:             map[string]string{"app": "workload"},
		Commands:           [][]string{{"sleep", "3600"}},
		NetworkAnnotations: map[string]string{"k8s.v1.cni.cncf.io/networks": "netperf/benchmark"},
		Annotations: map[string]string{
			"example.com/isolation":       "enabled",
			"sidecar.istio.io/inject":     "false",
			"k8s.v1.cni.cncf.io/networks": "user-network",
		},
	}

	deployment := newDeployment(params)
	annotations := deployment.Spec.Template.Annotations
	if got := annotations["k8s.v1.cni.cncf.io/networks"]; got != "netperf/benchmark" {
		t.Errorf("generated network annotation = %q, want %q", got, "netperf/benchmark")
	}
	if got := annotations["sidecar.istio.io/inject"]; got != "true" {
		t.Errorf("default annotation = %q, want true", got)
	}
	if got := annotations["example.com/isolation"]; got != "enabled" {
		t.Errorf("custom annotation = %q, want enabled", got)
	}
}
