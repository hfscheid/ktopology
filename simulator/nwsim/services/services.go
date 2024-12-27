package services
import (
  "fmt"
  "strings"
)

type TargetData struct {
  delay     int
  transform float64
}

type Service struct {
  id int
  queueSize int
  targets   map[*Service]TargetData
}
func New(id, qSize int) *Service {
  return &Service{
    id: id,
    queueSize: qSize,
    targets: make(map[*Service]TargetData),
  }
}

func (s *Service) AddTarget(t *Service, delay int, transform float64) {
  s.targets[t] = TargetData{delay, transform}
}

func (s *Service) Deployment() []byte {
  return []byte(fmt.Sprintf(
`apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploymet-%v
spec:
  replicas: 1
  selector:
    matchLabels:
      id: %v
  template:
    metadata:
      labels:
        id: %v
    spec:
      containers:
      - name: service
        image: service
        ports:
        - containerPort: 80
      envFrom:
      - configMapRef:
          name: configmap-%v`,
    s.id, s.id, s.id, s.id))
}

func (s *Service) Service() []byte {
  return []byte(fmt.Sprintf(
`apiVersion: v1
kind: Service
metadata:
  name: service-%v
spec:
  selector:
    app.kubernetes.io/name: deployment-%v
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80`,
  s.id, s.id))
}

func (s *Service) ConfigMap() []byte {
  strTargets, strDelays, strTransforms := toString(s.targets)
  return []byte(fmt.Sprintf(
`apiVersion: v1
kind: ConfigMap
metadata:
  name: configmap-%v
data:
  TARGETS: %v
  DELAYS: %v
  TRANSFORMS: %v`,
  s.id, strTargets, strDelays, strTransforms))
}

func toString(tData map[*Service]TargetData) (string, string, string) {
  targets := ""
  delays := ""
  transforms := ""
  for k, v := range tData {
    targets     += fmt.Sprintf("service-%v,", k.id)
    delays      += fmt.Sprintf("%v,", v.delay)
    transforms  += fmt.Sprintf("%v,", v.transform)
  }
  targets     = strings.TrimRight(targets, ",")
  delays      = strings.TrimRight(delays, ",")
  transforms  = strings.TrimRight(transforms, ",")
  return targets, delays, transforms
}
