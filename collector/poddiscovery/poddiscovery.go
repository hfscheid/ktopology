package poddiscovery
import (
  "context"
  "fmt"
  "log"
  "os"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
  // "k8s.io/client-go/tools/clientcmd"
)

var logger = log.New(os.Stdout, "[poddiscovery] ", log.Ltime)

type PodInfo struct {
  Name    string
  IP      string
  HostIP  string
}

func getK8sClient() *kubernetes.Clientset {
  config, err := rest.InClusterConfig()
  if err != nil {
    logger.Fatalf("Could not get in-cluster Kubernetes config: %v", err)
  }
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    logger.Fatalf("Could not create Kubernetes clientset: %v", err)
  }
  return clientset
}

func ListPods() ([]PodInfo, error) {
  k8sClient := getK8sClient()
  logger.Println("Polling Kubernetes API for Pod IPs...")
  pods, err := k8sClient.CoreV1().
    Pods("default").
    List(context.Background(), metav1.ListOptions{})
  podsItems := pods.Items
  if err != nil {
    return nil, fmt.Errorf("Could not list nodes: %v", err)
  }
  numPods := len(podsItems)
  logger.Printf("Found %v pods in namespace \"default\"\n", numPods)

  podInfos := make([]PodInfo, 0, numPods)
  for _, pod := range podsItems {
    podInfos = append(
      podInfos,
      PodInfo{
        Name:   pod.Name,
        IP:     pod.Status.PodIP,
        HostIP: pod.Status.HostIP,
      },
    )
  }
  return podInfos, nil
}
