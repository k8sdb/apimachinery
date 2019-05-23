package controller

import (
	"math"
	"path/filepath"

	"github.com/appscode/go/log"
	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/azure"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/reference"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/objectstore-api/osm"
)

const UtilVolumeName = "util-volume"

func (c *Controller) DeleteSnapshotData(snapshot *api.Snapshot) error {
	cfg, err := osm.NewOSMContext(c.Client, snapshot.Spec.Backend, snapshot.Namespace)
	if err != nil {
		return err
	}

	loc, err := stow.Dial(cfg.Provider, cfg.Config)
	if err != nil {
		return err
	}
	bucket, err := snapshot.Spec.Backend.Container()
	if err != nil {
		return err
	}
	container, err := loc.Container(bucket)
	if err != nil {
		return err
	}

	prefixLocation, _ := snapshot.Location() // error checked by .Container()
	prefix := filepath.Join(prefixLocation, snapshot.Name)
	prefix += "/" // A separator after prefix to prevent multiple snapshot's prefix matching. ref: https://github.com/kubedb/project/issues/377
	cursor := stow.CursorStart
	for {
		items, next, err := container.Items(prefix, cursor, 50)
		if err != nil {
			return err
		}
		for _, item := range items {
			if err := container.RemoveItem(item.ID()); err != nil {
				return err
			}
		}
		cursor = next
		if stow.IsCursorEnd(cursor) {
			break
		}
	}

	return nil
}

func (c *Controller) checkGoverningService(name, namespace string) (bool, error) {
	_, err := c.Client.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (c *Controller) CreateGoverningService(name, namespace string) error {
	// Check if service name exists
	found, err := c.checkGoverningService(name, namespace)
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	service := &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: core.ServiceSpec{
			Type:      core.ServiceTypeClusterIP,
			ClusterIP: core.ClusterIPNone,
		},
	}
	_, err = c.Client.CoreV1().Services(namespace).Create(service)
	return err
}

func (c *Controller) SetJobOwnerReference(snapshot *api.Snapshot, job *batch.Job) error {
	secret, err := c.Client.CoreV1().Secrets(snapshot.Namespace).Get(snapshot.OSMSecretName(), metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	} else {
		_, _, err := core_util.PatchSecret(c.Client, secret, func(in *core.Secret) *core.Secret {
			in.SetOwnerReferences([]metav1.OwnerReference{
				{
					APIVersion: batch.SchemeGroupVersion.String(),
					Kind:       "Job",
					Name:       job.Name,
					UID:        job.UID,
				},
			})
			return in
		})
		if err != nil {
			return err
		}
	}

	pvc, err := c.Client.CoreV1().PersistentVolumeClaims(snapshot.Namespace).Get(job.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	} else {
		_, _, err := core_util.PatchPVC(c.Client, pvc, func(in *core.PersistentVolumeClaim) *core.PersistentVolumeClaim {
			in.SetOwnerReferences([]metav1.OwnerReference{
				{
					APIVersion: batch.SchemeGroupVersion.String(),
					Kind:       "Job",
					Name:       job.Name,
					UID:        job.UID,
				},
			})
			return in
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetVolumeForSnapshot returns pvc or empty directory depending on StorageType.
// In case of PVC, this function will create a PVC then returns the volume.
func (c *Controller) GetVolumeForSnapshot(st api.StorageType, pvcSpec *core.PersistentVolumeClaimSpec, jobName, namespace string) (*core.Volume, error) {
	if st == api.StorageTypeEphemeral {
		ed := core.EmptyDirVolumeSource{}
		if pvcSpec != nil {
			if sz, found := pvcSpec.Resources.Requests[core.ResourceStorage]; found {
				ed.SizeLimit = &sz
			}
		}
		return &core.Volume{
			Name: UtilVolumeName,
			VolumeSource: core.VolumeSource{
				EmptyDir: &ed,
			},
		}, nil
	}

	volume := &core.Volume{
		Name: UtilVolumeName,
	}
	if len(pvcSpec.AccessModes) == 0 {
		pvcSpec.AccessModes = []core.PersistentVolumeAccessMode{
			core.ReadWriteOnce,
		}
		log.Infof(`Using "%v" as AccessModes in "%v"`, core.ReadWriteOnce, *pvcSpec)
	}

	claim := &core.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: *pvcSpec,
	}
	if pvcSpec.StorageClassName != nil {
		claim.Annotations = map[string]string{
			"volume.beta.kubernetes.io/storage-class": *pvcSpec.StorageClassName,
		}
	}

	if _, err := c.Client.CoreV1().PersistentVolumeClaims(claim.Namespace).Create(claim); err != nil {
		return nil, err
	}

	volume.PersistentVolumeClaim = &core.PersistentVolumeClaimVolumeSource{
		ClaimName: claim.Name,
	}

	return volume, nil
}

func (c *Controller) CreateStatefulSetPodDisruptionBudget(sts *appsv1.StatefulSet) error {
	_, err := c.Client.PolicyV1beta1().PodDisruptionBudgets(sts.Namespace).Get(sts.Name, metav1.GetOptions{})
	if !kerr.IsNotFound(err) {
		return err
	}

	maxUnavailable := int32(math.Ceil((float64(*sts.Spec.Replicas) - 1.0) / 2.0))
	pdb := policyv1beta1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sts.Name,
			Namespace: sts.Namespace,
			Labels:    sts.Labels,
		},

		Spec: policyv1beta1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: sts.Spec.Template.Labels,
			},
			MaxUnavailable: &intstr.IntOrString{IntVal: maxUnavailable},
		},
	}

	ref, rerr := reference.GetReference(clientsetscheme.Scheme, sts)
	if rerr != nil {
		return rerr
	}
	core_util.EnsureOwnerReference(&pdb.ObjectMeta, ref)

	//Create PDB
	_, err = c.Client.PolicyV1beta1().PodDisruptionBudgets(sts.Namespace).Create(&pdb)
	return err
}

func (c *Controller) CreateDeploymentPodDisruptionBudget(deployment *appsv1.Deployment) error {
	_, err := c.Client.PolicyV1beta1().PodDisruptionBudgets(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{})
	if !kerr.IsNotFound(err) {
		return err
	}
	deploymentReplicas := float64(*(deployment.Spec.Replicas))
	maxUnavailable := int32(math.Ceil((deploymentReplicas - 1.0) / 2.0))
	pdb := policyv1beta1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Labels:    deployment.Labels,
		},

		Spec: policyv1beta1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: deployment.Spec.Template.Labels,
			},
			MaxUnavailable: &intstr.IntOrString{IntVal: maxUnavailable},
		},
	}

	ref, rerr := reference.GetReference(clientsetscheme.Scheme, deployment)
	if rerr != nil {
		return rerr
	}
	core_util.EnsureOwnerReference(&pdb.ObjectMeta, ref)

	//Create PDB
	_, err = c.Client.PolicyV1beta1().PodDisruptionBudgets(deployment.Namespace).Create(&pdb)
	return err
}
