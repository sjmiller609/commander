package kubernetes

import (
// 	"fmt"
//
	kube "k8s.io/client-go/kubernetes"
// 	"k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type Secret struct {
	ClientSet *kube.Clientset
}
//
// func (n *Secret) Exists(namespace string) (bool, error) {
// 	_, err := n.ClientSet.Core().Namespaces().Get(namespace, metav1.GetOptions{})
// 	if err != nil {
// 		// Maybe better to list all namespaces and search it so we don't have to do a dirty check
// 		if err.Error() == fmt.Sprintf("namespaces \"%s\" not found", namespace) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return true, nil
// }
//
//
// func (n *Secret) Create(namespace string) (*v1.Namespace, error) {
// 	ns := &v1.Namespace{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: namespace,
// 		},
// 	}
// 	return n.ClientSet.Core().Namespaces().Create(ns)
// }
//
// func (n *Secret) Ensure(namespace string) (error) {
// 	exists, err := n.Exists(namespace)
// 	if err != nil {
// 		return err
// 	}
// 	if exists {
// 		return nil
// 	}
//
// 	_, err = n.Create(namespace)
// 	return err
// }
/*
// createSecret creates the Tiller secret resource.
func createSecret(client corev1.SecretsGetter, opts *Options) error {
	o, err := generateSecret(opts)
	if err != nil {
		return err
	}
	_, err = client.Secrets(o.Namespace).Create(o)
	return err
}

// generateSecret builds the secret object that hold Tiller secrets.
func generateSecret(opts *Options) (*v1.Secret, error) {

	labels := generateLabels(map[string]string{"name": "tiller"})
	secret := &v1.Secret{
		Type: v1.SecretTypeOpaque,
		Data: make(map[string][]byte),
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Labels:    labels,
			Namespace: opts.Namespace,
		},
	}
	var err error
	if secret.Data["tls.key"], err = read(opts.TLSKeyFile); err != nil {
		return nil, err
	}
	if secret.Data["tls.crt"], err = read(opts.TLSCertFile); err != nil {
		return nil, err
	}
	if opts.VerifyTLS {
		if secret.Data["ca.crt"], err = read(opts.TLSCaCertFile); err != nil {
			return nil, err
		}
	}
	return secret, nil
}
 */