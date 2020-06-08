/*
Copyright The KubeDB Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"context"
	"encoding/json"
	"fmt"

	api "kubedb.dev/apimachinery/apis/ops/v1alpha1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned/typed/ops/v1alpha1"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	kutil "kmodules.xyz/client-go"
)

func CreateOrPatchProxySQLOpsRequest(ctx context.Context, c cs.OpsV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.ProxySQLOpsRequest) *api.ProxySQLOpsRequest, opts metav1.PatchOptions) (*api.ProxySQLOpsRequest, kutil.VerbType, error) {
	cur, err := c.ProxySQLOpsRequests(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating ProxySQLOpsRequest %s/%s.", meta.Namespace, meta.Name)
		out, err := c.ProxySQLOpsRequests(meta.Namespace).Create(ctx, transform(&api.ProxySQLOpsRequest{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ProxySQLOpsRequest",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}), metav1.CreateOptions{
			DryRun:       opts.DryRun,
			FieldManager: opts.FieldManager,
		})
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchProxySQLOpsRequest(ctx, c, cur, transform, opts)
}

func PatchProxySQLOpsRequest(ctx context.Context, c cs.OpsV1alpha1Interface, cur *api.ProxySQLOpsRequest, transform func(*api.ProxySQLOpsRequest) *api.ProxySQLOpsRequest, opts metav1.PatchOptions) (*api.ProxySQLOpsRequest, kutil.VerbType, error) {
	return PatchProxySQLOpsRequestObject(ctx, c, cur, transform(cur.DeepCopy()), opts)
}

func PatchProxySQLOpsRequestObject(ctx context.Context, c cs.OpsV1alpha1Interface, cur, mod *api.ProxySQLOpsRequest, opts metav1.PatchOptions) (*api.ProxySQLOpsRequest, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching ProxySQLOpsRequest %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.ProxySQLOpsRequests(cur.Namespace).Patch(ctx, cur.Name, types.MergePatchType, patch, opts)
	return out, kutil.VerbPatched, err
}

func TryUpdateProxySQLOpsRequest(ctx context.Context, c cs.OpsV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.ProxySQLOpsRequest) *api.ProxySQLOpsRequest, opts metav1.UpdateOptions) (result *api.ProxySQLOpsRequest, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.ProxySQLOpsRequests(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.ProxySQLOpsRequests(cur.Namespace).Update(ctx, transform(cur.DeepCopy()), opts)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update ProxySQLOpsRequest %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update ProxySQLOpsRequest %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func UpdateProxySQLOpsRequestStatus(
	ctx context.Context,
	c cs.OpsV1alpha1Interface,
	meta metav1.ObjectMeta,
	transform func(*api.ProxySQLOpsRequestStatus) *api.ProxySQLOpsRequestStatus,
	opts metav1.UpdateOptions,
) (result *api.ProxySQLOpsRequest, err error) {
	apply := func(x *api.ProxySQLOpsRequest) *api.ProxySQLOpsRequest {
		return &api.ProxySQLOpsRequest{
			TypeMeta:   x.TypeMeta,
			ObjectMeta: x.ObjectMeta,
			Spec:       x.Spec,
			Status:     *transform(x.Status.DeepCopy()),
		}
	}

	attempt := 0
	cur, err := c.ProxySQLOpsRequests(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		var e2 error
		result, e2 = c.ProxySQLOpsRequests(meta.Namespace).UpdateStatus(ctx, apply(cur), opts)
		if kerr.IsConflict(e2) {
			latest, e3 := c.ProxySQLOpsRequests(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
			switch {
			case e3 == nil:
				cur = latest
				return false, nil
			case kutil.IsRequestRetryable(e3):
				return false, nil
			default:
				return false, e3
			}
		} else if err != nil && !kutil.IsRequestRetryable(e2) {
			return false, e2
		}
		return e2 == nil, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update status of ProxySQLOpsRequest %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}