# Coletor de Métricas

Este projeto coleta métricas de nodes e as armazena em um banco de dados MongoDB.

## Pré-requisitos

- Docker e Docker Compose instalados
- Go instalado

## Iniciando o Container Docker do MongoDB

1. Certifique-se de que o Docker está em execução.
2. Navegue até o diretório do projeto onde o arquivo `docker-compose.yml` está localizado.
3. Execute o seguinte comando para iniciar o container MongoDB:
```sh
docker compose up --build -d
```

## Iniciando o coletor

1. Navegue até o diretório do projeto onde os arquivos Go estão localizados.
2. Inicialize um novo módulo Go:
```sh
go mod init coletor
```

3. Adicione as dependências necessárias:
```sh
go get go.mongodb.org/mongo-driver/mongo
```

4. Compile e execute o programa Go:
```sh
go run main.go collector.go storage.go models.go
```