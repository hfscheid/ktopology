package main
import (
  "context"
  "log"
  "os"
  "math/rand"
  "time"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/clientcmd"
  appsv1 "k8s.io/api/apps/v1"
  typedappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

var logger = log.New(os.Stdout, "[alterator] ", log.Ltime)

func getK8sClient() *kubernetes.Clientset {
  config, err := clientcmd.BuildConfigFromFlags("", "/Users/henriquefurst/.kube/config")
  if err != nil {
    logger.Fatalf("Could not get Kubernetes config: %v", err)
  }
  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    logger.Fatalf("Could not create Kubernetes clientset: %v", err)
  }
  return clientset
}

func randomPick(ds []appsv1.Deployment) appsv1.Deployment {
  var pick appsv1.Deployment
  for {
    pick = ds[rand.Intn(len(ds))]
    if  pick.Name != "collector" &&
        pick.Name != "collector-mongodb" &&
        pick.Name != "sink" {
      break
    }
  }
  return pick
}

func randomUpdate(dc *typedappsv1.DeploymentInterface,
                  d *appsv1.Deployment) {
  replicas := *(d.Spec.Replicas)
  var newReplicas int32
  for {
    newReplicas = int32(rand.Intn(4)+1)
    if  newReplicas > 0 &&
        newReplicas != replicas {
      break
    }
  }
  d.Spec.Replicas = &newReplicas
  log.Printf("Alterating %v to %v replicas...\n", d.Name, newReplicas)
  (*dc).Update(context.Background(), d, metav1.UpdateOptions{})
}

func main() {
  k8sClient := getK8sClient()
  deploymentsClient := k8sClient.AppsV1().Deployments("default")
  for {
    deployments, _ := deploymentsClient.List(context.Background(), metav1.ListOptions{})
    deploymentsItems := deployments.Items
    for i := rand.Intn(3); i > 0; i-- {
      target := randomPick(deploymentsItems)
      randomUpdate(&deploymentsClient, &target)
    }
    time.Sleep(20*time.Second)
  }
}
