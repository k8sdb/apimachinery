package monitor

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	_ "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	tapi "github.com/k8sdb/apimachinery/api"
	kapi "k8s.io/kubernetes/pkg/api"
	kerr "k8s.io/kubernetes/pkg/api/errors"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/util/intstr"
)

const (
	k8sdbExporter         = "k8sdb-exporter"
	ImageK8sdbExporter    = "k8sdb/exporter"
	exporterContainerName = "exporter"
	postgresPortName      = "pg_port"
)

var exporterLabel = map[string]string{
	"k8sdb/exporter": "appscode/exporter",
}

type PrometheusController struct {
	monitoringClient v1alpha1.MonitoringV1alpha1Interface
	kubeCLient       clientset.Interface
}

func NewPrometheusController(KubeCLient clientset.Interface, MonitoringClient v1alpha1.MonitoringV1alpha1Interface) Monitor {
	return &PrometheusController{
		monitoringClient: MonitoringClient,
		kubeCLient:       KubeCLient,
	}
}

func (c *PrometheusController) AddMonitor(meta *kapi.ObjectMeta, spec *tapi.MonitorSpec) error {
	err := c.ensureExporter(meta)
	if err != nil {
		return err
	}
	if ok, err := c.supportPrometheusOperator(); err != nil {
		return err
	} else if !ok {
		return nil
	}
	return c.ensureMonitor(meta, spec)
}

func (c *PrometheusController) UpdateMonitor(meta *kapi.ObjectMeta, spec *tapi.MonitorSpec) error {
	err := c.ensureExporter(meta)
	if err != nil {
		return err
	}
	if ok, err := c.supportPrometheusOperator(); err != nil {
		return err
	} else if !ok {
		return nil
	}
	return c.ensureMonitor(meta, spec)
}

func (c *PrometheusController) DeleteMonitor(meta *kapi.ObjectMeta, spec *tapi.MonitorSpec) error {
	if ok, err := c.supportPrometheusOperator(); err != nil {
		return err
	} else if !ok {
		return nil
	}
	// Delete a service monitor for this DB TPR, if does not exist
	return nil
}

func (c *PrometheusController) ensureExporter(meta *kapi.ObjectMeta) error {
	// check if the global exporter is running or not
	// if not running, create a deployment of exporter pod
	exporter, err := c.kubeCLient.Extensions().Deployments(namespace()).Get(k8sdbExporter)
	if kerr.IsNotFound(err) {
		//create exporter
		if _, err = c.createk8sdbExporter(); err != nil {
			return err
		}
		if _, err = c.createServiceForExporter(); err != nil {
			return err
		}

	}
	return err
}

func (c *PrometheusController) supportPrometheusOperator() (bool, error) {
	// check if the prometheus TPR group exists
	return false, nil
}

func (c *PrometheusController) ensureMonitor(meta *kapi.ObjectMeta, spec *tapi.MonitorSpec) error {
	// Check if a service monitor exists,
	// If does not exist, create one.
	// If exists, then update it only if update is needed.
	return nil
}

// namespace returns the namespace of tiller
func namespace() string {
	if ns := os.Getenv("KUBEDB_OPERATOR_NAMESPACE"); ns != "" {
		return ns
	}

	// Fall back to the namespace associated with the service account token, if available
	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}

	return kapi.NamespaceDefault
}

func (c *PrometheusController) createk8sdbExporter() (*extensions.Deployment, error) {
	exporter := &extensions.Deployment{
		ObjectMeta: kapi.ObjectMeta{
			Name:      k8sdbExporter,
			Namespace: namespace(),
			Labels:    exporterLabel,
		},
		Spec: extensions.DeploymentSpec{
			Replicas: 1,
			Template: kapi.PodTemplate{
				Template: kapi.PodTemplateSpec{
					Spec: kapi.PodSpec{
						Containers: []kapi.Container{
							{
								Name:            exporterContainerName,
								Image:           ImageK8sdbExporter,
								ImagePullPolicy: kapi.PullIfNotPresent,
								Ports: []kapi.ContainerPort{
									{ // ports for postgres
										Name:          postgresPortName,
										Protocol:      kapi.ProtocolTCP,
										ContainerPort: 5432,
									},
									// TODO add port for elastic search
								},
							},
						},
					},
				},
			},
		},
	}
	return c.kubeCLient.Extensions().Deployments(namespace()).Create(exporter)
}

func (c *PrometheusController) createServiceForExporter() (*kapi.Service, error) {
	svc := &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name:      k8sdbExporter,
			Namespace: namespace(),
			Labels:    exporterLabel,
		},
		Spec: kapi.ServiceSpec{
			Type: kapi.ServiceTypeClusterIP,
			Ports: []kapi.ServicePort{
				{
					Name:       postgresPortName,
					Port:       9187,
					Protocol:   kapi.ProtocolTCP,
					TargetPort: intstr.FromString(postgresPortName),
				},
				// TODO Add ports for elasticsearch
			},
		},
	}
	return c.kubeCLient.Core().Services(namespace()).Create(svc)
}
