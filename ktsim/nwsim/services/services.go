package services
import (
  "fmt"
  "strings"
)

type TargetData struct {
  delay     int
  transform float64
  name      string
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

func (s *Service) AddTarget(t *Service,
                            delay int,
                            transform float64,
                            name string) {
  s.targets[t] = TargetData{
    delay,
    transform,
    name,
  }
}

func (s *Service) Deployment() []byte {
  return []byte(fmt.Sprintf(
`apiVersion: apps/v1
kind: Deployment
metadata:
  name: 'service-%v'
spec:
  replicas: 1
  selector:
    matchLabels:
      id: '%v'
  template:
    metadata:
      labels:
        id: '%v'
    spec:
      containers:
      - name: service
        image: service
        imagePullPolicy: Never
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
  name: 'service-%v'
spec:
  selector:
    id: '%v'
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
  name: 'configmap-%v'
data:
  ID: '%v'
  QUEUESIZE: '%v'
  TARGETS: '%v'
  DELAYS: '%v'
  TRANSFORMS: '%v'`,
  s.id, s.id, s.queueSize, strTargets, strDelays, strTransforms))
}

func toString(tData map[*Service]TargetData) (string, string, string) {
  targets := ""
  delays := ""
  transforms := ""
  for _, v := range tData {
    targets     += fmt.Sprintf("http://%v,", v.name)
    delays      += fmt.Sprintf("%v,", v.delay)
    transforms  += fmt.Sprintf("%.2f,", v.transform)
  }
  targets     = strings.TrimRight(targets, ",")
  delays      = strings.TrimRight(delays, ",")
  transforms  = strings.TrimRight(transforms, ",")
  return targets, delays, transforms
}
