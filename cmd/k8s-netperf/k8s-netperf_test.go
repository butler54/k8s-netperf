package main

import (
	"testing"

	ocpmetadata "github.com/cloud-bulldozer/go-commons/v2/ocp-metadata"
	"github.com/cloud-bulldozer/k8s-netperf/pkg/metrics"
)

func TestParseAnnotations(t *testing.T) {
	testCases := []struct {
		name    string
		values  []string
		want    map[string]string
		wantErr bool
	}{
		{name: "multiple annotations", values: []string{"example.com/one=value", "example.com/two=value=with=equals"}, want: map[string]string{"example.com/one": "value", "example.com/two": "value=with=equals"}},
		{name: "missing separator", values: []string{"example.com/key"}, wantErr: true},
		{name: "empty key", values: []string{"=value"}, wantErr: true},
		{name: "empty value", values: []string{"example.com/key="}, wantErr: true},
		{name: "whitespace-only key", values: []string{"   =value"}, wantErr: true},
		{name: "whitespace-only value", values: []string{"example.com/key=   "}, wantErr: true},
		{name: "malformed key", values: []string{"not a key=value"}, wantErr: true},
		{name: "duplicate key", values: []string{"example.com/key=one", "example.com/key=two"}, wantErr: true},
		{name: "managed istio key", values: []string{"sidecar.istio.io/inject=false"}, wantErr: true},
		{name: "managed network key", values: []string{"k8s.v1.cni.cncf.io/networks=netperf/network"}, wantErr: true},
		{name: "managed generated network status key", values: []string{"k8s.v1.cni.cncf.io/network-status=netperf/network"}, wantErr: true},
		{name: "managed ovn network key", values: []string{"k8s.ovn.org/pod-networks=netperf/network"}, wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseAnnotations(tc.values)
			if tc.wantErr {
				if err == nil {
					t.Fatal("parseAnnotations() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseAnnotations() error = %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("parseAnnotations() = %#v, want %#v", got, tc.want)
			}
			for key, want := range tc.want {
				if got[key] != want {
					t.Errorf("parseAnnotations()[%q] = %q, want %q", key, got[key], want)
				}
			}
		})
	}
}

func TestParseLabels(t *testing.T) {
	testCases := []struct {
		name    string
		values  []string
		want    map[string]string
		wantErr bool
	}{
		{name: "multiple labels", values: []string{"example.com/team=networking", "environment=test"}, want: map[string]string{"example.com/team": "networking", "environment": "test"}},
		{name: "missing separator", values: []string{"example.com/team"}, wantErr: true},
		{name: "empty key", values: []string{"=value"}, wantErr: true},
		{name: "empty value", values: []string{"example.com/team="}, wantErr: true},
		{name: "invalid key", values: []string{"not a key=value"}, wantErr: true},
		{name: "invalid value", values: []string{"example.com/team=not a value"}, wantErr: true},
		{name: "duplicate key", values: []string{"example.com/team=one", "example.com/team=two"}, wantErr: true},
		{name: "managed role selector", values: []string{"role=server"}, wantErr: true},
		{name: "managed app selector", values: []string{"app=server"}, wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseLabels(tc.values)
			if tc.wantErr {
				if err == nil {
					t.Fatal("parseLabels() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseLabels() error = %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("parseLabels() = %#v, want %#v", got, tc.want)
			}
			for key, want := range tc.want {
				if got[key] != want {
					t.Errorf("parseLabels()[%q] = %q, want %q", key, got[key], want)
				}
			}
		})
	}
}

func TestValidateWorkloadOptions(t *testing.T) {
	testCases := []struct {
		name              string
		runtimeClass      string
		runtimeClassSet   bool
		launchSecurity    string
		launchSecuritySet bool
		pod               bool
		vm                bool
		hostNetOnly       bool
		wantErr           bool
	}{
		{name: "pod runtime class", runtimeClass: "kata", runtimeClassSet: true, pod: true},
		{name: "empty runtime class", runtimeClassSet: true, pod: true, wantErr: true},
		{name: "runtime class with pods disabled", runtimeClass: "kata", runtimeClassSet: true, vm: true, wantErr: true},
		{name: "runtime class with host network only", runtimeClass: "kata", runtimeClassSet: true, pod: true, hostNetOnly: true, wantErr: true},
		{name: "snp launch security", launchSecurity: "snp", launchSecuritySet: true, vm: true},
		{name: "tdx launch security", launchSecurity: "tdx", launchSecuritySet: true, vm: true},
		{name: "launch security without VM", launchSecurity: "snp", launchSecuritySet: true, pod: true, wantErr: true},
		{name: "invalid launch security", launchSecurity: "sev", launchSecuritySet: true, vm: true, wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validateWorkloadOptions(nil, tc.runtimeClass, tc.runtimeClassSet, tc.launchSecurity, tc.launchSecuritySet, tc.pod, tc.vm, tc.hostNetOnly)
			if (err != nil) != tc.wantErr {
				t.Fatalf("validateWorkloadOptions() error = %v, wantErr %t", err, tc.wantErr)
			}
		})
	}
}

func TestValidateAirGappedOptions(t *testing.T) {
	tests := []struct {
		name, image, vmImage, offline, wantNamespace, wantName string
		imageSet, vmImageSet, vm                               bool
		wantErr                                                bool
	}{
		{name: "image override is verbatim", image: "mirror.local/path/netperf:tag", imageSet: true, vm: true},
		{name: "offline DataVolume", offline: "operator-images/prepared", vm: true, wantNamespace: "operator-images", wantName: "prepared"},
		{name: "offline DataVolume requires VM", offline: "operator-images/prepared", wantErr: true},
		{name: "offline DataVolume conflicts with explicit VM image", vmImage: "mirror.local/vm:tag", vmImageSet: true, offline: "operator-images/prepared", vm: true, wantErr: true},
		{name: "malformed DataVolume", offline: "prepared", vm: true, wantErr: true},
		{name: "DataVolume has empty namespace", offline: "/prepared", vm: true, wantErr: true},
		{name: "DataVolume has empty name", offline: "operator-images/", vm: true, wantErr: true},
		{name: "DataVolume has extra path separator", offline: "operator-images/prepared/extra", vm: true, wantErr: true},
		{name: "benchmark namespace source allowed", offline: "netperf/prepared", vm: true, wantNamespace: "netperf", wantName: "prepared"},
		{name: "empty explicit image", imageSet: true, vm: true, wantErr: true},
		{name: "empty explicit VM image", vmImageSet: true, vm: true, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ns, name, err := validateAirGappedOptions(tc.image, tc.imageSet, tc.vmImage, tc.vmImageSet, tc.offline, tc.vm)
			if (err != nil) != tc.wantErr {
				t.Fatalf("validateAirGappedOptions() error = %v, wantErr %t", err, tc.wantErr)
			}
			if ns != tc.wantNamespace || name != tc.wantName {
				t.Errorf("DataVolume = %s/%s, want %s/%s", ns, name, tc.wantNamespace, tc.wantName)
			}
		})
	}
}

func TestApplyClusterDistributionSetsPrometheusFlags(t *testing.T) {
	testCases := []struct {
		name       string
		dist       string
		openShift  bool
		microShift bool
	}{
		{
			name:      "openshift",
			dist:      ocpmetadata.DistributionOpenShift,
			openShift: true,
		},
		{
			name:       "microshift",
			dist:       ocpmetadata.DistributionMicroShift,
			microShift: true,
		},
		{
			name: "kubernetes",
			dist: ocpmetadata.DistributionKubernetes,
		},
		{
			name: "unknown",
			dist: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pcon := metrics.PromConnect{}
			applyClusterDistribution(&pcon, tc.dist)
			if pcon.OpenShift != tc.openShift {
				t.Fatalf("OpenShift = %t, want %t", pcon.OpenShift, tc.openShift)
			}
			if pcon.MicroShift != tc.microShift {
				t.Fatalf("MicroShift = %t, want %t", pcon.MicroShift, tc.microShift)
			}
		})
	}
}

func TestShouldDiscoverPrometheus(t *testing.T) {
	testCases := []struct {
		name                   string
		dist                   string
		url                    string
		metadataAgentAvailable bool
		clusterInfoDegraded    bool
		want                   bool
	}{
		{
			name:                   "discover openshift without explicit prometheus url",
			dist:                   ocpmetadata.DistributionOpenShift,
			metadataAgentAvailable: true,
			want:                   true,
		},
		{
			name:                   "keep openshift discovery with explicit url for token",
			dist:                   ocpmetadata.DistributionOpenShift,
			url:                    "http://127.0.0.1:9090",
			metadataAgentAvailable: true,
			want:                   true,
		},
		{
			name:                   "skip microshift discovery without explicit url",
			dist:                   ocpmetadata.DistributionMicroShift,
			metadataAgentAvailable: true,
		},
		{
			name:                   "skip microshift discovery with explicit url",
			dist:                   ocpmetadata.DistributionMicroShift,
			url:                    "http://127.0.0.1:9090",
			metadataAgentAvailable: true,
		},
		{
			name:                   "skip kubernetes discovery without explicit url",
			dist:                   ocpmetadata.DistributionKubernetes,
			metadataAgentAvailable: true,
		},
		{
			name:                   "skip kubernetes discovery with explicit url",
			dist:                   ocpmetadata.DistributionKubernetes,
			url:                    "http://127.0.0.1:9090",
			metadataAgentAvailable: true,
		},
		{
			name:                   "skip discovery when metadata agent is unavailable",
			dist:                   ocpmetadata.DistributionOpenShift,
			metadataAgentAvailable: false,
		},
		{
			name:                   "discover when cluster info degraded despite explicit url",
			url:                    "http://127.0.0.1:9090",
			metadataAgentAvailable: true,
			clusterInfoDegraded:    true,
			want:                   true,
		},
		{
			name:                   "skip discovery with explicit url when distribution is unknown",
			url:                    "http://127.0.0.1:9090",
			metadataAgentAvailable: true,
		},
		{
			name:                   "discover without explicit url when distribution is unknown",
			metadataAgentAvailable: true,
			want:                   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldDiscoverPrometheus(tc.dist, tc.url, tc.metadataAgentAvailable, tc.clusterInfoDegraded)
			if got != tc.want {
				t.Fatalf("shouldDiscoverPrometheus(%q, %q, %t, %t) = %t, want %t", tc.dist, tc.url, tc.metadataAgentAvailable, tc.clusterInfoDegraded, got, tc.want)
			}
		})
	}
}
