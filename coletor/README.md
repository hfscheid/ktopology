# Coletor de Métricas

Este projeto coleta métricas de nodes e as armazena em um banco de dados MongoDB.

## Pré-requisitos

- Docker e Docker Compose instalados
- Go instalado

## Iniciando os Containers Docker

1. Certifique-se de que o Docker está em execução.
2. Navegue até o diretório do projeto onde o arquivo `docker-compose.yml` está localizado.
3. Defina no `docker-compose.yml` as variáveis de ambiente MONGO_URI, POLL_INTERVAL, KUBECONFIG e TEST_MODE.
4. Execute o seguinte comando para iniciar os containers:
```sh
docker compose up --build -d
```

## Rodando o Programa Localmente

1. Certifique-se de que um container MongoDB está em execução.
2. Navegue até o diretório do projeto onde os arquivos Go estão localizados.
3. Inicialize um novo módulo Go:
```sh
go mod init coletor
```
4. Adicione as dependências necessárias:
```sh
go get ./...
```
5. Compile e execute o programa Go (Utilize a flag -test caso não esteja rodando um ambiente Kubernetes):
```sh
go run main.go collector.go storage.go models.go -test
```