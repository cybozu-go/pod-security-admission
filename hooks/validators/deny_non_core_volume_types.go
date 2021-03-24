package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// DenyNonCoreVolumeTypes is a Validator that denies usage of non-core volume types
func DenyNonCoreVolumeTypes(ctx context.Context, pod *corev1.Pod) admission.Response {
	for _, vol := range pod.Spec.Volumes {
		if vol.HostPath != nil {
			return admission.Denied("Volume type HostPath is not allowed to be used")
		}
		if vol.GCEPersistentDisk != nil {
			return admission.Denied("Volume type GCEPersistentDisk is not allowed to be used")
		}
		if vol.AWSElasticBlockStore != nil {
			return admission.Denied("Volume type AWSElasticBlockStore is not allowed to be used")
		}
		if vol.GitRepo != nil {
			return admission.Denied("Volume type GitRepo is not allowed to be used")
		}
		if vol.NFS != nil {
			return admission.Denied("Volume type NFS is not allowed to be used")
		}
		if vol.ISCSI != nil {
			return admission.Denied("Volume type ISCSI is not allowed to be used")
		}
		if vol.Glusterfs != nil {
			return admission.Denied("Volume type Glusterfs is not allowed to be used")
		}
		if vol.RBD != nil {
			return admission.Denied("Volume type RBD is not allowed to be used")
		}
		if vol.FlexVolume != nil {
			return admission.Denied("Volume type FlexVolume is not allowed to be used")
		}
		if vol.Cinder != nil {
			return admission.Denied("Volume type Cinder is not allowed to be used")
		}
		if vol.CephFS != nil {
			return admission.Denied("Volume type CephFS is not allowed to be used")
		}
		if vol.Flocker != nil {
			return admission.Denied("Volume type Flocker is not allowed to be used")
		}
		if vol.FC != nil {
			return admission.Denied("Volume type FC is not allowed to be used")
		}
		if vol.AzureFile != nil {
			return admission.Denied("Volume type AzureFile is not allowed to be used")
		}
		if vol.VsphereVolume != nil {
			return admission.Denied("Volume type VsphereVolume is not allowed to be used")
		}
		if vol.Quobyte != nil {
			return admission.Denied("Volume type Quobyte is not allowed to be used")
		}
		if vol.AzureDisk != nil {
			return admission.Denied("Volume type AzureDisk is not allowed to be used")
		}
		if vol.PortworxVolume != nil {
			return admission.Denied("Volume type PortworkVolume is not allowed to be used")
		}
		if vol.ScaleIO != nil {
			return admission.Denied("Volume type ScaleIO is not allowed to be used")
		}
		if vol.StorageOS != nil {
			return admission.Denied("Volume type StorageOS is not allowed to be used")
		}
		if vol.CSI != nil {
			return admission.Denied("Volume type CSI is not allowed to be used")
		}
	}
	return admission.Allowed("ok")
}
