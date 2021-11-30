package hashsum

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

func TestParse(t *testing.T) {
	type args struct {
		fileName       string
		checksumFormat ChecksumFormat
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			args: struct {
				fileName       string
				checksumFormat ChecksumFormat
			}{fileName: "testdata/bsd.checksum", checksumFormat: ChecksumFormatBSD},
			want: map[string]string{
				"CentOS-8-Container-8.3.2011-20201204.2.x86_64.tar.xz":               "5b141656a073acf7ada9a6e809cc54b5202a8cff9acaef879522a0dfb8273676",
				"CentOS-8-GenericCloud-8.3.2011-20201204.2.x86_64.qcow2":             "7ec97062618dc0a7ebf211864abf63629da1f325578868579ee70c495bed3ba0",
				"CentOS-8-Vagrant-8.0.1905-1.x86_64.vagrant-virtualbox.box":          "a0823c8c1d48024e44627e5fbfc55cb2c14bd8997aa002fc037202e7ec543e2b",
				"CentOS-8-Vagrant-8.1.1911-20200113.3.x86_64.vagrant-libvirt.box":    "aeafc5ad5cdb2c0ae5daa996384e4ee284c9e7b0a3b39f48efe92b2ea8e39aaa",
				"CentOS-8-Vagrant-8.1.1911-20200113.3.x86_64.vagrant-virtualbox.box": "88b225f61abda0c59db7d7b8dd980e3e45cd08ab8bac5ab2ae872dbb35bd80e4",
				"CentOS-8-Container-8.2.2004-20200611.2.x86_64.tar.xz":               "9f4b3d3ed01917e4d16f1e6c218f6b904bbd5714b38f147f4ff46657810c7555",
				"CentOS-8-GenericCloud-8.2.2004-20200611.2.x86_64.qcow2":             "d8984b9baee57b127abce310def0f4c3c9d5b3cea7ea8451fc4ffcbc9935b640",
				"CentOS-8-Vagrant-8.2.2004-20200611.2.x86_64.vagrant-libvirt.box":    "e91d44d96c64f015ae943d66525f4c9d763c28b91440d5741d79226c44c45f86",
				"CentOS-8-Container-8.4.2105-20210603.0.x86_64.tar.xz":               "029d18736e3c6fbc65fd4ffe79a8beefd945fb4d9bd1c909aa8909b91ad39ed3",
				"CentOS-8-GenericCloud-8.4.2105-20210603.0.x86_64.qcow2":             "3510fc7deb3e1939dbf3fe6f65a02ab1efcc763480bc352e4c06eca2e4f7c2a2",
				"CentOS-8-Vagrant-8.0.1905-1.x86_64.vagrant-libvirt.box":             "0b93b761fcaca720920095bdf606d3c4b87daf56ecd4e6213e5a51b90f4fe284",
				"CentOS-8-GenericCloud-8.1.1911-20200113.3.x86_64.qcow2":             "e2cf1081645b1089f574918fb808b32d247169ec4ec1a13bca9e14a74df6530e",
				"CentOS-8-ec2-8.3.2011-20201204.2.x86_64.qcow2":                      "f5cc7165940e991b0b7fa6f5bf866bcb941ecbb9a9942afb043b3f31940e88d9",
				"CentOS-8-Vagrant-8.3.2011-20201204.2.x86_64.vagrant-libvirt.box":    "9f26496a13bb560ae55d21c4e4b6f93bfb994a7b56923e6fc858845ac4c69e30",
				"CentOS-8-Vagrant-8.4.2105-20210603.0.x86_64.vagrant-virtualbox.box": "dfe4a34e59eb3056a6fe67625454c3607cbc52ae941aeba0498c29ee7cb9ac22",
				"CentOS-8-Container-8.1.1911-20200113.3-layer.x86_64.tar.xz":         "81f75692e369b8759a25ab58b836f7ed7f7339df21e3fbbe39abc0e5f71ed001",
				"CentOS-8-Container-8.1.1911-20200113.3.x86_64.tar.xz":               "6e2c208e29b29f0c7a2b29ecae9cbcff28c0bd2b178041b59c2e0763a08504ca",
				"CentOS-8-ec2-8.2.2004-20200611.2.x86_64.qcow2":                      "74d8e7cdc62b3ac5a1719c642a8aa4c9915ca65d86e4dff9a446be44acf13c37",
				"CentOS-8-ec2-8.1.1911-20200113.3.x86_64.qcow2":                      "1d14d3f52a24e174fa273fac5cb99a573deca8351cde0847cabefb1a8bb93a02",
				"CentOS-8-Vagrant-8.2.2004-20200611.2.x86_64.vagrant-virtualbox.box": "698b0d9c6c3f31a4fd1c655196a5f7fc224434112753ab6cb3218493a86202de",
				"CentOS-8-Vagrant-8.3.2011-20201204.2.x86_64.vagrant-virtualbox.box": "fee51a026c1caa9d88a8c74f09352ef4b7606952285cdf2888ea062a8eee499f",
				"CentOS-8-ec2-8.4.2105-20210603.0.x86_64.qcow2":                      "da9c4abe7009954f8b54339c6d671b24b1b2df7350b5d3e9a0aed0e127b74ca1",
				"CentOS-8-Vagrant-8.4.2105-20210603.0.x86_64.vagrant-libvirt.box":    "37cc017738bf12cafce3a97c4de73526452da6332cc5c0516723988644e62620",
			},
		},
		{
			args: struct {
				fileName       string
				checksumFormat ChecksumFormat
			}{fileName: "testdata/gnu.checksum", checksumFormat: ChecksumFormatGNU},
			want: map[string]string{
				"rhcos-4.9.0-x86_64-aws.x86_64.vmdk.gz":        "006896e8a02f6d5f0950cb97f5be904c0d052570254616e3ca05d3e4a76b2710",
				"rhcos-4.9.0-x86_64-azure.x86_64.vhd.gz":       "1c0de512132c239614ef9a8b9be6c8b5692ed405990dd64daa1f4f391bbb7382",
				"rhcos-4.9.0-x86_64-azurestack.x86_64.vhd.gz":  "68e219825af597580aaf60930c08966d3304182e259b4744ada54d1409865fd3",
				"rhcos-4.9.0-x86_64-gcp.x86_64.tar.gz":         "031cbf6a3c00e89383a42266d00e48501aebaa4b9aaf2a1f80e44e90d1b66a82",
				"rhcos-4.9.0-x86_64-ibmcloud.x86_64.qcow2.gz":  "c5c6be77aac71d93522a5099464f0cb8084db304dec5693933e8b970a7885185",
				"rhcos-4.9.0-x86_64-live-initramfs.x86_64.img": "54a07a62f336f760c61641a0ec54f70eefc6cc262399ef9bdd38376c88d8c9bd",
				"rhcos-4.9.0-x86_64-live-kernel-x86_64":        "d13269e6c60119397210418781b7057673c4018692d28a868e248a0b550ea247",
				"rhcos-4.9.0-x86_64-live-rootfs.x86_64.img":    "3b5ba1e98d9852907aaffe4acda9af7b293bdb88008f310631a8a5b4333ea378",
				"rhcos-4.9.0-x86_64-live.x86_64.iso":           "0e92c3ad698ef68057011f7cc5b9fd07356b8711a55f735aaae22c91b996c96e",
				"rhcos-4.9.0-x86_64-metal.x86_64.raw.gz":       "ef9a304cba0c0050486965e38b3c6c614c0646af0b0de493c99eb6a14703cb5a",
				"rhcos-4.9.0-x86_64-metal4k.x86_64.raw.gz":     "9b4548b8b87322dd4d659922cddd287060d6b4ac53992d121ac32442a471f793",
				"rhcos-4.9.0-x86_64-openstack.x86_64.qcow2.gz": "3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8",
				"rhcos-4.9.0-x86_64-ostree.x86_64.tar":         "6708a9fbf379a2e53a9b61ae18fad2fb472bd746d4a04638dbb412195de10a24",
				"rhcos-4.9.0-x86_64-qemu.x86_64.qcow2.gz":      "cae8928e0cd35b88fcec7c07b1072155bde17d7dd44985f8b0d9e3862c556602",
				"rhcos-4.9.0-x86_64-vmware.x86_64.ova":         "6c8bfdee5930f12368b9f46a11aea736a068208262f7747f3bac54eb581531f5",
				"rhcos-aws.x86_64.vmdk.gz":                     "006896e8a02f6d5f0950cb97f5be904c0d052570254616e3ca05d3e4a76b2710",
				"rhcos-azure.x86_64.vhd.gz":                    "1c0de512132c239614ef9a8b9be6c8b5692ed405990dd64daa1f4f391bbb7382",
				"rhcos-azurestack.x86_64.vhd.gz":               "68e219825af597580aaf60930c08966d3304182e259b4744ada54d1409865fd3",
				"rhcos-gcp.x86_64.tar.gz":                      "031cbf6a3c00e89383a42266d00e48501aebaa4b9aaf2a1f80e44e90d1b66a82",
				"rhcos-ibmcloud.x86_64.qcow2.gz":               "c5c6be77aac71d93522a5099464f0cb8084db304dec5693933e8b970a7885185",
				"rhcos-installer-initramfs.x86_64.img":         "54a07a62f336f760c61641a0ec54f70eefc6cc262399ef9bdd38376c88d8c9bd",
				"rhcos-installer-kernel-x86_64":                "d13269e6c60119397210418781b7057673c4018692d28a868e248a0b550ea247",
				"rhcos-installer-rootfs.x86_64.img":            "3b5ba1e98d9852907aaffe4acda9af7b293bdb88008f310631a8a5b4333ea378",
				"rhcos-live-initramfs.x86_64.img":              "54a07a62f336f760c61641a0ec54f70eefc6cc262399ef9bdd38376c88d8c9bd",
				"rhcos-live-kernel-x86_64":                     "d13269e6c60119397210418781b7057673c4018692d28a868e248a0b550ea247",
				"rhcos-live-rootfs.x86_64.img":                 "3b5ba1e98d9852907aaffe4acda9af7b293bdb88008f310631a8a5b4333ea378",
				"rhcos-live.x86_64.iso":                        "0e92c3ad698ef68057011f7cc5b9fd07356b8711a55f735aaae22c91b996c96e",
				"rhcos-metal.x86_64.raw.gz":                    "ef9a304cba0c0050486965e38b3c6c614c0646af0b0de493c99eb6a14703cb5a",
				"rhcos-metal4k.x86_64.raw.gz":                  "9b4548b8b87322dd4d659922cddd287060d6b4ac53992d121ac32442a471f793",
				"rhcos-openstack.x86_64.qcow2.gz":              "3466690807fb710102559ea57daac0484c59ed4d914996882d601b8bb7a7ada8",
				"rhcos-ostree.x86_64.tar":                      "6708a9fbf379a2e53a9b61ae18fad2fb472bd746d4a04638dbb412195de10a24",
				"rhcos-qemu.x86_64.qcow2.gz":                   "cae8928e0cd35b88fcec7c07b1072155bde17d7dd44985f8b0d9e3862c556602",
				"rhcos-vmware.x86_64.ova":                      "6c8bfdee5930f12368b9f46a11aea736a068208262f7747f3bac54eb581531f5",
			},
		},
		{
			args: struct {
				fileName       string
				checksumFormat ChecksumFormat
			}{fileName: "testdata/broken.checksum", checksumFormat: ChecksumFormatBSD},
			want: map[string]string{
				"CentOS-Stream-Container-Base-9-20211119.0.x86_64.tar.xz":          "bd329142ec8e7455cbfb641d286cc74baaf0fac54e8b0bdbc873fc01c61bf19d",
				"CentOS-Stream-GenericCloud-9-20211119.0.x86_64.qcow2":             "84e67ec05f085bbf2fe42d3a341bfff4a4800ef1957655443638522c4c73e02c",
				"CentOS-Stream-Vagrant-9-20211119.0.x86_64.vagrant-libvirt.box":    "6e732563b0997996415cfee3f4b7aa722750134e0c0b6eeb9572a23888de6ad3",
				"CentOS-Stream-Vagrant-9-20211119.0.x86_64.vagrant-virtualbox.box": "203e2ecad207632cd6866e9971febc1265801e6cb53acf1e37592693809ea8a1",
				"CentOS-Stream-ec2-9-20211119.0.x86_64.raw.xz":                     "b5fadd02e18a1e65134cc33eb6843820d4b7be57f0531f7a243717fc8887b456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			f, err := os.Open(tt.args.fileName)
			if err != nil {
				panic(err)
			}
			got, err := Parse(f, tt.args.checksumFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			g.Expect(got).To(Equal(tt.want))
		})
	}
}
