package validators

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// DenyNonCoreVolumeTypes is a Validator that denies usage of non-core volume types
type DenyNonCoreVolumeTypes struct{}

func (v DenyNonCoreVolumeTypes) Validate(ctx context.Context, pod *corev1.Pod) field.ErrorList {
	p := field.NewPath("spec").Child("volumes")
	var errs field.ErrorList

	for i, vol := range pod.Spec.Volumes {
		if vol.GCEPersistentDisk != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type GCEPersistentDisk is not allowed to be used"))
		}
		if vol.AWSElasticBlockStore != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type AWSElasticBlockStore is not allowed to be used"))
		}
		if vol.GitRepo != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type GitRepo is not allowed to be used"))
		}
		if vol.NFS != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type NFS is not allowed to be used"))
		}
		if vol.ISCSI != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type ISCSI is not allowed to be used"))
		}
		if vol.Glusterfs != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type Glusterfs is not allowed to be used"))
		}
		if vol.RBD != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type RBD is not allowed to be used"))
		}
		if vol.FlexVolume != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type FlexVolume is not allowed to be used"))
		}
		if vol.Cinder != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type Cinder is not allowed to be used"))
		}
		if vol.CephFS != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type CephFS is not allowed to be used"))
		}
		if vol.Flocker != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type Flocker is not allowed to be used"))
		}
		if vol.FC != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type FC is not allowed to be used"))
		}
		if vol.AzureFile != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type AzureFile is not allowed to be used"))
		}
		if vol.VsphereVolume != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type VsphereVolume is not allowed to be used"))
		}
		if vol.Quobyte != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type Quobyte is not allowed to be used"))
		}
		if vol.AzureDisk != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type AzureDisk is not allowed to be used"))
		}
		if vol.PortworxVolume != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type PortworxVolume is not allowed to be used"))
		}
		if vol.ScaleIO != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type ScaleIO is not allowed to be used"))
		}
		if vol.StorageOS != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type StorageOS is not allowed to be used"))
		}
		if vol.CSI != nil {
			errs = append(errs, field.Forbidden(p.Index(i), "Volume type CSI is not allowed to be used"))
		}
	}
	return errs
}
