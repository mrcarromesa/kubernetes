# gRPC project


- Dicas para configurar o vscode para o go [Setting Golang Plugin on VSCode for Autocomplete and Auto-import](https://medium.com/backend-habit/setting-golang-plugin-on-vscode-for-autocomplete-and-auto-import-30bf5c58138a)

- Para criar o projeto utilizar o comando:

```shell
go mod init github.com/user/repo
```

- Alterar o /user pelo seu usuário
- Alterar o /repo pelo seu repositório

- Precisamos instalar o [protocol buffers](https://developers.google.com/protocol-buffers/docs/gotutorial)

- Primeiro executamos o seguinte comando:

```shell
go get google.golang.org/protobuf/cmd/protoc-gen-go
```

- Feito isso ele adiciona essa dependencia no go.mod

- Após isso podemos executar:

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

- Adicionar:

```shell
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

- Feito isso ele adiciona essa dependencia no go.mod

- Após isso podemos executar:

```shell
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest  
```

- Com o pacote instalado nós temos acesso aos comandos:

  - protoc
  - protoc-gen-go
  - protocpgen-go-grpc


- Criamos o arquivo `proto/user.proto`

- Nele adicionamos a estrutura que teremos no tipo User que criamos.
- Esse arquivo é o que irá definir o contrato.
- O número que vem após cada atributo é o indice deles

- Ainda nesse arquivo adicionamos o service:

```proto
message User {
  /// ...
}

service UserService {
  rpc AddUser (User) returns (User);
}
```

- No caso criamos um metodo AddUser que recebe um User que é a estrutura que criamos nesse arquivo também.. e basicamente ele recebe um User e retorna esse User

- Criamos uma pasta chamada `pb` que será gerado os arquivos de protocol buffers da aplicação...
- E para o nome do pacote ser `pb` precisamos ajustar no arquivo `proto/user.proto`:

```proto
option go_package = ".;pb";
```

- Ali é ";" mesmo...

- Por fim executamos o comando:

```shell
protoc --proto_path=proto proto/*.proto --go_out=pb
```

- o atributo `--proto_path` colcamos o caminho dos nossos arquivos proto
- Eu informo para ele pegar todos os arquivos `proto/*.proto`
- E gerar na pasta `pb` pelo atributo --go_out=pb

- Caso de o seguinte erro:

```shell
command not found: protoc
```

- No caso do mac, considerando que já tenha o brew instalado execute o comando:

```shell
brew install protobuf
```

- Mais detalhes em [Installing protoc](http://google.github.io/proto-lens/installing-protoc.html)

- Adicionalmente pode ser necesssário realizar isso aqui também [protoc-gen-go: program not found or is not executable](https://stackoverflow.com/questions/57700860/protoc-gen-go-program-not-found-or-is-not-executable)

---

- Apos o comando: 

```shell
protoc --proto_path=proto proto/*.proto --go_out=pb
```

- ele irá gerar `pb/user.pb.go`

- Porém ele não gerou o grpc... para isso utilizamos o mesmo comando porém adicionando o atributo do grpc:

```shell
protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
```

---

- Criar o arquivo `cmd/server/server.go` servirá para subir o lado servidor

- Feito isso executar o comando:

```shell
go run cmd/server/server.go
```

----


## Criação dos services

- Lembrando que antes de mais nada precisamos ter gerado os protos conforme comandos anteriores

- Criamos a pasta services/

- Criamos o arquivo services/user.go

- Importante veja os comentários

```go
type UserService struct {
	pb.UnimplementedUserServiceServer // colocar essa parada aqui e boas pois se eu não implementar algo que esteja no proto não dará error!
}

func NewUserService() *UserService { // Serve como constructor
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) { // Estou recebendo o usuário e retornando ele
	// Insert - Database

	fmt.Println(req.Name)

	return &pb.User{
		Id:    "123",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}
```

- Feito isso implementamos esse service em cmd/server/server.go:

```go
func main() {
	lis, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService()) // <- Adicionar o service

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Could not serve: %v", err)
	}
}
```


----

### Executar o projeto

- Executar o comando 

```shell
go run cmd/server/server.go
```

- Utilizar o projeto [Evans](https://github.com/ktr0731/evans)
- Instalar ele na máquina
- E em outra aba ou terminal executar o comando:

```shell
evans -r repl --host localhost --port 50051
```

- repl é o modo que vamos trabalhar
- host é localhost
- a porta precisa ser a mesma definida em `cmd/server/server.go`:

```go
	lis, err := net.Listen("tcp", "localhost:50051")
```

- Ao executar aparecerá isso:

![Evans](./readme/assets/evans_hello_world.png?raw=true "Evans")


- Depois podemos escolher o service:

![Evans](./readme/assets/select_service.png?raw=true "Evans")


- service userService:

![Evans](./readme/assets/userService.png?raw=true "Evans")

- Depois adicionar o usuário:

![Evans](./readme/assets/call_addUser.png?raw=true "Evans")

- Depois salvamos o usuário e obtemos o retorno:

![Evans](./readme/assets/save_user.png?raw=true "Evans")


---

### Client

- Para realizar as chamadas gRPC criamos um client em `cmd/client/client.go`

- Para testarmos vamos executar o comando:

- No terminal 1:

```shell
go run cmd/server/server.go
```

- No terminal 2:

```shell
go run cmd/client/client.go
```

---

### Server streams

- Para implementar um server stream ajustamos isso no `proto/user.proto`:

```proto
message UserResultStream {
  string status = 1;
  User user = 2;
}

service UserService {
  rpc AddUser (User) returns (User);
  rpc AddUserVerbose (User) returns (stream UserResultStream);
}
```

- E recompilamos o proto:

```shell
protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
```

- Realizando isso será gerado as functions em `pb/user_grpc.pb.go`, `user.pb.go`

- E adicionamos a func `AddUserVerbose` em `services/user.go`

---

### Implementando o client para Server Stream

- Ajustamos o `cmd/client/client.go` para adicionar a func `AddUserVerbose` e no metodo main adicionamos a chamada para esse metodo.

- Por fim executamos o seguinte comando:

```shell
go run cmd/server/server.go
```

```shell
go run cmd/client/client.go
```

---

### ClientStream

- Para isso ajustamos o `proto/user.proto`:

```proto
message Users {
  repeated User user = 1;
}

service UserService {
  // ...
  rpc AddUsers(stream User) returns (Users);
}
```

- Executamos o comando: 

```shell
protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
```

- Após isso adicionamos a func `AddUsers` em `services/user.go`

- E implementamos o client em `cmd/client/client.go`:

```go
func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "R1",
			Name:  "Rodolfo",
			Email: "example1@email.com",
		},
		{
			Id:    "R2",
			Name:  "Rodolfo 2",
			Email: "example2@email.com",
		},
		{
			Id:    "R3",
			Name:  "Rodolfo 3",
			Email: "example3@email.com",
		},
		{
			Id:    "R4",
			Name:  "Rodolfo 4",
			Email: "example4@email.com",
		},
		{
			Id:    "R5",
			Name:  "Rodolfo 5",
			Email: "example5@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

```

E nesse mesmo arquivos adicionamos a chamada na func `main`:

```go
AddUsers(client)
```

- Por fim no terminal 1 executamos o comando:

```shell
go run cmd/server/server.go
```

- E no terminal 2 executar o comando:

```shell
go run cmd/client/client.go
```

- O client envia o stream de dados o servidor recebe e quando finaliza todo o envio o servidor retorna tudo para o client

---

### Stream bi-direcional

- No arquivo `proto/user.proto` adicionamos o seguinte:

```proto
rpc AddUserStreamBoth (stream User) returns (stream UserResultStream);
```

- Depois geramos os arquivos `.pb.go` executando o seguinte comando:

```shell
protoc --proto_path=proto proto/*.proto --go_out=pb --go-grpc_out=pb
```

- No arquivo `services/user.go` adicionamos a func `AddUserStreamBoth`

- Por fim vamos preparar o client `cmd/client/client.go`, criando a func `AddUserStreamBoth` e nele vamos utilizar o recurso do go chamado de `goroutines` é um tipo de thread controlado pelo golang e podemos criar milhões dessas threads do go, dessa forma vamos criar uma goroutines para enviar as informações e outra para receber as informações, dessa forma o sistema fica
esperando para sempre enviando e recebendo essa informação, utilizando uma função anonima:

```go
go func() {
	// TODO
}()

// A func acima, será executada o dará continuidade ao fluxo, ou seja tudo que está abaixo dela será executado, e ela ficará executando indefinidamente
```

- Adicionamos essas funções anonimas:

```go
// Para ficar enviando e
	// anonymos func
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		// Parei de enviar aqui
		stream.CloseSend()
	}()

	// Em paralelo/concorrente
	// ficar recebendo
	// Quando o servidor para de enviar ele cai no break e é encerrada
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}

			fmt.Printf("Recebendo user %v com status: %v", res.GetUser().GetName(), res.GetStatus())
		}
	}()
```

- Porém dessa forma o script será executado, chegou no final do script simplesmente tudo será encerrado, para evitar que isso aconteça vamos utilizar um recurso do golang, o `channel`

- Um channel é um local onde envia uma comunicação entre goroutines e outra goroutines e esse channel pode fazer ele ficar "emperrado", ou seja pedir para ele jogar o valor, para evitar que a aplicação seja encerrada...

```go
// criar channel
wait := make(chan int)

// ...

// Aguardar channel evita que a aplicação seja encerrada
<-wait

// ...

// Fechar / finalizar channel para encerrar a aplicação
close(wait)
```


- Para executar inciamos o server e o client:

```shell
go run cmd/server/server.go
```

```shell
go run cmd/client/client.go
```