package poddiscovery
import (
  "context"
  "fmt"
  "log"
  "os"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/rest"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/api/core/v1"
)

var logger = log.New(os.Stdout, "[poddiscovery] ", log.Ltime)

type PodInfo struct {
  Name    string
  IP      string
  HostIP  string
  Service string
}

func getK8sClient() *kubernetes.Clientset {
  var (
    config  *rest.Config
    err     error
  )
  if os.Getenv("LOCAL") == "true" {
    config, err = clientcmd.BuildConfigFromFlags("", "/Users/henriquefurst/.kube/config")
  } else {
    config, err = rest.InClusterConfig()
  }
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
  if err != nil {
    return nil, fmt.Errorf("Could not list Services: %v", err)
  }
  podsItems := pods.Items
  numPods := len(podsItems)
  logger.Printf("Found %v pods in namespace \"default\"\n", numPods)

  logger.Println("Polling Kubernetes API for Services...")
  services, err := k8sClient.CoreV1().
    Services("default").
    List(context.Background(), metav1.ListOptions{})
  if err != nil {
    return nil, fmt.Errorf("Could not list Services: %v", err)
  }
  servicesItems := services.Items
  podInfos := matchPodsAndServices(podsItems, servicesItems)
  return podInfos, nil
}

func matchPodsAndServices(podsItems []v1.Pod,
                          servicesItems []v1.Service) []PodInfo {
  serviceLabels := make(map[string]string)
  for _, service := range servicesItems {
    label := service.Spec.Selector["id"]
    serviceLabels[label] = fmt.Sprintf("http://%v", service.Name)
  }
  numPods := len(podsItems)
  podInfos := make([]PodInfo, 0, numPods)
  for _, pod := range podsItems {
    label := pod.Labels["id"]
    serviceAddr := serviceLabels[label]
    podInfos = append(
      podInfos,
      PodInfo{
        Name:     pod.Name,
        IP:       pod.Status.PodIP,
        HostIP:   pod.Status.HostIP,
        Service:  serviceAddr,
      },
    )
  }
  return podInfos
}
