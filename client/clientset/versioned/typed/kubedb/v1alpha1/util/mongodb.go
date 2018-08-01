package util

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	"github.com/golang/glog"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateOrPatchMongoDB(c cs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.MongoDB) *api.MongoDB) (*api.MongoDB, kutil.VerbType, error) {
	cur, err := c.MongoDBs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating MongoDB %s/%s.", meta.Namespace, meta.Name)
		out, err := c.MongoDBs(meta.Namespace).Create(transform(&api.MongoDB{
			TypeMeta: metav1.TypeMeta{
				Kind:       "MongoDB",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchMongoDB(c, cur, transform)
}

func PatchMongoDB(c cs.KubedbV1alpha1Interface, cur *api.MongoDB, transform func(*api.MongoDB) *api.MongoDB) (*api.MongoDB, kutil.VerbType, error) {
	return PatchMongoDBObject(c, cur, transform(cur.DeepCopy()))
}

func PatchMongoDBObject(c cs.KubedbV1alpha1Interface, cur, mod *api.MongoDB) (*api.MongoDB, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonmergepatch.CreateThreeWayJSONMergePatch(curJson, modJson, curJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching MongoDB %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.MongoDBs(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateMongoDB(c cs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.MongoDB) *api.MongoDB) (result *api.MongoDB, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.MongoDBs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {

			result, e2 = c.MongoDBs(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update MongoDB %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update MongoDB %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func UpdateMongoDBStatus(c cs.KubedbV1alpha1Interface, cur *api.MongoDB, transform func(*api.MongoDBStatus) *api.MongoDBStatus, useSubresource ...bool) (result *api.MongoDB, err error) {
	if len(useSubresource) > 1 {
		return nil, errors.Errorf("invalid value passed for useSubresource: %v", useSubresource)
	}

	modFunc := func() *api.MongoDB {
		return &api.MongoDB{
			TypeMeta:   cur.TypeMeta,
			ObjectMeta: cur.ObjectMeta,
			Spec:       cur.Spec,
			Status:     *transform(cur.Status.DeepCopy()),
		}
	}

	if len(useSubresource) == 1 && useSubresource[0] {
		attempt := 0
		err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
			attempt++
			var e2 error
			mod := modFunc()
			result, e2 = c.MongoDBs(cur.Namespace).UpdateStatus(mod)
			if kerr.IsNotFound(e2) {
				return false, e2
			} else if kerr.IsConflict(e2) {
				cur, _ = c.MongoDBs(cur.Namespace).Get(cur.Name, metav1.GetOptions{})
				return false, nil
			}
			return e2 == nil, nil
		})

		if err != nil {
			err = fmt.Errorf("failed to update MySQLStatus %s/%s after %d attempts due to %v", cur.Namespace, cur.Name, attempt, err)
		}
		return
	}

	mod := modFunc()
	result, _, err = PatchMongoDBObject(c, cur, mod)
	return
}
