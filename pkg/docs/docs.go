package docs

import (
	_ "embed"
	"strings"
	"text/template"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	v1 "kubevirt.io/api/core/v1"
)

//go:embed data/cloudinit.txt
var cloudinitFixture string

//go:embed data/ignition.txt
var ignitionFixture string

//go:embed data/descripiton.tpl
var defaultTemplate string

func BasicVirtualMachine(name, image, userData string) *v1.VirtualMachine {
	always := v1.RunStrategyAlways
	return &v1.VirtualMachine{
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualMachine",
			APIVersion: "kubevirt.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.VirtualMachineSpec{
			RunStrategy: &always,
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				Spec: v1.VirtualMachineInstanceSpec{
					TerminationGracePeriodSeconds: pointer.Int64Ptr(180),
					Domain: v1.DomainSpec{
						Resources: v1.ResourceRequirements{
							Requests: map[k8sv1.ResourceName]resource.Quantity{
								k8sv1.ResourceMemory: resource.MustParse("1Gi"),
								k8sv1.ResourceCPU:    resource.MustParse("2"),
							},
						},
						Devices: v1.Devices{
							Disks: []v1.Disk{
								{
									Name: "containerdisk",
									DiskDevice: v1.DiskDevice{
										Disk: &v1.DiskTarget{
											Bus: "virtio",
										},
									},
								},
								{
									Name: "cloudinit",
									DiskDevice: v1.DiskDevice{
										Disk: &v1.DiskTarget{
											Bus: "virtio",
										},
									},
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "containerdisk",
							VolumeSource: v1.VolumeSource{
								ContainerDisk: &v1.ContainerDiskSource{
									Image: image,
								},
							},
						},
						{
							Name: "cloudinit",
							VolumeSource: v1.VolumeSource{
								CloudInitNoCloud: &v1.CloudInitNoCloudSource{
									UserData: userData,
								},
							},
						},
					},
				},
			},
		},
	}
}

func CloudInit() string {
	return cloudinitFixture
}

func Ignition() string {
	return ignitionFixture
}

func Template() *template.Template {
	funcMap := template.FuncMap{
		"ToTitle": strings.Title,
	}
	tpl, err := template.New("description").Funcs(funcMap).Parse(defaultTemplate)
	if err != nil {
		panic(err)
	}
	return tpl
}

type TemplateData struct {
	Name        string
	Description string
	Example     string
}
