# {{ .Name | ToTitle }} Containerdisk Images

{{ .Description }}

## Documentation

This image is maintained by [KubeVirt](https://kubevirt.io/) and automatically created from https://github.com/kubevirt/containerdisks.
<br />
<br />
For how to get started with `KubeVirt` visit the [user guide](https://kubevirt.io/user-guide/):
  * [Installation](https://kubevirt.io/user-guide/operations/installation/#installing-kubevirt-on-kubernetes)
  * [Containerdisks](https://kubevirt.io/user-guide/virtual_machines/disks_and_volumes/#containerdisk)
  * [virtctl](https://kubevirt.io/user-guide/user_workloads/virtctl_client_tool)
  * [Creating VirtualMachines by using virtctl](https://kubevirt.io/user-guide/user_workloads/creating_vms/#creating-virtualmachines-by-using-virtctl)

## Examples

### Creating a VirtualMachine and importing this containerdisk with virtctl

You can create a VirtualMachine that will import this containerdisk by running the following command:

```shell
virtctl create vm {{ if .Instancetype }}--instancetype={{ .Instancetype }} {{ end }}{{ if .Preference }}--preference={{ .Preference }} {{ end }}--volume-import=type:registry,url:docker://{{ .Image }},size:10Gi | kubectl create -f -
```

### Creating a VirtualMachine without persistence with virtctl

You can create a VirtualMachine that will use this containerdisk as an ephemeral volume by running the following command:

```shell
virtctl create vm {{ if .Instancetype }}--instancetype={{ .Instancetype }} {{ end }}{{ if .Preference }}--preference={{ .Preference }} {{ end }}--volume-containerdisk=src:{{ .Image }} | kubectl create -f -
```

### Using this containerdisk in a VirtualMachine definition

```yaml
{{ .Example -}}
```