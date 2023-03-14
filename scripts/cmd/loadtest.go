package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeConfigPath := filepath.Join(homeDir, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err)
	}
	config.QPS = 100_000

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	groupVersionResource := schema.GroupVersionResource{
		Group:    "mygroup.github.com",
		Version:  "v1",
		Resource: "mycustomresources",
	}

	const numResources = 100_000

	var wg sync.WaitGroup
	wg.Add(numResources)

	for i := 0; i < numResources; i++ {
		if i%1000 == 0 {
			fmt.Println("creating resource", i)
		}
		name := fmt.Sprintf("mycustomresource.%s", randString(32))

		go func(wg *sync.WaitGroup) {
			_, err := dynamicClient.Resource(groupVersionResource).Namespace("default").Apply(
				context.Background(),
				name,
				&unstructured.Unstructured{
					Object: map[string]any{
						"apiVersion": "mygroup.github.com/v1",
						"kind":       "MyCustomResource",
						"metadata": map[string]any{
							"name": name,
						},
						"spec": map[string]any{
							"cluster":  "us-east-1",
							"status":   "active",
							"database": "my-database-host:port",
						},
					},
				},
				metav1.ApplyOptions{FieldManager: "fieldManagerName"},
			)
			wg.Done()

			if err != nil {
				fmt.Println("got error", err)
			}
		}(&wg)
	}

	wg.Wait()
}

func randString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyz"

	buffer := strings.Builder{}
	buffer.Grow(n)

	for i := 0; i < n; i++ {
		if err := buffer.WriteByte(letters[rand.Intn(len(letters))]); err != nil {
			panic(err)
		}
	}

	return buffer.String()
}
