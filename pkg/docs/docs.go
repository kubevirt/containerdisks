package docs

import (
	_ "embed"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	v1 "kubevirt.io/api/core/v1"
)

type TemplateData struct {
	Name        string
	Description string
	Example     string
}

type UserData struct {
	Username       string
	AuthorizedKeys []string
}

type Option func(vm *v1.VirtualMachine)

//go:embed data/cloudinit.tpl
var cloudinitTemplate string

//go:embed data/ignition.tpl
var ignitionTemplate string

//go:embed data/description.tpl
var descriptionTemplate string

func NewVM(name, image string, opts ...Option) *v1.VirtualMachine {
	vm := BasicVM(name, image)

	for _, opt := range opts {
		opt(vm)
	}

	return vm
}

func BasicVM(name, image string) *v1.VirtualMachine {
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
					TerminationGracePeriodSeconds: pointer.Int64(180),
					Domain: v1.DomainSpec{
						Resources: v1.ResourceRequirements{
							Requests: map[k8sv1.ResourceName]resource.Quantity{
								k8sv1.ResourceMemory: resource.MustParse("1Gi"),
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
					},
				},
			},
		},
	}
}

func WithRng() Option {
	return func(vm *v1.VirtualMachine) {
		vm.Spec.Template.Spec.Domain.Devices.Rng = &v1.Rng{}
	}
}

func withCloudInit(volumeSource v1.VolumeSource) Option {
	return func(vm *v1.VirtualMachine) {
		vm.Spec.Template.Spec.Domain.Devices.Disks = append(
			vm.Spec.Template.Spec.Domain.Devices.Disks,
			v1.Disk{
				Name: "cloudinit",
				DiskDevice: v1.DiskDevice{
					Disk: &v1.DiskTarget{
						Bus: "virtio",
					},
				},
			},
		)
		vm.Spec.Template.Spec.Volumes = append(
			vm.Spec.Template.Spec.Volumes,
			v1.Volume{
				Name:         "cloudinit",
				VolumeSource: volumeSource,
			},
		)
	}
}

func WithCloudInitNoCloud(userData string) Option {
	return withCloudInit(v1.VolumeSource{
		CloudInitNoCloud: &v1.CloudInitNoCloudSource{
			UserData: userData,
		},
	})
}

func WithCloudInitConfigDrive(userData string) Option {
	return withCloudInit(v1.VolumeSource{
		CloudInitConfigDrive: &v1.CloudInitConfigDriveSource{
			UserData: userData,
		},
	})
}

func WithSecureBoot() Option {
	return func(vm *v1.VirtualMachine) {
		vm.Spec.Template.Spec.Domain.Features = &v1.Features{
			SMM: &v1.FeatureState{
				Enabled: pointer.Bool(true),
			},
		}
		vm.Spec.Template.Spec.Domain.Firmware = &v1.Firmware{
			Bootloader: &v1.Bootloader{
				EFI: &v1.EFI{
					SecureBoot: pointer.Bool(true),
				},
			},
		}
	}
}

func Template() *template.Template {
	funcMap := template.FuncMap{
		"ToTitle": cases.Title,
	}

	return template.Must(
		template.New("description").Funcs(funcMap).Parse(descriptionTemplate),
	)
}

func CloudInit(data *UserData) string {
	tpl := template.Must(
		template.New("cloudinit").Parse(cloudinitTemplate),
	)

	return mustExecute(tpl, data)
}

func Ignition(data *UserData) string {
	funcMap := template.FuncMap{
		"Quote": func(items []string) []string {
			for i, item := range items {
				items[i] = strconv.Quote(item)
			}
			return items
		},
		"Join": strings.Join,
	}

	tpl := template.Must(
		template.New("ignition").Funcs(funcMap).Parse(ignitionTemplate),
	)

	return mustExecute(tpl, data)
}

func mustExecute(tpl *template.Template, data interface{}) string {
	res := strings.Builder{}

	if err := tpl.Execute(&res, data); err != nil {
		panic(err)
	}

	return res.String()
}
