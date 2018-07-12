package kubernetes

import (
	"fmt"

	kube "k8s.io/client-go/kubernetes"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type Secret struct {
	ClientSet *kube.Clientset
}

func (s *Secret) Create(name string, key string, value string, namespace string, releaseName string) (*v1.Secret, error) {
	labels := map[string]string{
		"component": "airflow-secret",
		"tier": "airflow",
		"release": releaseName,
	}
	data := make(map[string][]byte)
	data[key] = []byte(value)

	secret := &v1.Secret{
		Type: v1.SecretTypeOpaque,
		Data: data,
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Labels:    labels,
			Namespace: namespace,
		},
	}

	return s.ClientSet.CoreV1().Secrets(namespace).Create(secret)
}

func (s *Secret) DeleteByRelease(releaseName string, namespace string) error {
	deleteOptions := &meta_v1.DeleteOptions{
	}
	listOptions := meta_v1.ListOptions{
		LabelSelector: fmt.Sprintf("release=%s", releaseName),
	}
	return s.ClientSet.CoreV1().Secrets(namespace).DeleteCollection(deleteOptions, listOptions)
}