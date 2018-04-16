package helm

import (
	"fmt"

	"k8s.io/helm/pkg/helm/portforwarder"
)

func (c *Client) OpenTunnel() error {
	if c.settings.TillerHost == "" {
		tunnel, err := portforwarder.New(c.settings.TillerNamespace, c.kubeClient.ClientSetInterface(), c.kubeClient.Config)
		if err != nil {
			return err
		}
		c.tillerTunnel = tunnel

		c.settings.TillerHost = fmt.Sprintf("127.0.0.1:%d", tunnel.Local)
		fmt.Printf("Created tunnel using local port: '%d'\n", tunnel.Local)
	}
	return nil
}

func (c *Client) CloseTunnel() {
	if c.tillerTunnel != nil {
		c.tillerTunnel.Close()
	}
}