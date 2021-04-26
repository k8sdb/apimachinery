/*
Copyright AppsCode Inc. and Contributors

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

package v1alpha2

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

//const (
//	RedisShardAffinityTemplateVar = "SHARD_INDEX"
//)

func (rs RedisSentinel) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedisSentinel))
}

var _ apis.ResourceInfo = &RedisSentinel{}

func (rs RedisSentinel) OffshootName() string {
	return rs.Name
}

func (rs RedisSentinel) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      rs.ResourceFQN(),
		meta_util.InstanceLabelKey:  rs.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (rs RedisSentinel) OffshootLabels() map[string]string {
	out := rs.OffshootSelectors()
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, out, rs.Labels)
}

func (rs RedisSentinel) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralRedisSentinel, kubedb.GroupName)
}

func (rs RedisSentinel) ResourceShortCode() string {
	return ResourceCodeRedisSentinel
}

func (rs RedisSentinel) ResourceKind() string {
	return ResourceKindRedisSentinel
}

func (rs RedisSentinel) ResourceSingular() string {
	return ResourceSingularRedisSentinel
}

func (rs RedisSentinel) ResourcePlural() string {
	return ResourcePluralRedisSentinel
}

func (rs RedisSentinel) ServiceName() string {
	return rs.OffshootName()
}

func (rs RedisSentinel) GoverningServiceName() string {
	return meta_util.NameWithSuffix(rs.ServiceName(), "pods")
}

func (rs RedisSentinel) ConfigSecretName() string {
	return rs.OffshootName()
}

//func (rs RedisSentinel) BaseNameForShard() string {
//	return fmt.Sprintf("%s-shard", rs.OffshootName())
//}

//func (r Redis) StatefulSetNameWithShard(i int) string {
//	return fmt.Sprintf("%s%d", r.BaseNameForShard(), i)
//}

//func (rs RedisSentinel) Address() string {
//	return fmt.Sprintf("%v.%v.svc:%d", rs.Name, rs.Namespace, RedisDatabasePort)
//}

type redisSentinelApp struct {
	*RedisSentinel
}

func (rs redisSentinelApp) Name() string {
	return rs.RedisSentinel.Name
}

func (rs redisSentinelApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularRedisSentinel))
}

func (rs RedisSentinel) AppBindingMeta() appcat.AppBindingMeta {
	return &redisSentinelApp{&rs}
}

type redisSentinelStatsService struct {
	*RedisSentinel
}

func (rs redisSentinelStatsService) GetNamespace() string {
	return rs.RedisSentinel.GetNamespace()
}

func (rs redisSentinelStatsService) ServiceName() string {
	return rs.OffshootName() + "-stats"
}

func (rs redisSentinelStatsService) ServiceMonitorName() string {
	return rs.ServiceName()
}

//func (p redisStatsService) ServiceMonitorAdditionalLabels() map[string]string {
//	return p.OffshootLabels()
//}

//func (r redisStatsService) Path() string {
//	return DefaultStatsPath
//}
//
//func (r redisStatsService) Scheme() string {
//	return ""
//}
//
//func (r Redis) StatsService() mona.StatsAccessor {
//	return &redisStatsService{&r}
//}
//
//func (r Redis) StatsServiceLabels() map[string]string {
//	lbl := meta_util.FilterKeys(kubedb.GroupName, r.OffshootSelectors(), r.Labels)
//	lbl[LabelRole] = RoleStats
//	return lbl
//}
//
//func (r *Redis) SetDefaults(topology *core_util.Topology) {
//	if r == nil {
//		return
//	}
//
//	// perform defaulting
//	if r.Spec.Mode == "" {
//		r.Spec.Mode = RedisModeStandalone
//	} else if r.Spec.Mode == RedisModeCluster {
//		if r.Spec.Cluster == nil {
//			r.Spec.Cluster = &RedisClusterSpec{}
//		}
//		if r.Spec.Cluster.Master == nil {
//			r.Spec.Cluster.Master = pointer.Int32P(3)
//		}
//		if r.Spec.Cluster.Replicas == nil {
//			r.Spec.Cluster.Replicas = pointer.Int32P(1)
//		}
//	}
//	if r.Spec.StorageType == "" {
//		r.Spec.StorageType = StorageTypeDurable
//	}
//	if r.Spec.TerminationPolicy == "" {
//		r.Spec.TerminationPolicy = TerminationPolicyDelete
//	}
//
//	if r.Spec.PodTemplate.Spec.ServiceAccountName == "" {
//		r.Spec.PodTemplate.Spec.ServiceAccountName = r.OffshootName()
//	}
//
//	labels := r.OffshootSelectors()
//	if r.Spec.Mode == RedisModeCluster {
//		labels[RedisShardKey] = r.ShardNodeTemplate()
//	}
//	r.setDefaultAffinity(&r.Spec.PodTemplate, labels, topology)
//
//	r.Spec.Monitor.SetDefaults()
//
//	r.SetTLSDefaults()
//	SetDefaultResourceLimits(&r.Spec.PodTemplate.Spec.Resources, DefaultResourceLimits)
//}
//
//func (r *Redis) SetTLSDefaults() {
//	if r.Spec.TLS == nil || r.Spec.TLS.IssuerRef == nil {
//		return
//	}
//	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(RedisServerCert), r.CertificateName(RedisServerCert))
//	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(RedisClientCert), r.CertificateName(RedisClientCert))
//	r.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(r.Spec.TLS.Certificates, string(RedisMetricsExporterCert), r.CertificateName(RedisMetricsExporterCert))
//}
//
//func (r *RedisSpec) GetPersistentSecrets() []string {
//	return nil
//}
//
//func (r *Redis) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
//	if podTemplate == nil {
//		return
//	} else if podTemplate.Spec.Affinity != nil {
//		topology.ConvertAffinity(podTemplate.Spec.Affinity)
//		return
//	}
//
//	podTemplate.Spec.Affinity = &corev1.Affinity{
//		PodAntiAffinity: &corev1.PodAntiAffinity{
//			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
//				// Prefer to not schedule multiple pods on the same node
//				{
//					Weight: 100,
//					PodAffinityTerm: corev1.PodAffinityTerm{
//						Namespaces: []string{r.Namespace},
//						LabelSelector: &metav1.LabelSelector{
//							MatchLabels: labels,
//						},
//
//						TopologyKey: corev1.LabelHostname,
//					},
//				},
//				// Prefer to not schedule multiple pods on the node with same zone
//				{
//					Weight: 50,
//					PodAffinityTerm: corev1.PodAffinityTerm{
//						Namespaces: []string{r.Namespace},
//						LabelSelector: &metav1.LabelSelector{
//							MatchLabels: labels,
//						},
//						TopologyKey: topology.LabelZone,
//					},
//				},
//			},
//		},
//	}
//}
//
//func (r Redis) ShardNodeTemplate() string {
//	if r.Spec.Mode == RedisModeStandalone {
//		panic("shard template is not applicable to a standalone redis server")
//	}
//	return fmt.Sprintf("${%s}", RedisShardAffinityTemplateVar)
//}
//
//// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
//func (r *Redis) CertificateName(alias RedisCertificateAlias) string {
//	return meta_util.NameWithSuffix(r.Name, fmt.Sprintf("%s-cert", string(alias)))
//}
//
//// MustCertSecretName returns the secret name for a certificate alias
//func (r *Redis) MustCertSecretName(alias RedisCertificateAlias) string {
//	if r == nil {
//		panic("missing Redis database")
//	} else if r.Spec.TLS == nil {
//		panic(fmt.Errorf("Redis %s/%s is missing tls spec", r.Namespace, r.Name))
//	}
//	name, ok := kmapi.GetCertificateSecretName(r.Spec.TLS.Certificates, string(alias))
//	if !ok {
//		panic(fmt.Errorf("Redis %s/%s is missing secret name for %s certificate", r.Namespace, r.Name, alias))
//	}
//	return name
//}
//
//func (r *Redis) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
//	// Desire number of statefulSets
//	expectedItems := 1
//	if r.Spec.Cluster != nil {
//		expectedItems = int(pointer.Int32(r.Spec.Cluster.Master))
//	}
//	return checkReplicas(lister.StatefulSets(r.Namespace), labels.SelectorFromSet(r.OffshootLabels()), expectedItems)
//}
