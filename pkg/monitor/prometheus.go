package monitor

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	_ "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	tapi "github.com/k8sdb/apimachinery/api"
	"k8s.io/client-go/pkg/api/v1"
	kapi "k8s.io/kubernetes/pkg/api"
	kerr "k8s.io/kubernetes/pkg/api/errors"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/util/intstr"
)

const (
	k8sdbExporter         = "k8sdb-exporter"
	ImageK8sdbExporter    = "k8sdb/exporter"
	exporterContainerName = "k8sdbExporter"
	exporterPort          = "k8sdbExporter"
)

var exporterLabel = map[string]string{
	"run": "k8sdbExporter",
}

type PrometheusController struct {
	monitoringClient  v1alpha1.MonitoringV1alpha1Client
	kubeCLient        clientset.Interface
	exporterNamespace string
}

func NewPrometheusController(KubeCLient clientset.Interface, MonitoringClient v1alpha1.MonitoringV1alpha1Client) Monitor {
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
		return errors.New("This cluster lacks prometheus operator support")
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
		return errors.New("This cluster lacks prometheus operator support")
	}
	// Delete a service monitor for this DB TPR, if does not exist
	return c.monitoringClient.ServiceMonitors(spec.Prometheus.Namespace).Delete(getServiceMonitorName(meta), nil)
}

func (c *PrometheusController) ensureExporter(meta *kapi.ObjectMeta) error {
	// check if the global exporter is running or not
	// if not running, create a deployment of exporter pod
	_, err := c.kubeCLient.Extensions().Deployments(c.exporterNamespace).Get(k8sdbExporter)
	if kerr.IsNotFound(err) {
		//create exporter
		if _, err = c.createK8sdbExporter(); err != nil {
			return err
		}
		if err = c.createServiceForExporter(); err != nil {
			return err
		}

	}
	return err
}

func (c *PrometheusController) supportPrometheusOperator() (bool, error) {
	// check if the prometheus TPR group exists
	_, err := c.kubeCLient.Extensions().ThirdPartyResources().Get("prometheus." + prom.TPRGroup)
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *PrometheusController) ensureMonitor(meta *kapi.ObjectMeta, spec *tapi.MonitorSpec) error {
	// Check if a service monitor exists,
	// If does not exist, create one.
	// If exists, then update it only if update is needed.
	serviceMonitorName := getServiceMonitorName(meta)
	svcMonitor, err := c.monitoringClient.ServiceMonitors(spec.Prometheus.Namespace).Get(serviceMonitorName)
	if err != nil && !kerr.IsNotFound(err) {
		return err
	}
	if err == nil {
		return c.checkServiceMonitorUpdate(svcMonitor, spec)
	}
	//else create service monitor
	serviceMonitor := &prom.ServiceMonitor{
		ObjectMeta: v1.ObjectMeta{
			Name:      serviceMonitorName,
			Namespace: spec.Prometheus.Namespace,
			Labels:    spec.Prometheus.Labels,
		},
		Spec: prom.ServiceMonitorSpec{
			NamespaceSelector: prom.NamespaceSelector{
				MatchNames: []string{meta.Namespace},
			},
			Endpoints: []prom.Endpoint{
				{
					Port:     exporterPort,
					Interval: spec.Prometheus.Interval,
					Path:     fmt.Sprintf("/k8sdb.com/v1beta1/namespaces/:%s/:%s/:%s/metrics", meta.Namespace, getTypeFromSelfLink(meta.SelfLink), meta.Name),
				},
			},
		},
	}
	_, err = c.monitoringClient.ServiceMonitors(spec.Prometheus.Namespace).Create(serviceMonitor)
	return err
}

func (c *PrometheusController) createK8sdbExporter() (*extensions.Deployment, error) {
	exporter := &extensions.Deployment{
		ObjectMeta: kapi.ObjectMeta{
			Name:      k8sdbExporter,
			Namespace: c.exporterNamespace,
			Labels:    exporterLabel,
		},
		Spec: extensions.DeploymentSpec{
			Replicas: 1,
			Template: kapi.PodTemplateSpec{
				Spec: kapi.PodSpec{
					Containers: []kapi.Container{
						{
							Name:            exporterContainerName,
							Image:           ImageK8sdbExporter,
							ImagePullPolicy: kapi.PullIfNotPresent,
							Ports: []kapi.ContainerPort{
								{
									Name:          exporterPort,
									Protocol:      kapi.ProtocolTCP,
									ContainerPort: 9187,
								},
							},
						},
					},
				},
			},
		},
	}
	return c.kubeCLient.Extensions().Deployments(c.exporterNamespace).Create(exporter)
}

func (c *PrometheusController) createServiceForExporter() error {
	found, err := c.checkService()
	if err != nil {
		return err
	}
	if found {
		return nil
	}
	svc := &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name:      k8sdbExporter,
			Namespace: c.exporterNamespace,
			Labels:    exporterLabel,
		},
		Spec: kapi.ServiceSpec{
			Type: kapi.ServiceTypeClusterIP,
			Ports: []kapi.ServicePort{
				{
					Name:       exporterPort,
					Port:       9187,
					Protocol:   kapi.ProtocolTCP,
					TargetPort: intstr.FromString(exporterPort),
				},
			},
			Selector: exporterLabel,
		},
	}
	_, err = c.kubeCLient.Core().Services(c.exporterNamespace).Create(svc)
	return err
}

func (c *PrometheusController) checkService() (bool, error) {
	svc, err := c.kubeCLient.Core().Services(c.exporterNamespace).Get(k8sdbExporter)
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	if reflect.DeepEqual(svc.Spec.Selector, exporterLabel) /*svc.Spec.Selector != exporterLabel */ {
		return false, fmt.Errorf("Service %v already exist but selector mismatch with k8sdbexporter", svc.Name)
	}
	return true, nil

}

func (c *PrometheusController) checkServiceMonitorUpdate(svcMonitor *v1alpha1.ServiceMonitor, spec *tapi.MonitorSpec) error {
	var needUpdate bool
	if svcMonitor.Namespace != spec.Prometheus.Namespace {
		needUpdate = true
		svcMonitor.Namespace = spec.Prometheus.Namespace
	}
	if reflect.DeepEqual(svcMonitor.Labels, spec.Prometheus.Labels) /*svcMonitor.Labels != spec.Prometheus.Labels */ {
		needUpdate = true
		svcMonitor.Labels = spec.Prometheus.Labels
	}
	if needUpdate {
		_, err := c.monitoringClient.ServiceMonitors(spec.Prometheus.Namespace).Update(svcMonitor)
		return err
	}
	return nil
}

func getTypeFromSelfLink(selfLink string) string {
	if len(selfLink) == 0 {
		return ""
	}
	s := strings.Split(selfLink, "/")
	return s[len(s)-2]
}

func getServiceMonitorName(meta *kapi.ObjectMeta) string {
	return fmt.Sprintf("kubedb-%s-%s", meta.Namespace, meta.Name)
}
