package kubi_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-maestro/kubi"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("when KUBECONFIG points to a existing file", func(t *testing.T) {
		pkgPath, _ := os.Getwd()
		kubeconfig := fmt.Sprintf("%s/data/kubeconfig", pkgPath)

		os.Setenv("KUBECONFIG", kubeconfig)

		client, err := kubi.NewClient()

		t.Run("shouldn't return error", func(t *testing.T) {
			assert.Nil(t, err)
		})

		t.Run("should return the client", func(t *testing.T) {
			assert.IsType(t, &kubi.Client{}, client)
		})

		t.Run("should set client.RestConfig according", func(t *testing.T) {
			assert.Equal(t, client.RestConfig.Host, "https://0.0.0.0:46271")
		})

		os.Unsetenv("KUBECONFIG")
	})

	t.Run("when KUBECONFIG points to a non-existing file", func(t *testing.T) {
		pkgPath, _ := os.Getwd()
		kubeconfig := fmt.Sprintf("%s/i/doesn't/exists", pkgPath)

		os.Setenv("KUBECONFIG", kubeconfig)

		client, err := kubi.NewClient()

		t.Run("shouldn't return the client", func(t *testing.T) {
			assert.Nil(t, client)
		})

		t.Run("should return an error", func(t *testing.T) {
			assert.Error(t, err)
			assert.IsType(t, &os.PathError{}, err)
		})

		os.Unsetenv("KUBECONFIG")
	})

	t.Run("when KUBECONFIG wasn't set", func(t *testing.T) {
		t.Run("when application is running inside Kubernetes", func(t *testing.T) {
			assert.NotEmpty(t, os.Getenv("KUBERNETES_SERVICE_HOST"))
			assert.NotEmpty(t, os.Getenv("KUBERNETES_SERVICE_PORT"))

			assert.FileExists(t, "/var/run/secrets/kubernetes.io/serviceaccount/token")
			assert.FileExists(t, "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")

			k8sEndpoint := fmt.Sprintf("https://%s:%s",
				os.Getenv("KUBERNETES_SERVICE_HOST"),
				os.Getenv("KUBERNETES_SERVICE_PORT"),
			)

			client, err := kubi.NewClient()

			t.Run("should return the client", func(t *testing.T) {
				assert.IsType(t, &kubi.Client{}, client)
			})

			t.Run("should set client.RestConfig according", func(t *testing.T) {
				assert.Equal(t, client.RestConfig.Host, k8sEndpoint)
			})

			t.Run("shouldn't return an error", func(t *testing.T) {
				assert.Nil(t, err)
			})
		})

		t.Run("when application isn't running inside Kubernetes", func(t *testing.T) {
			os.Unsetenv("KUBERNETES_SERVICE_HOST")
			os.Unsetenv("KUBERNETES_SERVICE_PORT")

			client, err := kubi.NewClient()

			t.Run("should return an error", func(t *testing.T) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "without kubernetes access")
			})

			t.Run("shouldn't return the client", func(t *testing.T) {
				assert.Nil(t, client)
			})
		})
	})
}
