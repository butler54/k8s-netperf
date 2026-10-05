package k8s

import (
	"context"
	"strings"
	"testing"

	"github.com/cloud-bulldozer/k8s-netperf/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreateDeploymentWorkloadOptions(t *testing.T) {
	testCases := []struct {
		name             string
		annotations      map[string]string
		workloadLabels   map[string]string
		runtimeClass     string
		hostNetwork      bool
		wantRuntimeClass string
	}{
		{name: "defaults preserved"},
		{
			name:             "custom annotations and runtime class",
			annotations:      map[string]string{"example.com/isolation": "enabled"},
			workloadLabels:   map[string]string{"example.com/team": "networking"},
			runtimeClass:     "kata",
			wantRuntimeClass: "kata",
		},
		{
			name:         "host network omits runtime class",
			runtimeClass: "kata",
			hostNetwork:  true,
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
				RuntimeClass:   tc.runtimeClass,
				HostNetwork:    tc.hostNetwork,
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
			if tc.wantRuntimeClass == "" {
				if deployment.Spec.Template.Spec.RuntimeClassName != nil {
					t.Errorf("RuntimeClassName = %q, want nil", *deployment.Spec.Template.Spec.RuntimeClassName)
				}
			} else if deployment.Spec.Template.Spec.RuntimeClassName == nil || *deployment.Spec.Template.Spec.RuntimeClassName != tc.wantRuntimeClass {
				t.Errorf("RuntimeClassName = %v, want %q", deployment.Spec.Template.Spec.RuntimeClassName, tc.wantRuntimeClass)
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

func TestBenchmarkPodImage(t *testing.T) {
	tests := []struct{ override, want string }{
		{"", k8sNetperfImage},
		{"registry.example.local/custom/netperf:v1", "registry.example.local/custom/netperf:v1"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := benchmarkPodImage(&config.PerfScenarios{PodImage: tc.override}); got != tc.want {
				t.Errorf("benchmarkPodImage() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestOfflineDataVolumeValidationAndClone(t *testing.T) {
	ctx := context.Background()
	source := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "cdi.kubevirt.io/v1beta1", "kind": "DataVolume",
		"metadata": map[string]interface{}{"name": "prepared", "namespace": "images"},
		"spec":     map[string]interface{}{"pvc": map[string]interface{}{"resources": map[string]interface{}{"requests": map[string]interface{}{"storage": "20Gi"}}}},
		"status":   map[string]interface{}{"phase": "Succeeded"},
	}}
	dyn := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), source)
	if err := CreateOfflineDataVolumeClone(ctx, dyn, "images", "prepared", "client-disk", "run-123"); err != nil {
		t.Fatalf("CreateOfflineDataVolumeClone() error = %v", err)
	}
	clone, err := dyn.Resource(dataVolumeGVR).Namespace(namespace).Get(ctx, "client-disk", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get clone: %v", err)
	}
	if got, _, _ := unstructured.NestedString(clone.Object, "spec", "source", "pvc", "namespace"); got != "images" {
		t.Errorf("clone source namespace = %q, want images", got)
	}
	if got, _, _ := unstructured.NestedString(clone.Object, "spec", "source", "pvc", "name"); got != "prepared" {
		t.Errorf("clone source name = %q, want prepared", got)
	}
	if clone.GetLabels()[offlineCloneLabel] != "true" {
		t.Errorf("clone label = %q, want true", clone.GetLabels()[offlineCloneLabel])
	}
	if clone.GetLabels()[offlineCloneRunLabel] != "run-123" {
		t.Errorf("clone run label = %q, want run-123", clone.GetLabels()[offlineCloneRunLabel])
	}
	if got, _, _ := unstructured.NestedString(clone.Object, "spec", "pvc", "resources", "requests", "storage"); got != "20Gi" {
		t.Errorf("clone storage = %q, want 20Gi", got)
	}
	if _, err := dyn.Resource(dataVolumeGVR).Namespace("images").Get(ctx, "prepared", metav1.GetOptions{}); err != nil {
		t.Errorf("source DataVolume was not preserved: %v", err)
	}
}

func TestValidateOfflineDataVolumeRejectsUnavailableSource(t *testing.T) {
	ctx := context.Background()
	for _, source := range []*unstructured.Unstructured{nil, {Object: map[string]interface{}{
		"apiVersion": "cdi.kubevirt.io/v1beta1", "kind": "DataVolume",
		"metadata": map[string]interface{}{"name": "pending", "namespace": "images"},
		"status":   map[string]interface{}{"phase": "ImportInProgress"},
	}}} {
		dyn := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme())
		name := "missing"
		if source != nil {
			dyn = dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), source)
			name = "pending"
		}
		if err := ValidateOfflineDataVolume(ctx, dyn, "images", name); err == nil {
			t.Error("ValidateOfflineDataVolume() error = nil, want unavailable-source error")
		}
	}
}

func TestCreateOfflineDataVolumeCloneDoesNotCreateCloneFromInvalidSource(t *testing.T) {
	ctx := context.Background()
	source := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "cdi.kubevirt.io/v1beta1", "kind": "DataVolume",
		"metadata": map[string]interface{}{"name": "prepared", "namespace": "images"},
		"spec":     map[string]interface{}{"pvc": map[string]interface{}{}},
		"status":   map[string]interface{}{"phase": "Succeeded"},
	}}
	dyn := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), source)

	err := CreateOfflineDataVolumeClone(ctx, dyn, "images", "prepared", "client-disk", "run-123")
	if err == nil || !strings.Contains(err.Error(), "no requested storage") {
		t.Fatalf("CreateOfflineDataVolumeClone() error = %v, want missing-storage error", err)
	}
	if _, err := dyn.Resource(dataVolumeGVR).Namespace(namespace).Get(ctx, "client-disk", metav1.GetOptions{}); !apierrors.IsNotFound(err) {
		t.Fatalf("clone exists or lookup failed with %v, want not found", err)
	}
}

func TestCreateOfflineDataVolumeClonePreservesStorageAndRejectsUnownedClone(t *testing.T) {
	ctx := context.Background()
	source := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "cdi.kubevirt.io/v1beta1", "kind": "DataVolume",
		"metadata": map[string]interface{}{"name": "prepared", "namespace": "images"},
		"spec": map[string]interface{}{"storage": map[string]interface{}{
			"storageClassName": "fast", "accessModes": []interface{}{"ReadWriteMany"}, "volumeMode": "Block",
			"resources": map[string]interface{}{"requests": map[string]interface{}{"storage": "20Gi"}},
			"selector":  map[string]interface{}{"matchLabels": map[string]interface{}{"disk": "benchmark"}},
		}}, "status": map[string]interface{}{"phase": "Succeeded"},
	}}
	dyn := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), source)
	if err := CreateOfflineDataVolumeClone(ctx, dyn, "images", "prepared", "server-disk-run", "run-123"); err != nil {
		t.Fatal(err)
	}
	clone, err := dyn.Resource(dataVolumeGVR).Namespace(namespace).Get(ctx, "server-disk-run", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got, _, _ := unstructured.NestedString(clone.Object, "spec", "storage", "storageClassName"); got != "fast" {
		t.Errorf("storage class = %q, want fast", got)
	}
	if got, _, _ := unstructured.NestedString(clone.Object, "spec", "storage", "volumeMode"); got != "Block" {
		t.Errorf("volume mode = %q, want Block", got)
	}
	if modes, _, _ := unstructured.NestedStringSlice(clone.Object, "spec", "storage", "accessModes"); len(modes) != 1 || modes[0] != "ReadWriteMany" {
		t.Errorf("access modes = %v, want [ReadWriteMany]", modes)
	}
	if got, _, _ := unstructured.NestedString(clone.Object, "spec", "storage", "resources", "requests", "storage"); got != "20Gi" {
		t.Errorf("storage request = %q, want 20Gi", got)
	}
	if got, _, _ := unstructured.NestedString(clone.Object, "spec", "storage", "selector", "matchLabels", "disk"); got != "benchmark" {
		t.Errorf("selector disk = %q, want benchmark", got)
	}
	if err := CreateOfflineDataVolumeClone(ctx, dyn, "images", "prepared", "server-disk-run", "other-run"); err == nil {
		t.Fatal("expected reused clone from another run to be rejected")
	}
}

func TestOfflineCloneNameUsesRunUUID(t *testing.T) {
	if got, want := offlineCloneName("client", "ABC-123"), "client-disk-abc123"; got != want {
		t.Errorf("offlineCloneName() = %q, want %q", got, want)
	}
}

func TestDestroyBenchmarkResourcesPreservesSourceDataVolume(t *testing.T) {
	ctx := context.Background()
	client := fake.NewClientset(
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "benchmark", Namespace: namespace, Labels: map[string]string{benchmarkManagedLabel: "true"}}},
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "benchmark", Namespace: namespace, Labels: map[string]string{benchmarkManagedLabel: "true"}}},
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "operator-workload", Namespace: namespace}},
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "operator-service", Namespace: namespace}},
	)
	source := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "cdi.kubevirt.io/v1beta1", "kind": "DataVolume",
		"metadata": map[string]interface{}{"name": "source", "namespace": namespace},
	}}
	clone := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "cdi.kubevirt.io/v1beta1", "kind": "DataVolume",
		"metadata": map[string]interface{}{"name": "clone", "namespace": namespace, "labels": map[string]interface{}{offlineCloneLabel: "true"}},
	}}
	scheme := runtime.NewScheme()
	dyn := dynamicfake.NewSimpleDynamicClientWithCustomListKinds(scheme, map[schema.GroupVersionResource]string{
		dataVolumeGVR: "DataVolumeList",
		vmiGVR:        "VirtualMachineInstanceList",
		routeGVR:      "RouteList",
	}, source, clone)

	if err := DestroyBenchmarkResources(client, dyn); err != nil {
		t.Fatalf("DestroyBenchmarkResources() error = %v", err)
	}
	if _, err := dyn.Resource(dataVolumeGVR).Namespace(namespace).Get(ctx, "source", metav1.GetOptions{}); err != nil {
		t.Errorf("source DataVolume was deleted: %v", err)
	}
	if _, err := dyn.Resource(dataVolumeGVR).Namespace(namespace).Get(ctx, "clone", metav1.GetOptions{}); !apierrors.IsNotFound(err) {
		t.Errorf("clone DataVolume lookup error = %v, want not found", err)
	}
	if _, err := client.AppsV1().Deployments(namespace).Get(ctx, "benchmark", metav1.GetOptions{}); !apierrors.IsNotFound(err) {
		t.Errorf("managed Deployment lookup error = %v, want not found", err)
	}
	if _, err := client.CoreV1().Services(namespace).Get(ctx, "benchmark", metav1.GetOptions{}); !apierrors.IsNotFound(err) {
		t.Errorf("managed Service lookup error = %v, want not found", err)
	}
	if _, err := client.AppsV1().Deployments(namespace).Get(ctx, "operator-workload", metav1.GetOptions{}); err != nil {
		t.Errorf("unmanaged Deployment was deleted: %v", err)
	}
	if _, err := client.CoreV1().Services(namespace).Get(ctx, "operator-service", metav1.GetOptions{}); err != nil {
		t.Errorf("unmanaged Service was deleted: %v", err)
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

func TestDeploymentFailureError(t *testing.T) {
	testCases := []struct {
		name       string
		dp         DeploymentParams
		conditions []appsv1.DeploymentCondition
		want       string
	}{
		{
			name: "unavailable runtime class",
			dp:   DeploymentParams{Name: "server", RuntimeClass: "kata"},
			conditions: []appsv1.DeploymentCondition{{
				Type:    appsv1.DeploymentReplicaFailure,
				Status:  corev1.ConditionTrue,
				Reason:  "FailedCreate",
				Message: `pods "server" is forbidden: runtimeclass "kata" not found`,
			}},
			want: `pod workload with runtime class "kata"`,
		},
		{
			name: "unschedulable workload",
			dp:   DeploymentParams{Name: "client"},
			conditions: []appsv1.DeploymentCondition{{
				Type:    appsv1.DeploymentProgressing,
				Status:  corev1.ConditionFalse,
				Reason:  "ProgressDeadlineExceeded",
				Message: "ReplicaSet has timed out progressing.",
			}},
			want: "workload may be unschedulable",
		},
		{
			name: "no failure",
			dp:   DeploymentParams{Name: "server"},
			conditions: []appsv1.DeploymentCondition{{
				Type:   appsv1.DeploymentAvailable,
				Status: corev1.ConditionTrue,
			}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: tc.dp.Name},
				Status:     appsv1.DeploymentStatus{Conditions: tc.conditions},
			}
			err := deploymentFailureError(d, tc.dp)
			if tc.want == "" {
				if err != nil {
					t.Fatalf("deploymentFailureError() error = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("deploymentFailureError() error = %v, want substring %q", err, tc.want)
			}
		})
	}
}

func TestWaitForReadyReturnsAlreadyReadyDeployment(t *testing.T) {
	client := fake.NewClientset(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "ready", Namespace: "netperf"},
		Status:     appsv1.DeploymentStatus{ReadyReplicas: 1},
	})

	ready, err := WaitForReady(client, DeploymentParams{Name: "ready", Namespace: "netperf"})
	if err != nil {
		t.Fatalf("WaitForReady() error = %v", err)
	}
	if !ready {
		t.Fatal("WaitForReady() = false, want true")
	}
}

func TestDeploymentReadyTimeoutErrorIncludesImageReference(t *testing.T) {
	err := deploymentReadyTimeoutError(DeploymentParams{
		Name:      "server",
		Namespace: "netperf",
		Image:     "mirror.example.local/bench/netperf:v1",
	})
	if got := err.Error(); !strings.Contains(got, "mirror.example.local/bench/netperf:v1") || !strings.Contains(got, "image pull failures") {
		t.Errorf("deploymentReadyTimeoutError() = %q, want image reference and actionable guidance", got)
	}
}
