# Requisitos para simulação de rede
- [docker](https://docs.docker.com/engine/install/)
- [minikube](https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Fx86-64%2Fstable%2Fbinary+download)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
- [Go](https://go.dev/doc/install)

# Passo 0: configurar ambiente de execução
Se assegure que o ambiente minikube utilize um repositório local de imagens
Docker:

```
# inicie o cluster miniube
minikube start
# utilize o repositório de images local do minikube com docker
eval $(minikube docker-env)
```


# Passo 1: gerar topologia da rede simulada
Dentro do diretório `k8stopology/simulator/nwsim` execute o comando:

`go run . n`

onde `n` é o número desejado de vértices na rede. O programa popula o
diretório `k8stopology/simulator/nwsim/manifests` com os manifestos kubernetes
para a instanciação dos Pods e serviços correspondentes.

# Passo 2: gerar imagem dos serviços
Dentro do diretório `k8stopology/simulator/svcsim` execute o seguinte comando
para criar a imagem Docker dos serviços simulados:

`docker build -t service .`

> A imagem só será gerada no repositório correto, legível pelo minikube, se o
> passo 0 for observado.

> O nome das imagens é referenciado em todos os manifestos gerados no passo 1
> (Deployment.spec.template.spec.containers[0].image). É fundamental que 
> correspondam perfeitamente.

De forma semelhante, crie a imagem do gerador-coletor de fluxo executando o
seguinte comando a partir do diretório `k8stopology/simulator/fluxgen`:

`docker build -t sink .`

> Aqui também é fundamental observância do passo 0 e do devido nome da imagem.

# Passo 3: popular o cluster
```
# a partir de k8stopology/simulator:
kubectl apply -f nwsim/manifests/
# aguarde a criação de toda a rede antes de criar o gerador-coletor de fluxo
kubectl apply -f fluxgen/manifest.yaml
```

O gerador-coletor de fluxo começa a enviar requisições automaticamente a partir
de sua criação, e a rede está em funcionamento. A frequência com que o fluxo é
gerado é configurada via ConfigMap manifesto em
`k8stopology/simulator/fluxgen/manifest.yaml`.

> A rede ainda não é acessível por fora. Toda interação é feita através do
> comando `kubectl` ou da criação de Pods que interajam com os demais. Uma
> exceção é o uso de **port-forwarding**, que permite associar um porto local a
> um porto de um serviço da rede.
