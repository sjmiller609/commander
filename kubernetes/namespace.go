package kubernetes

import (
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
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

func (n *Namespace) Create(namespace string, labels map[string]string) (*v1.Namespace, error) {

	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: labels,
		},
	}
	return n.ClientSet.Core().Namespaces().Create(ns)
}

func (n *Namespace) Delete(namespace string) error {
	// set options
	propagation := metav1.DeletePropagationForeground
	gracePeriod := int64(0)

	// create k8 DeleteOptions object
	options := metav1.DeleteOptions{
		PropagationPolicy:  &propagation,
		GracePeriodSeconds: &gracePeriod,
	}

	return n.ClientSet.Core().Namespaces().Delete(namespace, &options)
}

func (n *Namespace) Ensure(namespace string, labels map[string]string) error {
	exists, err := n.Exists(namespace)
	if err != nil {
		return err
	}
	if exists {
		n.UpdateLabels(namespace, labels)
		return nil
	}

	_, err = n.Create(namespace, labels)
	return err
}

func (n *Namespace) UpdateLabels(namespace string, labels map[string]string) error {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: labels,
		},
	}
	_, err := n.ClientSet.Core().Namespaces().Update(ns)
	return err
}
