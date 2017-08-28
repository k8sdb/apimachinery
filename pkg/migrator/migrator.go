package migrator

/*


import (
	"errors"
	"fmt"
	"time"
	"github.com/appscode/log"
	"github.com/hashicorp/go-version"
	aci "github.com/k8sdb/apimachinery/api"
	tcs "github.com/k8sdb/apimachinery/client/clientset"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type migrationState struct {
	tprRegDeleted bool
	crdCreated    bool
}

type migrator struct {
	kubeClient       clientset.Interface
	apiExtKubeClient apiextensionsclient.Interface
	extClient        tcs.ExtensionInterface

	migrationState *migrationState
}

func NewMigrator(kubeClient clientset.Interface, apiExtKubeClient apiextensionsclient.Interface, extClient tcs.ExtensionInterface) *migrator {
	return &migrator{
		migrationState:   &migrationState{},
		kubeClient:       kubeClient,
		apiExtKubeClient: apiExtKubeClient,
		extClient:        extClient,
	}
}

func (m *migrator) isMigrationNeeded(runtimeObjs ...aci.RuntimeObject) (bool, error) {
	v, err := m.kubeClient.Discovery().ServerVersion()
	if err != nil {
		return false, err
	}

	ver, err := version.NewVersion(v.String())
	if err != nil {
		return false, err
	}

	mv := ver.Segments()[1]

	if mv == 7 {
		for _, runtime := range runtimeObjs {
			_, err := m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(
				runtime.ResourceName()+"."+aci.V1alpha1SchemeGroupVersion.Group,
				metav1.GetOptions{},
			)
			if err != nil {
				if !kerr.IsNotFound(err) {
					return false, err
				}
			} else {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *migrator) RunMigration(runtimeObjs ...aci.RuntimeObject) error {
	needed, err := m.isMigrationNeeded(runtimeObjs...)
	if err != nil {
		return err
	}

	if needed {
		if err := m.migrateTPR2CRD(runtimeObjs...); err != nil {
			return m.rollback(runtimeObjs...)
		}
	}

	return nil
}

func (m *migrator) migrateTPR2CRD(runtimeObjs ...aci.RuntimeObject) error {
	log.Debugln("Performing TPR to CRD migration.")

	log.Debugln("Deleting TPRs.")
	if err := m.deleteTPRs(runtimeObjs...); err != nil {
		return errors.New("Failed to Delete TPRs")
	}

	m.migrationState.tprRegDeleted = true

	log.Debugln("Creating CRDs.")
	if err := m.createCRDs(runtimeObjs...); err != nil {
		return errors.New("Failed to create CRDs")
	}

	m.migrationState.crdCreated = true

	log.Debugln("Waiting for CRDs to be ready.")
	if err := m.waitForCRDsReady(len(runtimeObjs)); err != nil {
		return errors.New("Failed to be ready CRDs")
	}

	return nil
}

func (m *migrator) deleteTPRs(runtimeObjs ...aci.RuntimeObject) error {
	tprClient := m.kubeClient.ExtensionsV1beta1().ThirdPartyResources()

	deleteTPR := func(runtime aci.RuntimeObject) error {
		name := runtime.ResourceName() + "." + aci.V1alpha1SchemeGroupVersion.Group
		if err := tprClient.Delete(name, &metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("failed to remove %s TPR", name)
		}
		return nil
	}

	for _, runtime := range runtimeObjs {
		if err := deleteTPR(runtime); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) createCRDs(runtimeObjs ...aci.RuntimeObject) error {
	for _, runtime := range runtimeObjs {
		if err := m.createCRD(runtime); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) createCRD(runtime aci.RuntimeObject) error {
	crd := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: runtime.ResourceType() + "." + aci.V1alpha1SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "kubedb",
			},
		},
		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   aci.V1alpha1SchemeGroupVersion.Group,
			Version: aci.V1alpha1SchemeGroupVersion.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural:   runtime.ResourceType(),
				Singular: runtime.ResourceCode(),
				Kind:     runtime.ResourceKind(),
			},
		},
	}

	crdClient := m.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	_, err := crdClient.Create(crd)
	if err != nil && !kerr.IsAlreadyExists(err) {
		return fmt.Errorf(`Failed to create CRD "%v"`, crd.Spec.Names.Kind)
	}

	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crdEst, err := crdClient.Get(crd.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crdEst.Status.Conditions {
			switch cond.Type {
			case extensionsobj.Established:
				if cond.Status == extensionsobj.ConditionTrue {
					return true, err
				}
			case extensionsobj.NamesAccepted:
				if cond.Status == extensionsobj.ConditionFalse {
					fmt.Printf("Name conflict. Reason: %v\n", cond.Reason)
				}
			}
		}
		return false, err
	})

	return nil
}

func (m *migrator) waitForCRDsReady(expectedCRD int) error {
	labelMap := map[string]string{
		"app": "kubedb",
	}

	return wait.Poll(3*time.Second, 10*time.Minute, func() (bool, error) {
		crdList, err := m.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions().List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labelMap).String(),
		})
		if err != nil {
			return false, err
		}

		if len(crdList.Items) == expectedCRD {
			return true, nil
		}

		return false, nil
	})
}

func (m *migrator) rollback(runtimeObjs ...aci.RuntimeObject) error {
	log.Debugln("Rolling back migration.")

	ms := m.migrationState

	if ms.crdCreated {
		log.Debugln("Deleting CRDs.")
		err := m.deleteCRDs()
		if err != nil {
			return errors.New("Failed to delete CRDs")
		}
	}

	if ms.tprRegDeleted {
		log.Debugln("Creating TPRs.")
		err := m.CreateTPRs()
		if err != nil {
			return errors.New("Failed to recreate TPRs")
		}

		err = m.WaitForTPRsReady(len(runtimeObjs))
		if err != nil {
			return errors.New("Failed to be ready TPRs")
		}
	}

	return nil
}

func (m *migrator) deleteCRDs(runtimeObjs ...aci.RuntimeObject) error {
	crdClient := m.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()

	deleteCRD := func(runtime aci.RuntimeObject) error {
		name := runtime.ResourceType() + "." + aci.V1alpha1SchemeGroupVersion.Group
		err := crdClient.Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf(`Failed to delete CRD "%s""`, name)
		}
		return nil
	}

	for _, runtime := range runtimeObjs {
		if err := deleteCRD(runtime); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) CreateTPRs(runtimeObjs ...aci.RuntimeObject) error {
	for _, runtime := range runtimeObjs {
		if err := m.createTPR(runtime); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) createTPR(runtime aci.RuntimeObject) error {
	name := runtime.ResourceName() + "." + aci.V1alpha1SchemeGroupVersion.Group
	_, err := m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(name, metav1.GetOptions{})
	if !kerr.IsNotFound(err) {
		return err
	}

	thirdPartyResource := &extensions.ThirdPartyResource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "ThirdPartyResource",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": "kubedb",
			},
		},
		Description: "Searchlight by AppsCode - Alerts for Kubernetes",
		Versions: []extensions.APIVersion{
			{
				Name: aci.V1alpha1SchemeGroupVersion.Version,
			},
		},
	}

	_, err = m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Create(thirdPartyResource)
	return err
}

func (m *migrator) WaitForTPRsReady(expectedTPR int) error {
	labelMap := map[string]string{
		"app": "kubedb",
	}

	return wait.Poll(3*time.Second, 10*time.Minute, func() (bool, error) {
		crdList, err := m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labelMap).String(),
		})
		if err != nil {
			return false, err
		}

		if len(crdList.Items) == expectedTPR {
			return true, nil
		}

		return false, nil
	})
}
*/
