package kubernetes

import (
	"fmt"

	kube "k8s.io/client-go/kubernetes"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type Namespace struct {
	ClientSet *kube.Clientset
}

func (n *Namespace) Exists(namespace string) (bool, error) {
	_, err := n.ClientSet.Core().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		// Maybe better to list all namespaces and search it so we don't have to do a dirty check
		if err.Error() == fmt.Sprintf("namespaces \"%s\" not found", namespace) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}


func (n *Namespace) Create(namespace string) (*v1.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	return n.ClientSet.Core().Namespaces().Create(ns)
}

func (n *Namespace) Ensure(namespace string) (error) {
	exists, err := n.Exists(namespace)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = n.Create(namespace)
	return err
}