package fake

import (
	"github.com/k8sdb/apimachinery/client/clientset"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apimachinery/registered"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	testing "k8s.io/kubernetes/pkg/client/testing/core"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"
)

type FakeExtensionClient struct {
	*testing.Fake
}

var _ clientset.ExtensionInterface = &FakeExtensionClient{}

func NewFakeExtensionClient(objects ...runtime.Object) *FakeExtensionClient {
	o := testing.NewObjectTracker(api.Scheme, api.Codecs.UniversalDecoder())
	for _, obj := range objects {
		if obj.GetObjectKind().GroupVersionKind().Group == "k8sdb.com" {
			if err := o.Add(obj); err != nil {
				panic(err)
			}
		}
	}

	fakePtr := testing.Fake{}
	fakePtr.AddReactor("*", "*", testing.ObjectReaction(o, registered.RESTMapper()))

	fakePtr.AddWatchReactor("*", testing.DefaultWatchReactor(watch.NewFake(), nil))

	return &FakeExtensionClient{&fakePtr}
}

var _ clientset.ExtensionInterface = &FakeExtensionClient{}

func (m *FakeExtensionClient) Snapshots(ns string) clientset.SnapshotInterface {
	return &FakeSnapshot{m.Fake, ns}
}

func (m *FakeExtensionClient) DeletedDatabases(ns string) clientset.DeletedDatabaseInterface {
	return &FakeDeletedDatabase{m.Fake, ns}
}

func (m *FakeExtensionClient) Elastics(ns string) clientset.ElasticInterface {
	return &FakeElastic{m.Fake, ns}
}

func (m *FakeExtensionClient) Postgreses(ns string) clientset.PostgresInterface {
	return &FakePostgres{m.Fake, ns}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeExtensionClient) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
