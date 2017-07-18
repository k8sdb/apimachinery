package fake

import (
	"github.com/k8sdb/apimachinery/client/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/testing"
)

type FakeExtensionClient struct {
	*testing.Fake
}

var _ clientset.ExtensionInterface = &FakeExtensionClient{}

func NewFakeExtensionClient(objects ...runtime.Object) *FakeExtensionClient {
	o := testing.NewObjectTracker(api.Registry, api.Scheme, api.Codecs.UniversalDecoder())
	for _, obj := range objects {
		if obj.GetObjectKind().GroupVersionKind().Group == "kubedb.com" {
			if err := o.Add(obj); err != nil {
				panic(err)
			}
		}
	}

	fakePtr := testing.Fake{}
	fakePtr.AddReactor("*", "*", testing.ObjectReaction(o, api.Registry.RESTMapper()))

	fakePtr.AddWatchReactor("*", testing.DefaultWatchReactor(watch.NewFake(), nil))

	return &FakeExtensionClient{&fakePtr}
}

var _ clientset.ExtensionInterface = &FakeExtensionClient{}

func (m *FakeExtensionClient) Snapshots(ns string) clientset.SnapshotInterface {
	return &FakeSnapshot{m.Fake, ns}
}

func (m *FakeExtensionClient) DormantDatabases(ns string) clientset.DormantDatabaseInterface {
	return &FakeDormantDatabase{m.Fake, ns}
}

func (m *FakeExtensionClient) Elasticsearches(ns string) clientset.ElasticsearchInterface {
	return &FakeElasticsearch{m.Fake, ns}
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
