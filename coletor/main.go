package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var PollInterval time.Duration
var testMode bool

func init() {
	pollIntervalStr := os.Getenv("POLL_INTERVAL")
	if pollIntervalStr == "" {
		PollInterval = 30 * time.Second
	} else {
		pollIntervalInt, err := strconv.Atoi(pollIntervalStr)
		if err != nil {
			log.Fatalf("Erro ao converter POLL_INTERVAL para inteiro: %v", err)
		}
		PollInterval = time.Duration(pollIntervalInt) * time.Second
	}

	flag.BoolVar(&testMode, "test", false, "Rodar em modo de teste")
	flag.Parse()

	if !testMode {
		testEnv := os.Getenv("TEST_MODE")
		if testEnv == "true" {
			testMode = true
		}
	}

	fmt.Printf("Poll interval definido para: %v\n", PollInterval)
	fmt.Printf("Modo de teste: %v\n", testMode)
}

func main() {
	if testMode {
		for {
			nodeURL := "http://test:8080/metrics"
			nodeID := "test-node"
			metrics, err := CollectMetrics(nodeURL)
			if err != nil {
				log.Printf("Erro ao coletar métricas de %s: %v", nodeURL, err)
			} else {
				err = StoreMetrics(metrics, nodeID)
				if err != nil {
					log.Printf("Erro ao armazenar métricas de %s: %v", nodeURL, err)
				}
			}
			time.Sleep(PollInterval)
		}
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			kubeconfig := os.Getenv("KUBECONFIG")
			if kubeconfig == "" {
				kubeconfig = clientcmd.RecommendedHomeFile
			}
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				log.Fatalf("Erro ao construir configuração do Kubernetes: %v", err)
			}
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Erro ao criar cliente do Kubernetes: %v", err)
		}

		for {
			pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Printf("Erro ao listar nodes: %v", err)
				time.Sleep(PollInterval)
				continue
			}

			for _, pod := range pods.Items {
				podIP := pod.Status.PodIP
				podURL := fmt.Sprintf("http://%s:8080/metrics", podIP)
				nodeID := pod.Name

				go func(podURL, nodeID string) {
					metrics, err := CollectMetrics(podURL)
					if err != nil {
						log.Printf("Erro ao coletar métricas de %s: %v", podURL, err)
						return
					}

					err = StoreMetrics(metrics, nodeID)
					if err != nil {
						log.Printf("Erro ao armazenar métricas de %s: %v", podURL, err)
					}
				}(podURL, nodeID)
			}

			time.Sleep(PollInterval)
		}
	}
}
