# Advanced Usage

## Using External Server
This enables k8s-netperf to use the IP address provided via the `--serverIP` option as server address and the client sends requests to this IP address. This allows dataplane testing between ocp internal client pod and external server.

> *Note: User has to create a server with the provided IP address and run the intented k8s-netperf driver (i.e uperf, iperf or netperf). User has to enable respective ports on this server to allow the traffic from the client*

Once the external server is ready to accept the traffic, users can orhestrate k8s-netperf by running:

```bash
k8s-netperf --serverIP=44.243.95.221
```

## Running with VMs
Running k8s-netperf against Virtual Machines (OpenShift CNV) requires

- OpenShift CNV must be deployed and users should be able to define VMIs
- SSH keys to be present in the home directory `(~/.ssh/id_rsa.pub)`
- OpenShift Routes - k8s-netperf uses this to reach the VMs (k8s-netperf will create the route for the user, but we need Routes)

If the two above are in place, users can orhestrate k8s-netperf to launch VMs by running:

```bash
k8s-netperf --vm
```

### Offline VM DataVolumes

For fully disconnected VM benchmarks, use a bootable Fedora-compatible VM disk in any namespace. The
disk must already include SSH access for `fedora`,
`netperf`/`netserver`, `iperf3`, `uperf`, and their runtime libraries. Build or import the disk using
your approved connected build environment, then deliver it in-cluster as a CDI DataVolume. Do not
expect the benchmark to install packages, clone source, or download build artifacts at boot.

Confirm the source is ready before running:

```shell
$ kubectl get datavolume -n operator-images prepared-netperf
NAME               PHASE
prepared-netperf   Succeeded
$ k8s-netperf --vm --pod=false --offline-data-volume=operator-images/prepared-netperf
```

The source DataVolume must have phase `Succeeded`. k8s-netperf verifies it before creating VMIs,
creates one labeled clone DataVolume per client/server VMI in `netperf`, and attaches each clone as
that VMI's `disk0`. Offline cloud-init only configures guest access/networking, verifies the local
tools, and starts local server processes; it contains no `dnf`, `git`, `curl`, package download, or
source-build fallback. Before running workloads, k8s-netperf also checks the selected client tools and
reports a missing prerequisite rather than triggering an internet download. When the source is outside
`netperf`, cleanup removes the benchmark namespace and its clone volumes. When the source is in
`netperf`, cleanup preserves the source and removes only benchmark-owned clone volumes.

For the standard KubeVirt masquerade pod network, k8s-netperf omits cloud-init network data entirely.
The prepared Fedora image must retain DHCP configuration for its primary guest interface. Explicit
cloud-init network data is used only for secondary or static network configurations.

`--offline-data-volume` is VM-only, uses `NAMESPACE/NAME`, and cannot be combined with an explicitly
specified `--vm-image` because the clone replaces the container disk.

### Workload isolation options

Use `--annotation KEY=VALUE` repeatedly to add the same annotations to every benchmark pod template or VMI. k8s-netperf rejects empty, malformed, duplicate, and tool-managed annotation keys so it can preserve its Istio and network configuration.

Use `--label KEY=VALUE` repeatedly to add labels to every benchmark pod template or VMI. Label keys and values must be valid Kubernetes labels. Empty, malformed, or duplicate labels are rejected before resources are created. The `app` and `role` labels are managed by k8s-netperf for workload and service selection, so they cannot be supplied with `--label`.

```bash
k8s-netperf --label example.com/team=networking --label environment=staging
```

For pod benchmarks, `--runtime-class NAME` sets the pod `runtimeClassName`, for example:

```bash
k8s-netperf --runtime-class kata --annotation example.com/isolation=enabled
```

`--runtime-class` is pod-only and cannot be used with `--pod=false` or `--hostNet`.
When running all scenarios with `--all`, k8s-netperf applies the runtime class only to
pod-network workloads; host-network workloads use the cluster's default runtime.

For VM benchmarks, KubeVirt 1.8.4 or newer is required. Use `--launch-security snp` or `--launch-security tdx` with `--vm` to request the corresponding confidential-VM launch security:

```bash
k8s-netperf --vm --pod=false --launch-security tdx
```

Launch security is VM-only. SNP and TDX VMIs explicitly enable ACPI, use a host-passthrough CPU, the `q35` machine type, and UEFI with Secure Boot disabled. The selected `launchSecurity` member is the only difference. k8s-netperf does not configure TDX attestation. Cluster admission errors are returned to the operator. Ensure that the selected feature gate, scheduling, and hardware capability are available before running the benchmark.

## Using User Defined Network - UDN (only on OCP 4.18 and above)
To run k8s-netperf using a UDN primary network for the test instead of the default network of OVN-k:

For a layer3 UDN:
```
$ k8s-netperf --udnl3
```

For a layer2 UDN:
```
$ k8s-netperf --udnl2
```

It works also with VMs:
```
$ k8s-netperf --udnl2 --vm --udnPluginBinding=l2bridge
```

> Warning! Support of k8s Services with UDN is not fully supported yet, you may faced inconsistent results when using a service in your tests. 

## Cluster User Defined Network - C-UDN
k8s-netperf is able to deploy a C-UDN and then attach network interfaces to the pods (or VMs) that allow use the C-UDN as a secondary network.
The subnet of the C-UDN is `20.0.0.0/16` and the driver will use this secondary interface for the test. C-UDN is automatically cleanup at the end of the execution.
```
# For a layer2 C-UDN:
$ k8s-netperf --cudn layer2

# For a layer3 C-UDN:
$ k8s-netperf --cudn layer3

# Both of the setup works with VMs too:
$ k8s-netperf --cudn layer2 --vm
# or
$ k8s-netperf --cudn layer3 --vm
```
## SR-IOV Network Testing
To run k8s-netperf over SR-IOV Virtual Functions, the SR-IOV Network Operator must be installed on the cluster. k8s-netperf will automatically create the required `SriovNetworkNodePolicy` and `SriovNetwork` CRs, and clean them up after the test.

### Pod SR-IOV
Pass the PF (Physical Function) interface name to the `--sriov` flag:
```bash
k8s-netperf --sriov ens2f0np0
```

On a single-node cluster:
```bash
k8s-netperf --sriov ens2f0np0 --local
```

By default, the SR-IOV policy targets nodes with the `node-role.kubernetes.io/worker` label. If SR-IOV NICs are only available on nodes with a different role label, use `--sriov-node-selector`:
```bash
k8s-netperf --sriov ens2f0np0 --sriov-node-selector cnf-worker
```

### VM SR-IOV
SR-IOV is also supported with KubeVirt VMs. When `--vm` is used with `--sriov`, VFs are created with `deviceType: vfio-pci` and passed through to the guest via PCI passthrough. The IP address assigned by whereabouts to the virt-launcher pod is configured inside the guest VM using `nmcli` over virtctl SSH.

```bash
k8s-netperf --sriov ens1f0 --vm --use-virtctl --pod=false
```

Requirements for VM SR-IOV:
- OpenShift CNV (KubeVirt) must be deployed
- The PF must use a driver supported by the VM containerdisk (e.g. Intel `ice`/`iavf` drivers are included in Fedora 39; Mellanox `mlx5_core` is not)
- SSH keys must be present in `~/.ssh/id_rsa.pub`
- `virtctl` is recommended (`--use-virtctl`) for SSH access into the guest VMs

> Note: Pod and VM SR-IOV tests cannot run in the same invocation because pods require `deviceType: netdevice` while VMs require `deviceType: vfio-pci`. Run them separately with `--pod=false` for VM-only or without `--vm` for pod-only.

> Note: `--sriov` is mutually exclusive with `--bridge`, `--macvlan`, `--ib-write-bw`, `--hostNet` and UDN flags (`--udnl2`, `--udnl3`, `--cudn`).

## MACVLAN Network Testing
To run k8s-netperf over a MACVLAN interface, pass the master (host) interface name to the `--macvlan` flag. k8s-netperf will automatically create a MACVLAN `NetworkAttachmentDefinition` with [whereabouts](https://github.com/k8snetworkplumbingwg/whereabouts) IPAM in the `netperf` namespace and clean it up after the test.

Prerequisites:
- [Multus CNI](https://github.com/k8snetworkplumbingwg/multus-cni) must be installed on the cluster
- The `macvlan` CNI plugin must be available on the nodes (`/opt/cni/bin/macvlan`)
- The whereabouts IPAM plugin must be installed

```bash
k8s-netperf --macvlan eth0
```

On a single-node cluster:
```bash
k8s-netperf --macvlan eth0 --local
```

> Note: `--macvlan` is mutually exclusive with `--bridge`, `--sriov`, `--ib-write-bw`, `--hostNet` and UDN flags (`--udnl2`, `--udnl3`, `--cudn`).

## Using a Linux Bridge Interface
When using `--bridge`, a NetworkAttachmentDefinition defining a bridge interface is attached to the VMs and is used for the test. It requires the name of the bridge as it is defined in the NetworkNodeConfigurationPolicy, NMstate operator is required.

For example:
```yaml
apiVersion: nmstate.io/v1alpha1
kind: NodeNetworkConfigurationPolicy
metadata:
  name: br0-eth1
spec:
  desiredState:
    interfaces:
      - name: br0
        description: Linux bridge with eno2 as a port
        type: linux-bridge
        state: up
        ipv4:
          dhcp: true
          enabled: true
        bridge:
          options:
            stp:
              enabled: false
          port:
            - name: eno2
```

Then you can launch a test using the bridge interface:
```bash
./bin/amd64/k8s-netperf --vm --bridge br0
```

By default, it will read the `bridgeNetwork.json` file from the git repository. If the default IP addresses (10.10.10.12/24 and 10.10.10.14/24) are not available for your setup, it is possible to change it by passing a JSON file as a parameter with `--bridgeNetwork`, like follow:
```bash
k8s-netperf --vm --bridge br0 --bridgeNetwork /path/to/my/bridgeConfig.json
```

## Using OVN Localnet
When using `--localnet`, k8s-netperf creates a Localnet ClusterUserDefinedNetwork (C-UDN) and attaches a secondary interface to the VMs for the test. The value passed to `--localnet` must match the OVN-K external / physical network name configured in a NodeNetworkConfigurationPolicy bridge mapping. NMState operator is required. Localnet is VM-only (`--vm --pod=false`).

For example, map the localnet name `physnet` to the existing OVS bridge `br-ex`:
```yaml
apiVersion: nmstate.io/v1
kind: NodeNetworkConfigurationPolicy
metadata:
  name: mapping-localnet
spec:
  nodeSelector:
    node-role.kubernetes.io/worker: ''
  desiredState:
    ovn:
      bridge-mappings:
      - localnet: physnet
        bridge: br-ex
        state: present
```

Alternatively, create a dedicated OVS bridge on a NIC (e.g. `ens7f0`) and map a localnet to it:
```yaml
apiVersion: nmstate.io/v1
kind: NodeNetworkConfigurationPolicy
metadata:
  name: ens7f0-ovs-underlay
spec:
  nodeSelector:
    node-role.kubernetes.io/worker: ""
  desiredState:
    interfaces:
      - name: br-ens7f0
        description: "OVS Bridge dedicated to ens7f0 for CUDN Localnets"
        type: ovs-bridge
        state: up
        bridge:
          allow-extra-patch-ports: true
          port:
            - name: ens7f0
    ovn:
      bridge-mappings:
        - localnet: physnet-ens7f0
          bridge: br-ens7f0
          state: present
```

Then you can launch a test using the localnet name from the bridge mapping:
```bash
./bin/amd64/k8s-netperf --vm --pod=false --localnet physnet
# or, with the dedicated bridge example above:
./bin/amd64/k8s-netperf --vm --pod=false --localnet physnet-ens7f0
```

By default, it will read the `localnetNetwork.json` file from the git repository. If the default IP addresses (192.168.200.10/24 and 192.168.200.11/24) are not available for your setup, it is possible to change it by passing a JSON file as a parameter with `--localnet-config`, like follow:
```bash
k8s-netperf --vm --pod=false --localnet physnet --localnet-config /path/to/my/localnetConfig.json
```

> Note: `--localnet` is mutually exclusive with `--bridge`, `--sriov`, `--macvlan`, `--ib-write-bw`, `--hostNet` and UDN flags (`--udnl2`, `--udnl3`, `--cudn`).


## Privileged pods

If your use case requires running pods with privileged security context, use the `--privileged` flag:
```
$ k8s-netperf --privileged
```

## RoCEv2 testing

Using the `ib_write_bw` driver, your hardware should include RDMA devices:
```
$ k8s-netperf --ib-write-bw nic:gid --privileged --hostNet
```

### RoCEv2 local testing

On Fedora systems:
```bash
sudo modprobe rdma_rxe
sudo rdma link add rxe0 type rxe netdev eth0  # Name of your active interface
kind create cluster --config testing/kind-config-rdma.yaml
kubectl label node kind-control-plane node-role.kubernetes.io/worker=""
k8s-netperf --config config.yaml --hostNet --privileged --ib-write-bw rxe0:1 --local
```

Cleanup:
```bash
kind delete cluster
sudo rdma link delete rxe0
```
