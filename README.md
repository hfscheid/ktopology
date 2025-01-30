O *Monitor* encontra-se no diretório *ktmonitor*. Sua instalação em um cluster é
feita primeiro criando-se sua imagem Docker:
```
# dentro do diretório ktmonitor:
docker build -t ktmonitor .
```

A imagem gerada deve se encontrar em um repositório Docker utilizado pelo
_cluster_ Kubernetes de interesse. Para _clusters_ locais criados via
*minikube*, por exemplo, faz-se:
```
minikube start
eval $(minikube docker-env)
cd ktmonitor
docker build -t ktmonitor .
```

O diretório *ktsim* implementa a rede simulada. Para rodar a simulação completa,
devem-se criar as imagens Docker dos microsserviços e do gerador de fluxo:
```
minikube start
eval $(minikube docker-env)
# no diretório ktsim/svcsim
docker build -t service .

# no diretório ktsim/fluxgen
docker build -t sink .
```

Para se criarem os serviços no cluster kubernetes, é necessário primeiro gerar
os manifestos de uma topologia aleatória:
```
# no diretório ktsim/nwsim
# gerando uma rede com 14 microsserviços
go run . 14
kubectl apply -f manifests/
```
Em redes grandes é possível que algum manifesto seja defeituoso. Basta repetir o
comando e gerar uma nova rede.

O alterador da rede é um programa Go. Para se executá-lo depois de criada a
rede:
```
# no diretório ktsim/alterator
go run .
```

# Rodando toda a simulação
```
# lendo o diretório raiz
export ROOT=pwd

# criando o cluster e configurando seu repositório docker
minikube start
eval $(minikube docker-env)

# criando a rede
cd $ROOT/ktsim/svcsim
docker build -t service .
cd $ROOT/ktsim/fluxgen
docker build -t sink .
cd $ROOT/ktsim/nwsim
go run . 14

# criando o monitor
cd $ROOT/ktmonitor
docker build -t ktmonitor .

# cirando os pods no cluster
cd $ROOT
kubectl apply -f ktsim/nwsim/manifests
kubectl apply -f ktsim/fluxgen/
kubectl apply -f ktmonitor

# executando o alterador
cd $ROOT/ktsim/alterator
go run .
```

Para se pegar os dados, pegue o ID do ktmonitor:
```
kubectl get pods
```
E use seu ID para copiar os dados:
```
kubectl cp ktmonitor-[...]:/usr/src/app/data.json ./data.jsonl
```

Para terminar a simulação, encerre o alterador. Em seguida:
```
kubectl delete -f ktmonitor -f ktsim/fluxgen -f ktsim/nwsim/manifests
minikube stop
```
