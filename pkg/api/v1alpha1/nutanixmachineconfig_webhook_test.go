package v1alpha1_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/eks-anywhere/pkg/api/v1alpha1"
	"github.com/aws/eks-anywhere/pkg/utils/ptr"
)

func nutanixMachineConfig() *v1alpha1.NutanixMachineConfig {
	return &v1alpha1.NutanixMachineConfig{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-nmc",
		},
		Spec: v1alpha1.NutanixMachineConfigSpec{
			OSFamily:       v1alpha1.Ubuntu,
			VCPUsPerSocket: 2,
			VCPUSockets:    4,
			BootType:       v1alpha1.NutanixBootTypeLegacy,
			MemorySize:     resource.MustParse("8Gi"),
			Image: v1alpha1.NutanixResourceIdentifier{
				Type: v1alpha1.NutanixIdentifierName,
				Name: ptr.String("ubuntu-image"),
			},
			Cluster: v1alpha1.NutanixResourceIdentifier{
				Type: v1alpha1.NutanixIdentifierName,
				Name: ptr.String("cluster-1"),
			},
			Subnet: v1alpha1.NutanixResourceIdentifier{
				Type: v1alpha1.NutanixIdentifierName,
				Name: ptr.String("subnet-1"),
			},
			Project: &v1alpha1.NutanixResourceIdentifier{
				Type: v1alpha1.NutanixIdentifierName,
				Name: ptr.String("project-1"),
			},
			AdditionalCategories: []v1alpha1.NutanixCategoryIdentifier{
				{
					Key:   "category-1",
					Value: "value-1",
				},
				{
					Key:   "category-2",
					Value: "value-2",
				},
			},
			SystemDiskSize: resource.MustParse("100Gi"),
			Users: []v1alpha1.UserConfiguration{
				{
					Name:              "test-user",
					SshAuthorizedKeys: []string{"ssh AAA..."},
				},
			},
		},
	}
}

func TestValidateCreate_Valid(t *testing.T) {
	ctx := context.Background()
	g := NewWithT(t)
	config := nutanixMachineConfig()
	g.Expect(config.ValidateCreate(ctx, config)).Error().To(Succeed())
}

func TestValidateCreate_Invalid(t *testing.T) {
	ctx := context.Background()
	g := NewWithT(t)
	tests := []struct {
		name string
		fn   func(*v1alpha1.NutanixMachineConfig)
	}{
		{
			name: "invalid bootType",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.BootType = "invalid"
			},
		},
		{
			name: "invalid name",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Name = ""
			},
		},
		{
			name: "invalid os family",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.OSFamily = "invalid"
			},
		},
		{
			name: "invalid vcpus per socket",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.VCPUsPerSocket = 0
			},
		},
		{
			name: "invalid vcpus sockets",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.VCPUSockets = 0
			},
		},
		{
			name: "invalid memory size",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.MemorySize = resource.MustParse("0Gi")
			},
		},
		{
			name: "invalid image type",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Image.Type = "invalid"
			},
		},
		{
			name: "invalid image name",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Image.Name = nil
			},
		},
		{
			name: "invalid cluster type",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Cluster.Type = "invalid"
			},
		},
		{
			name: "invalid cluster name",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Cluster.Name = nil
			},
		},
		{
			name: "invalid subnet type",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Subnet.Type = "invalid"
			},
		},
		{
			name: "invalid subnet name",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Subnet.Name = nil
			},
		},
		{
			name: "invalid system disk size",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.SystemDiskSize = resource.MustParse("0Gi")
			},
		},
		{
			name: "no user",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Users = []v1alpha1.UserConfiguration{}
			},
		},
		{
			name: "no user name",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Users = []v1alpha1.UserConfiguration{
					{
						SshAuthorizedKeys: []string{"ssh AAA..."},
					},
				}
			},
		},
		{
			name: "no ssh authorized key",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Users = []v1alpha1.UserConfiguration{
					{
						Name: "eksa",
					},
				}
			},
		},
		{
			name: "invalid ssh authorized key",
			fn: func(config *v1alpha1.NutanixMachineConfig) {
				config.Spec.Users = []v1alpha1.UserConfiguration{
					{
						Name:              "eksa",
						SshAuthorizedKeys: []string{""},
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := nutanixMachineConfig()
			tt.fn(config)
			_, err := config.ValidateCreate(ctx, config)
			g.Expect(err).To(HaveOccurred(), "expected error for %s", tt.name)
		})
	}
}

func TestNutanixMachineConfigWebhooksValidateUpdateReconcilePaused(t *testing.T) {
	ctx := context.Background()
	g := NewWithT(t)
	oldConfig := nutanixMachineConfig()
	newConfig := nutanixMachineConfig()
	newConfig.Spec.Cluster.Name = ptr.String("new-cluster")
	oldConfig.Annotations = map[string]string{
		"anywhere.eks.amazonaws.com/paused": "true",
	}
	g.Expect(newConfig.ValidateUpdate(ctx, newConfig, oldConfig)).Error().To(Succeed())
}

func TestValidateUpdate(t *testing.T) {
	ctx := context.Background()
	g := NewWithT(t)
	oldConfig := nutanixMachineConfig()
	newConfig := nutanixMachineConfig()
	newConfig.Spec.VCPUSockets = 8
	g.Expect(newConfig.ValidateUpdate(ctx, newConfig, oldConfig)).Error().To(Succeed())

	oldConfig = nutanixMachineConfig()
	oldConfig.SetManagedBy("mgmt-cluster")
	g.Expect(newConfig.ValidateUpdate(ctx, newConfig, oldConfig)).Error().To(Succeed())

	oldConfig = nutanixMachineConfig()
	oldConfig.SetControlPlane()
	g.Expect(newConfig.ValidateUpdate(ctx, newConfig, oldConfig)).Error().To(HaveOccurred())
}

func TestValidateUpdate_Invalid(t *testing.T) {
	ctx := context.Background()
	g := NewWithT(t)
	tests := []struct {
		name string
		fn   func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig)
	}{
		{
			name: "different os family",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				new.Spec.OSFamily = v1alpha1.Bottlerocket
			},
		},
		{
			name: "different bootType",
			fn: func(newConfig *v1alpha1.NutanixMachineConfig, oldConfig *v1alpha1.NutanixMachineConfig) {
				oldConfig.SetControlPlane()
				newConfig.Spec.BootType = v1alpha1.NutanixBootTypeUEFI
			},
		},
		{
			name: "different cluster",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				new.Spec.Cluster = v1alpha1.NutanixResourceIdentifier{
					Type: v1alpha1.NutanixIdentifierName,
					Name: ptr.String("cluster-2"),
				}
			},
		},
		{
			name: "different subnet",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				new.Spec.Subnet = v1alpha1.NutanixResourceIdentifier{
					Type: v1alpha1.NutanixIdentifierName,
					Name: ptr.String("subnet-2"),
				}
			},
		},
		{
			name: "old cluster is managed",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				new.Spec.OSFamily = v1alpha1.Bottlerocket
				old.SetManagedBy("test")
			},
		},
		{
			name: "mismatch vcpu sockets on control plane cluster",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				old.SetControlPlane()
				new.Spec.VCPUSockets++
			},
		},
		{
			name: "mismatch vcpu per socket on control plane cluster",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				old.SetControlPlane()
				new.Spec.VCPUsPerSocket++
			},
		},
		{
			name: "mismatch memory size on control plane cluster",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				old.SetControlPlane()
				new.Spec.MemorySize.Add(resource.MustParse("1Gi"))
			},
		},
		{
			name: "mismatch system disk size on control plane cluster",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				old.SetControlPlane()
				new.Spec.SystemDiskSize.Add(resource.MustParse("1Gi"))
			},
		},
		{
			name: "mismatch users on control plane cluster",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				old.SetControlPlane()
				new.Spec.Users = append(new.Spec.Users, v1alpha1.UserConfiguration{
					Name: "another-user",
				})
			},
		},
		{
			name: "invalid bootType for control plane cluster",
			fn: func(newConfig *v1alpha1.NutanixMachineConfig, oldConfig *v1alpha1.NutanixMachineConfig) {
				oldConfig.SetControlPlane()
				newConfig.Spec.BootType = "invalid"
			},
		},
		{
			name: "invalid vcpus per socket",
			fn: func(new *v1alpha1.NutanixMachineConfig, old *v1alpha1.NutanixMachineConfig) {
				new.Spec.VCPUsPerSocket = 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldConfig := nutanixMachineConfig()
			newConfig := nutanixMachineConfig()
			tt.fn(newConfig, oldConfig)
			_, err := newConfig.ValidateUpdate(ctx, newConfig, oldConfig)
			g.Expect(err).To(HaveOccurred(), "expected error for %s", tt.name)
		})
	}
}

func TestValidateUpdate_OldObjectNotMachineConfig(t *testing.T) {
	ctx := context.Background()
	g := NewWithT(t)
	oldConfig := nutanixDatacenterConfig()
	newConfig := nutanixMachineConfig()
	_, err := newConfig.ValidateUpdate(ctx, newConfig, oldConfig)
	g.Expect(err).To(HaveOccurred())
}

func TestNutanixMachineConfigWebhooksValidateDelete(t *testing.T) {
	ctx := context.Background()
	g := NewWithT(t)
	config := nutanixMachineConfig()
	g.Expect(config.ValidateDelete(ctx, config)).Error().To(Succeed())
}

func TestNutanixMachineConfigValidateCreateCastFail(t *testing.T) {
	g := NewWithT(t)

	// Create a different type that will cause the cast to fail
	wrongType := &v1alpha1.Cluster{}

	// Create the config object that implements CustomValidator
	config := &v1alpha1.NutanixMachineConfig{}

	// Call ValidateCreate with the wrong type
	warnings, err := config.ValidateCreate(context.TODO(), wrongType)

	// Verify that an error is returned
	g.Expect(warnings).To(BeNil())
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("expected a NutanixMachineConfig"))
}

func TestNutanixMachineConfigValidateUpdateCastFail(t *testing.T) {
	g := NewWithT(t)

	// Create a different type that will cause the cast to fail
	wrongType := &v1alpha1.Cluster{}

	// Create the config object that implements CustomValidator
	config := &v1alpha1.NutanixMachineConfig{}

	// Call ValidateUpdate with the wrong type
	warnings, err := config.ValidateUpdate(context.TODO(), wrongType, &v1alpha1.NutanixMachineConfig{})

	// Verify that an error is returned
	g.Expect(warnings).To(BeNil())
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("expected a NutanixMachineConfig"))
}

func TestNutanixMachineConfigValidateDeleteCastFail(t *testing.T) {
	g := NewWithT(t)

	// Create a different type that will cause the cast to fail
	wrongType := &v1alpha1.Cluster{}

	// Create the config object that implements CustomValidator
	config := &v1alpha1.NutanixMachineConfig{}

	// Call ValidateDelete with the wrong type
	warnings, err := config.ValidateDelete(context.TODO(), wrongType)

	// Verify that an error is returned
	g.Expect(warnings).To(BeNil())
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("expected a NutanixMachineConfig"))
}
