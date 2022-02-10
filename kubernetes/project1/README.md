# Kubernetes

### Projeto em go

- Criamos um simples projeto em go http
- Para executar basta rodar o comando:

```shell
go run server.go
```

- Criamos um arquivo `Dockerfile`:

```Dockerfile
FROM golang:1.15 # utilizando a imagem go
COPY . . # Mandei copiar tudo que está na aplicação para dentro da pasta do container
RUN go build -o server . # Gerar o build e ao realizar o build, será gerado o executavel, e quando for gerado o executavel...
CMD ["./server"] # poderemos executa-lo 
```

- Com isso executamos um docker build:

```shell
docker build -t NOME_USUARIO/NOME_APLICACAO .
```

- Por fim podemos verificar se está funcionando executando o comando:

```shell
docker run --rm -p 80:80 NOME_USUARIO/NOME_APLICACAO
```

- Como já está funcionando podemos dar um push da nossa imagem:

```shell
docker push NOME_USUARIO/NOME_APLICACAO
```

- Antes de executar o push é necessário realizar o login:

```shell
docker login
```

- Criar tag:

```shell
docker tag NOME_IMAGEM_LOCAL NOME_DA_IMAGEM_DESEJADA_PODE_SER_O_MESMO_IMAGEM_LOCAL:TAG
```

-Ex.:

```shell
docker tag carromesa/go-with-kube carromesa/go-with-kube:v2
```

---

### Pods

- Criar o arquivo `k8s/pod.yaml`

-O kubernetes tem uma api, tudo que acontece no kube é feita através da api, em alguns momentos os recursos que formos utilizar tem uma versão diferente da api
- Podemos criar labels, que nos ajuda a filtrar depois no kube
- Não precisa decorar esses arquivos... o importante é entender a ideia...

- Podemos rodar mais containers no pod... mas de forma geral é um container para um pod

- Criado o arquivo, verificamos se estamos no ambiente certo:

```shell
kubectl get nodes
```

- E o comando que é muito utilizado:

```shell
kubectl apply -f k8s/pod.yaml
```

- Basicamente estou pedindo:
  - Aplicar um arquivo de configuração
  - -f de file
  - E o arquivo de configuração


- Para verificar se foi criado utilizamos o comando:

```shell
kubectl get pods
```

```shell
kubectl get pod
```

```shell
kubectl get po
```

- O que ficou ali:

```shell
NAME       READY   STATUS    RESTARTS   AGE
goserver   1/1     Running   0          16m
```

- Ele diz que temos um pod chamado `goserver`
- Termos um pod rodando `1/1`
- Status rodando
- Ele não restartou
- idade desde que ele está rodando


- Para permitir o acesso rápido ao pod uma forma rápida... podemos executar o comando:

```shell
kubectl port-forward pod/goserver 8000:80
```

- O que estamos fazendo, estamos liberando a porta 8000 da nossa máquina para acessar a porta 80 do pod
- e o `pod/goserver` é o nome do pod, considera a segunda parte... 

- Para deletar o pod executamos o comando:

```shell
kubectl delete pod goserver
```

- `goserver` é o nome do nosso pod

- E para verificar se foi apagado mesmo pode utilizar novamente o comando:

```shell
kubectl get pod
```

- Se desse algum problema de o pod ter travado, eu desse algum problema... para ele ser removido
- Ele não seria criado novamente
- Então é um pouco arriscado criarmos apenas um pod e deixa-lo lá rodando, vai funcionar? vai mas se der algum problema
dá ruim para nós... :(

- Em geral trabalhamos um pouco diferente na questão dos pods no kubernetes, pois quando removemos ele, o mesmo não recriado automaticamente, e em caso de problemas, ele para remover em geral fica dificil de recria-lo novamente

## ReplicaSet

- É um objeto que gerencia os pods, eu posso pedir quantas replicas eu quero manter, e quando for removido um dos pods o replicaSet irá cria-lo novamente então o replicaSet ficará sempre monitorando, e se der algum problema ele recria, ele mata e cria de novo

- para isso criamos o arquivo `k8s/replicaset.yaml`

- Dentro desse arquivo temos o selector:

```yaml
spec:
  selector: 
    matchLabels:
      app: goserver
```

- Através dele conseguimos filtrar as labels, e isso é util para encaminharmos o trafego para determinados pods / replicaSets

- Resumindo...

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: goserver
  labels:
    app: goserver
spec:
  selector: 
    matchLabels:
      app: goserver
  replicas: 2
```

- Essa primeira parte é a "configuração" do nosso ReplicaSet,
- Dizemos que:
  - Queremos duas replicas
  - ReplicaSet se chama goserver
  - Tem o label que também é goserver

- Nesse arquivo também é informado as replicas que eu preciso e as replicas do que basicamente..., que é informado em template as configurações do pod:

```yaml
replicas: 2
  template:
    metadata:
      labels:
        app: "goserver"
    spec: # especificacao do container
      containers:
        - name: goserver
          image: "carromesa/go-with-kube"
```

- Feito isso podemos rodar o comando para criar o nosso replicaSet:

```shell
kubectl apply -f k8s/replicaset.yaml
```

- E ele criou o nosso replicaSet

- Se executarmos o comando:

```shell
kubectl get pods --all-namespaces
```

```shell
NAMESPACE            NAME                                         READY   STATUS    RESTARTS   AGE
default              goserver-6xckg                               1/1     Running   0          17m
default              goserver-8d4b5                               1/1     Running   0          17m
```

- Nesse caso no NAME temos o nome do replicaSet + um nome aleatorio, para não conflitar os pods

- Também posso executar o seguinte comando:

```shell
kubectl get replicasets
```

- Para listar as replicas!

- O ReplicaSet gerencia os pods... podemos testar isso por verificar o nome atual dos pods e remover um deles, podemos verificar em seguida que será criado novamente o pod apagado com outro nome!

```shell
kubectl delete pod goserver-6xckg 
```

- Ou seja quando removi um ele automaticamente já criou o outro... pois ele quer sempre garantir
manter 2 funcionando

- E se alteramos o arquivo `k8s/replicaset.yaml` e ajustarmos o nr de replica e rodar o comando novametne:

```shell
kubectl apply -f k8s/replicaset.yaml
```

- Ele cria os pods adicionais!

---

### Alteração na imagem

- Se realizarmos alterações na imagem informada no replicaset ajustarmos o replicaset também para obter a nova imagem, e realizarmos o comando de apply, os pods existentes ainda continuarão rodando a imagem anterior, para que os pods rodem a imagem nova é necessário remove-los

```shell
kubectl delete pod NOME_DO_POD
```

- Para pegar informações a cerca do pod criado podemos utilizar o seguinte comando:

```shell
kubectl describe pod NOME_DO_POD
```

- O fato de eu alterar a imagem informada no replicaset e ela não ser replicada aos pods pode ser um problema, para resolver isso precisamos "matar" todos os pods, para que o replicaset recrie-os com a nova configuração.

- O kubernetes tem um recurso que faz isso que é o Deployment

### Deployment

- Quando é alterado a versão dele, ele remove todos os replicasets e cria tudo novamente com as novas configurações

- E para isso é algo simples...

- Primeiro criamos o arquivo `k8s/deployment.yaml`

- Ele terá basicamente o mesmo conteúdo do `k8s/replicaset.yaml` o que muda é apenas essa parte:

```yaml
kind: Deployment
```

- Para removermos os replicasets existentes utilizamos o seguinte comando:

```shell
kubectl delete replicaset NOME_DO_REPLICASET
```

- Para rodar o deployment utilizamos o seguinte comando:

```shell
kubectl apply -f k8s/deployment.yaml
```

- Podemos executar o comando:

```shell
kubectl get pods
```

- Para verificar se os pods foram criados

- E podemos verificar que temos o seguinte:

```shell
NAME                        READY   STATUS    RESTARTS   AGE
goserver-7dbc694dc5-4nmgf   1/1     Running   0          47s
goserver-7dbc694dc5-5xb6p   1/1     Running   0          47s
```

- Temos esse padrão:

`DEPLOYMENT-REPLICASET-POD`

- Podemos verificar isso por executar o seguinte comando:

```shell
kubectl get deployments
```

- Retorna o nome do deployment

```shell
kubectl get replicasets
```

- Retorna o nome do replicaset.

- Agora quando alterarmos a imagem de deployment em `deployment.yaml`:

```yaml
spec: # especificacao do container
      containers:
        - name: goserver
          image: "NOME_DA_IMAGEM"
```

- E quando rodamos o comando:

```shell
kubectl apply -f k8s/deployment.yaml
```

- Ele remove todos os pods e cria tudo novamente...

- Porém o replicaset não é removido, o antigo só estará "zerado"

---

### Rollout

- Quando precisamos por algum motivo voltar a versão dos pods...

- Para verificar o historico de rollout/revisions dos pods utilizamos o seguinte comando:

```shell
kubectl rollout history deployment NOME_DO_DEPLOYMENT
```

EX.:

```shell
kubectl rollout history deployment goserver
```

- Esse comando irá listar o historico das revisions:

```shell
REVISION  CHANGE-CAUSE
1         <none>
2         <none>
```

- Para realizar o rollout para versão anterior executar o seguinte comando:

```shell
kubectl rollout undo deployment NOME_DO_DEPLOYMENT
```

- O `undo` volta para ultima versão

- Mas se eu desejar voltar para uma versão especifica executo o seguinte comando:

```shell
kubectl rollout undo deployment NOME_DO_DEPLOYMENT --to-revision=NR_DA_REVISAO
```

- Para ver mais detalhes do deployment para ver até um histórico do que foi feito, e de disponibilidade podemos executar o seguinte comando:

```shell
kubectl describe deployment NOME_DO_DEPLOYMENT
```

---

### Criar Service

- Criar o arquivo `k8s/service.yaml`
- O kind será `kind: Service` pois se trata de um service
- o `selector` utilizamos para que o kubernetes identifique os pods que equivalem ao informados, imaginando que tenhamos vários pods, 200 por exemplo, com o selector podemos especificar que queremos apenas o equivalem no que foi informado..,
- Ou seja ele filtra todos os pods que estarão incorporados/associados nesse serviços:

```yaml
selector:
    app: goserver
```

- Nesse caso ele irá pegar todos os pods em que o label `app` for `goserver`, esse é um filtro, é uma forma que eu posso diferenciar um serviço do outro.

- Outro coisa importante é o type..., por padrão temos 4 tipos de services, `ClusterIP`, `NodePort`, `LoadBalance`, `HadlessService`.

- No final ficará assim:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: goserver-service
spec:
  selector:
    app: goserver
  type: ClusterIP
  ports:
  - name: goserver-service # O name é interessante colocar sempre no formato `algumacoisa-algumaoutracoisa`
    port: 80
    protocol: TCP

```

- Na parte de ports não é obrigatório mas é interessante colocar o name, pois um serviço pode export várias portas

- Por fim para colocar para rodar executamos o seguinte  comando:

```shell
kubectl apply -f k8s/service.yaml
```

- E daí ele cria o nosso serviço...

- Para verificar que os serviços foram criados... utilizamos o seguinte comando:

```shell
kubectl get svc
```

- Nós podemos ver que temos pelo menos dois services... um chamado kubernetes, é o service padrão para fazer requisição chamda de api... etc,

- E temos o outro que criamos que tem um IP interno do servidor, dessa forma todos que tiverem dentro do nosso kubernetes, eles irão bater no service e o service irá escolher quais dos pods ele irá acessar, ou eu chamo o próprio ip ou pelo nome do service mesmo pois o kubernetes tem resolução de DNS...

- Imaginando que eu tenha um service para forncer banco de dados eu faço:

```js
root:root@goserver-service
```

- dessa forma o `goserver-service` ele é "visivel",

- Em geral utilizamos pelo nome...

- Para testarmos redirecionamos uma porta da nossa máquina para esse service... dessa forma:

```shell
kubectl port-forward svc/goserver-service 8000:80
```

- `svc/` = Service
- `/goserver-service` = o NOME do serviço
- `8000:80` porta da minha máquina vs porta do pod

- E para testarmos utilizando o comando:

```shell
curl http://localhost:8000
```

- Dessa forma chamamos o service para que ele destribua melhor a carga escolhendo o pod certinho!

----

### Port e targetPort

- No arquivo .yaml podemos informar além da port o targetPort, ou seja, quando for acessada o porta x do meu container será redirecionado para porta y:

```yaml
    ports:
      - name: name-here
        port: 80
        targetPort: 80000
        protocol: TCP
```

- No caso acima estamos informando que tudo que acessar a porta 80 será direcionado para porta 8000!

- Quando não iformamos nada por padrão o targetPort será o mesmo que o port

- Para testar podemos executar o comando:

```shell
kubectl port-forward svc/goserver-service 9000:80
```

- Dessa forma estamos dizendo que ao acessar a porta 9000 da minha máquina ele acesse a porta 80 do container, e quando o container receber o acesso na porta 80, automáticamente ele direciona para porta 8000, por causa do targetPort.

- O `port` é a porta do service e não do container, e com o targetPort, eu escolho qual porta o service irá acessar para direcionar para os pods

- `browser` -9000:80-> `service` -80:8000-> container/pod

### Api do kubernetes

- O `kubectl` nada mais é do que um binário executável, um client, um comand-line interface  que se comunica com a api do kubernetes através de certificados autenticados;

- Essa api do kubernetes pode ser acessada diretamente...

- No nosso caso o kubernetes está em uma rede fechada, para entrar nessa rede do kubernetes, podemos utilizar o kubectl proxy, que gera um proxy na máquina para conseguir acessar o kubernetes:

```shell
kubectl proxy --port=8080
```

- Ao acessar no browser o endereço: `http://localhost:8080` iremos ver todas as urls do kubernetes

---

### NodePort

- O NodePort gera uma porta maior que 30000 e menor que 32767 e libera essa porta em todos os nodes do cluster dessa forma com o ip idependente do node que entrar conseguirá acessar o serviço, normalmente é utilizado para demosntração, fazer um serviço que vai sair do ar... é a forma mais arcaica, será muito raro para utilizar em produção.

- Para utilizar podemos ajustar o arquivo service.yaml:

```yaml
spec:
  selector:
    app: goserver
  type: NodePort
  ports:
  - name: goserver-service
    port: 80
    targetPort: 8000
    protocol: TCP
    nodePort: 30001 # Se nada for informado será gerado automaticamente
```

- Feito isso só executar o comando:

```shell
kubectl apply -f k8s/service.yaml
```

- Podemos ver que foi aplicado por executar o comando:

```shell
kubectl get svc
```

Então se eu pegar o ip de qualquer node utilizando a porta informada nós conseguiremos acessar o serviço

---


### LoadBalancer

- O LoadBalancer ele gera um IP para poder acessar a aplicação de fora, ele é normalmente utilizado quando utiliza um cluster gerenciado.
- Ele automaticamente gera um IP externo e todos que acessarem por esse IP teram acesso a esse servidor 

- Antes de começar... caso precise deletar um service é só executar o comando:

```shell
kubectl delete svc NOME_DO_SERVICE
```

- No arquivo `k8s/service.yaml` alterar o tipo para LoadBalancer:

```yaml
spec:
  selector:
    app: goserver
  type: LoadBalancer
```

- Agora podemos executar o comando:

```shell
kubectl apply -f k8s/service.yaml
```

- Executando o comando:

```shell
kubectl get svc
```

- Podemos ver que o LoadBalancer gera um IP interno o `CLUSTER-IP` e um ip externo o `EXTERNAL-IP` porém o ip externo localmente pode ser que ele não gere :(

---

### Variaveis de ambiente

- Primeiro vamos modificar o arquivo server.go para adicionar algumas variaveis de ambiente:

```go
func Hello(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	age := os.Getenv("AGE")

	fmt.Fprintf(w, "Hello, I`m %s. I`m %s.", name, age)
}
```

- E vamos gerar o build do arquivo:

```shell
docker build -t NOME_USUARIO/NOME_APLICACAO:v2 .
```

- Realizamos um push

```shell
docker push NOME_USUARIO/NOME_APLICACAO:v2
```

- Ajustamos o arquivo `k8s/deployment.yaml`:

```yaml
env:
  - name: NAME
    value: "Rodolfo"
  - name: AGE
    value: "33"
```

- E vamos executar o deploy:

```shell
kubectl apply -f k8s/deployment.yaml
```

- Para testar utilizamos o seguinte comando:

```shell
kubectl port-forward svc/NOME_DO_SERVICE 9000:80
```

- E acessar no browser: http://localhost:9000

### ConfigMap

- Para utilizar o configMap só criar um arquivo `k8s/configmap-env.yaml`, e ajustar no `k8s/deployment.yaml`:

```yaml
spec: # especificacao do container
  containers:
    - name: goserver
      image: "carromesa/go-with-kube:v3"
      env: # envs
        - name: NAME
          valueFrom:
            configMapKeyRef:
              name: goserver-env
              key: NAME
        - name: AGE
          valueFrom:
            configMapKeyRef:
              name: goserver-env
              key: AGE
```

- E toda vez que eu alterar o valor de uma variavel de ambiente eu altero no configmap, porém é importante que só alterar no configmap não é suficiente... é necessário realziar o create do config map e o apply do deployment para aplicar as alterações:

- Criar o config map execute:

```shell
kubectl apply -f k8s/configmap-env.yaml
```

- Realizar o apply do deployment:

```shell
kubectl apply -f k8s/deployment.yaml
```

---

- Existe uma forma mais fácil ainda que é carregar todas as envs do configmap no deployment:

```yaml
spec: # especificacao do container
  containers:
    - name: goserver
      image: "carromesa/go-with-kube:v3"
      envFrom: # envs
        - configMapRef:
            name: goserver-env
```

- Feito isso só realizar o apply do deployment.


----

### Injetar ConfigMap na aplicação

- Digamos que eu tenho o ngnix e tenho um arquivo conf que deve substituir o padrão do ngnix, para não precisar alterar diretamente no ngnix e recriar a imagem toda vez...
podemos utilizar o configMap para ser um arquivo fisico que será injetado no container para substituir as confs padrões, é um recurso muito utilizado em kubernetes.

- Para realizarmos um teste ajustamos o arquivo `server.go`:

```go
func main() {
	http.HandleFunc("/configmap", ConfigMap)
// ...
}
func ConfigMap(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadFile("myfamily/family.txt")

	if err != nil {
		log.Fatalf("Error reading file", err)
	}

	fmt.Fprintf(w, "My Family: %s.", string(data))
}
```

- Realizamos o build:

```shell
docker build -t NOME_USUARIO/NOME_APLICACAO:v5 .
```

- Realizamos um push

```shell
docker push NOME_USUARIO/NOME_APLICACAO:v5
```

- Criamos o arquivo `configmap-family.yaml`

- Realizamos o apply -f k8s/configmap-family.yaml:

```shell
kubectl apply -f k8s/configmap-family.yaml
```

- E realizamos um ajuste no deployment:

```yaml

```

- Quando falamos em injetar configmap, falamos de volume

- no deployment.yaml adicionamos a parte do volume:

```yaml
volumes:
        - name: config
          configMap:
            name: configmap-family
            items:
            - key: members
              path: "family.txt"
```

- Embora no `server.go` o caminho do arquivo seja `myfamily/family.txt` no `deployment.yaml`, estamos apenas mapeando.

- Feito isso vamos "montar" o volume:

```yaml
volumeMounts:
          - mountPath: "/go/myfamily" # aqui é onde eu quero que fique os meus arquivos, ou seja o volumes.configMap.items.path = family.txt vai para essa pasta!
            name: "config" # será o nome do volume.name que está logo abaixo
```

- Feito isso realizamos o apply do deployment:

```shell
kubectl apply -f k8s/configmap-family.yaml
```

```shell
kubectl apply -f k8s/deployment.yaml
```

- podemos verificar se o pod está rodando por executar o comando:

```shell
kubectl get po
```

- Vamos executar o comando:

```shell
kubectl port-forward svc/goserver-service 9000:80
```

- Podemos ver no browser que tudo continua funcionando!

- Caso necessário podemos apagar o deployment:

```shell
kubectl delete deploy goserver
```

- Para acessar o bash de um kubectl é algo bem parecido com o que é feito no docker:

```shell
kubectl exec -it NOME_DO_POD -- bash
```

- Lembrando que para obter o nome dos pods só executar o comando:

```shell
kubectl get pods
```

- Para verificar os logs de um pod podemos executar o comando:

```shell
kubectl logs POD_NAME
```

----

### Secret

- Semelhante ao configmap, porém os dados ficam mais ofuscados, no geral o que iremos utilizar é o secret do tipo opaco

- Para nosso exemplo vamos criar mais uma function dentro de `server.go`: `Secret`

- E depois vamos gerar mais uma nova versão:

```shell
docker build -t NOME_USUARIO/NOME_APLICACAO:v5.2 .
```

- Depois realizamos o push:

```shell
docker push NOME_USUARIO/NOME_APLICACAO:v5.2
```

- Vamos criar o arquivo `k8s/secret.yaml`
- As variaveis deverão estar em base64 para isso basta executar o comando para obter o base64:

```shell
echo "QUALQUER_COISA" | base64
```

- Para fazer o caminho inverso só executar o comando:

```shell
echo "MTIzNDU2Cg==" | base64 --decode
```

- Base64 é padrão do kubernetes, para deixar mais seguro daí exitem outros sistemas que fazem essa parte de segurança

- Para aplicar o secret podemos executar o seguinte comando:

```shell
kubectl apply -f k8s/secret.yaml
```

- No arquivo `k8s/deployment.yaml`, ajustamos a image para pegar a nova imagem gerada...

- E para utilizar o secret ajustamos:

```yaml
envFrom:
  # sem o secret utiliza o configMapRef normalmente
  - configMapRef:
      name: goserver-env
  # se for utilizar o secret utilizar o secretRef abaixo
  - secretRef:
      name: goserver-secret
```

- Por fim executamos o deployment:

```shell
kubectl apply -f k8s/deployment.yaml
```

- Para testar executamos o comando:

```shell
kubectl port-forward svc/goserver-service 9000:80
```

- E acessamos no browser: http://localhost:9000/secret

- E ele deve mostrar as variaveis certinho!

- Tudo aquilo ali no secret virou variavel de ambiente, para verificarmos isso podemos pegar o pod utilizando o comando `kubectl get po`

- E acessa-lo:

```shell
kubectl exec -it NOME_DO_POD -- bash
```

- E dentro dele executar o comando:

```shell
echo $USER
```

e assim obtemos a variavel!


---


## Healt check

- No arquivo `server.go` vamos criar uma rota chamada `healthz`
- Esse nome geralmente é utilizado por padrão para verificar a saúde da aplicação...

```go
func Healtz(w http.ResponseWriter, r *http.Request) {

	duration := time.Since(startedAt)

	if duration.Seconds() > 25 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
```

- Basicamente criamos uma função para testes,  para quando a aplicação estiver rodando por mais de 25 segundos ele de error 500, o qual iremos tratar com Liveness Probres, posteriormente...

- Feito isso vamos gerar um novo build dessa nossa aplicação:

```shell
docker build -t NOME_USUARIO/NOME_APLICACAO:v5.3 .
```

- Depois realizamos o push:

```shell
docker push NOME_USUARIO/NOME_APLICACAO:v5.3
```

- Por fim ajustamos a versão no `k8s/deployment.yaml`:

```shell
kubectl apply -f k8s/deployment.yaml
```

- E podemos ver que subiu através do comando:

```shell
kubectl get po
```

- E para testarmos podemos executar o comando:

```shell
kubectl port-forward svc/NOME_DO_SERVICE 9000:80
```

- E executar no browser: `localhost:9000/healthz`

---

### Liveness Probres

- para que o kubernetes identifique que há algum problema no healtcheck e realize uma ação acerca disso utilizamos o liveness, para isso ajustamos o arquivo `deployment.yaml`:

```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8000 # colocamos a porta do container e não da aplicação!!!
  periodSeconds: 5 # tempo em segundos testar de quanto em quanto tempo
  failureThreshold: 1 # quantas vezes pode falhar antes de reiniciar
  timeoutSeconds: 1 # Testando de forma integrada, colocar mais tempo, testando apenas o sistema pode ser um timeout menor
  successThreshold: 1 # quantas vezes tem que testar para dizer que a aplicação está ok
```

- Temos 3 formas de comunicação, HTTP, CLI e TCP

- Veja os comentários para entender o código acima!

- Feito isso agora vamos testar:

```shell
kubectl delete deploy goserver
```

```shell
kubectl apply -f k8s/deployment.yaml && watch -n1 kubectl get pods
```

- O diferente nesse comando é essa parte:

```shell
 && watch -n1 kubectl get pods
```

- Nesse caso estamos pedindo para monitorar os pods!

- Caso o comando `watch` não exista precisamos instalar via homebrew:

```shell
brew install watch
```

----

### Obter historico do pod

- Obter os pods:

```shell
kubectl get po
```

- copiar o name do pod que deseja obter o historico...

- E executar o comando:

```shell
kubectl describe pod NOME_DO_POD
```

- E então podemos ver os historicos referente ao pod.

---

### Readiness

- Verificar se a aplicação está 100% pronta para ser utilizada...
- Se o banco de dados está ok, se pode receber trafego... etc...

- Para realizar o teste vamos modificar o server.go:

```go
if duration.Seconds() < 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
```

- Estamos simulando que a aplicação demora 10 segundos para subir

- Feito isso vamos gerar uma nova imagem e subir:

```shell
docker build -t carromesa/go-with-kube:v5.4 .
```

- Depois subimos a imagem:

```shell
docker push carromesa/go-with-kube:v5.4  
```

- Feito isso modificamos o `k8s/deployment.yaml`,

- Testamos utilizando o comando:

```shell
kubectl apply -f k8s/deployment.yaml && watch -n1 kubectl get pods
```

- Porém para permitir que só seja enviado trafego para aplicação apenas quando ela estiver pronta precisamos utilizar o readiness..

- No arquivo `k8s/deployment.yaml` vamos ajustar para adicionar o seguinte:

```yaml
readinessProbe:
    httpGet:
      path: /healthz
      port: 8000
    periodSeconds: 5 # tempo em segundos testar de quanto em quanto tempo
    failureThreshold: 1 # quantas vezes pode falhar antes de reiniciar
    timeoutSeconds: 1 # Testando de forma integrada, colocar mais tempo, testando apenas o sistema pode ser um timeout menor
    successThreshold: 1 # quantas vezes tem que testar para dizer que a aplicação está ok
```

- Agora vamos criar esse pod...

```shell
kubectl apply -f k8s/deployment.yaml && watch -n1 kubectl get pods
```

- A aplicação só será liberado o trafego apenas quando o readiness garantir que a aplicação está ok!

- Outro parametro que podemos colocar no deployment.yaml no readiness e no liveness é o `initialDelaySeconds`, onde eu posso informa o tempo de espera antes de começar a fazer o primeiro teste

### Liveness e Readiness

- Podemos ativar ambos no arquivo `deployment.yaml`

- Porém é bom ser cuidadoso ao trabalhar com ambos, pq dependendo da configuração o liveness pode não permitir que o readiness fique online...

- Quando utilizamos o liveness e o Readiness, em algum momento o readiness precisa ficar pronto, e o livenes precisa aguardar que isso ocorra se não quando o readiness estiver para ficar pronto o liveness vai lá e mata o processo e fica nesse loop infinito

- Algo que podemos fazer é adicionar para o liveness um initialDelay para ele inciar só quando o readiness estiver pronto

- Outro ponto importante é que o readiness ele não verifica apenas na inicialização do container, ele verifica o tempo todo, ele irá ficar de tempo em tempo conforme definido verificando se o container está read, ou seja mesmo depois de o container estiver depois de um tempo rodando, ele está realizando a verificação, pois ele quer ver se está read, ele não quer ver se está live ele quer ver ser está read, se não estiver read ele irá desviar o trafego, ele desvia o trafego enquanto o liveness tenta reiniciar, então...

- O Readiness tira o trafego fora o liveness ele recria o processo.

- Vamos simular por colocar mais uma condição de teste na aplicação ele irá dar erro após 30 segundos que a aplicação estiver no ar, para isso vamos ajustar o arquivo `server.go` adicionamos mais uma condição no if:

```go
if duration.Seconds() < 10 || duration.Seconds() > 30
```

- Geramos uma nova imagem:

```shell
docker build -t carromesa/go-with-kube:v5.5 .    
```

- E o push:

```shell
docker push carromesa/go-with-kube:v5.5
```

- se mantermos da forma como está o `initialDelaySeconds` em 10 segundos para cada um pode ser tenhamos o problema de que o readiness não irá conseguir ficar pronto pois o liveness não irá deixar... ele deve reiniciar antes..., para contornar isso podemos aumentar o delay do liveness em 15 segundos, e ajustar também a imagem do docker para a versão 5.5, e executar o comando:

```shell
kubectl apply -f k8s/deployment.yaml
```

---

### startupProp

- Para resolver os problemas listados pelo readiness e o liveness a partir da versão 1.16 do kubernetes foi adicionado o startupProp,

- Ele funciona como o readiness porém apenas no processo de inicialização, quando ficar pronto, e quando estiver pronto e ai sim que o liveness e o readiness irá atuar

- Adicionamos o startupProp no deployment.yaml:

```yaml
startupProbe:
  httpGet:
    path: /healthz
    port: 8000
  periodSeconds: 3 # tempo em segundos testar de quanto em quanto tempo
  failureThreshold: 30 # quantas vezes pode falhar antes de reiniciar
```

- Com isso podemos até remover a prop `initialDelaySeconds` dos demais pois toda vez que o pod for iniciado ou reiniciado/recriado o startupProbe será executado:

- Por fim podemos executar o comando:

```shell
kubectl apply -f k8s/deployment.yaml
```



----

## Resources e HPA

- Antes de nossa aplicação ir para produção precisamos ter em mente, quantos pods eu posso ter no meu cluster, quantos sistemas eu posso ter no mesmo cluster
- E também outro ponto importante é ter uma forma de "medir" a aplicação, para isso utilizamos o metrics-server.

- Quando trabalhamos em cloud o metrics-server já vem por padrão com o gks, do gcp, o eks do aws e o aks do azure, porém no kind não vem

- Repositorio do kubernetes metric-server: https://github.com/kubernetes-sigs/metrics-server

- no repositório para utilizar o metric-server ele pede para executar o comando:

```shell
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

- Porém como tem as questões de TLS etc... vamos fazer de uma forma diferente...

- Acessamos a pasta k8s, e executamos o comando:

```shell
wget https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

- será criado o arquivo k8s/components.yaml e ajustamos para metrics-server.yaml

- No arquivo na parte de deployment:

```yaml
kind: Deployment
```

ajustamos os args do containers:

adicionamos:

```yaml
- --kublet-insecure-tls
```

- E com isso ele irá permitir trabalhar de modo inseguro sem tls

- Por fim executamos o comando:

```shell
kubectl apply -f k8s/metrics-server.yaml
```

- Para verificar se está funcionando executamos o comando:

```shell
kubectl get apiservices
```

- Irá mostrar todos os serviços que temos disponível e um deles é o do metric-server:

```shell
v1beta1.metrics.k8s.io 
```

- Se o Available estiver true então está tudo certo!

----

## Recursos do sistema

- Definir o recurso necessário que para que o serviço rode e os limites também

- Adicionamos isso no arquivo `k8s/deployment.yaml`:

```yaml
resources:
  requests: # => MINIMO -> SEQUESTRANDO / RESERVANDO os recursos para o POD
    cpu: 100m
    memory: 20Mi
```

- A unidade de medida para CPU é em milicores, por exemplo uma vCPU de 1 significa que ele tem 1000m e se o meu sistema precisa utilizar metade podemos definir 500m! ou podemos colocar em decimal
se o meu sistema utiliza metade definimos que ele utiliza 0.5 = 50%

- Se eu colocar números absurdos a aplicação ficará pendente até que os minimos sejam atendidos;

- Outro parametro é o `limits` ou seja até quanto que ele pode usar no máximo.

- no caso da cpu o ideal é que o limit não ultrapasse a soma de pods que utilizaram a cpu

```yaml
limits:
  cpu: 500m
  memory: 25Mi
```

- Já o caso de memória é limitado

- Feito isso podemos realizar o deployment:

```shell
kubectl apply -f deployment.yaml
```

- Para acompanhar como está o consumo podemos executar o comando:

```shell
kubectl top pod NOME_DO_POD
```

---

### HPA (Horizontal Pod AutoScaling)

- Não usa apenas a cpu como de escala, é possível utilizar outras metricas e customizadas, na maioria das vezes o HPA de CPU de funcionar.

- Para começar criamos o arquivo `k8s/hpa.yaml`

- Essa parte informamos a especificação:

```yaml
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: goserver
  minReplicas: 1
  maxReplicas: 5
  targetCPUUtilizationPercentage: 30
```

- Vamos nos basear o kind sempre no Deployment e não no pod, pois ele que é o responsável em criar os pods e "gerenciar"
- em geral o ideal no minReplicas é de 2 ou mais
- o maxReplicas é necessário informar se não ele irá tentar adicionar pods de forma infinita
- e o target fazemos para cpu quando chegar em uma porcentagem.

- Feito isso podemos rodar o comando:

```shell
kubectl apply -f k8s/hpa.yaml
```

- Podemos verificar o HPA utilizando o comando:

```shell
kubectl get hpa
```

---

### Teste de stress com fortio

- Para isso foi ajustado a função: Healtz do server.go:

```go
func Healtz(w http.ResponseWriter, r *http.Request) {

	duration := time.Since(startedAt)

	// Testes de estress
	if duration.Seconds() < 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
```

- E realizado o build da nova imagem:

```shell
docker build -t carromesa/go-with-kube:v5.6 .
```

- E o push:

```shell
docker push carromesa/go-with-kube:v5.6 
```

- E ajustado no arquivo `k8s/deployment.yaml`:

```yaml
containers:
  - name: goserver
    image: "carromesa/go-with-kube:v5.6"
```

- E executado o apply:

```shell
kubectl apply -f k8s/deployment.yaml
```

---

### Fortio

- Repositório: https://github.com/fortio/fortio

- Ferramenta desenvolvida em go ajuda a criar Teste de stress

- Para iniciar vamos executar o seguinte comando:

```shell
kubectl run -it fortion --rm --image=fortio/fortio -- load -qps 800 -t 120s -c 70 "http://goserver-service/healthz"
```

- O comando é parecido com um comando para acessar uma imagem do docker!

- O diferente é a parte do laod, ou seja carregar, -qps = queries por segundo,o -t = segundos de duração, o -c = a conexões simultaneas ou seja quantas threads vão acessar para gerar as 800 queries.

- Na url passamos qual é o service que estamos acessando, conforme o name especificado no arquivo `k8s/service.yaml` na prop `name` e adicionamos a `/` para especificar qual recurso no sistema que iremos testar nesse caso iremos testar o healthz conforme especificado no `server.go`:

```go
func main() {
	http.HandleFunc("/healthz", Healtz)
  // ...
}
```

- E em outra aba do terminal executamos o comando para monitorar:

```shell
watch -n1 kubectl get hpa
```

- que irá acompanhar a cada um segundo


----

## Criar volume persistente

- Criar uma arquivo `k8s/pvc.yaml` (pvc = persiste volume claim)

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: goserver-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```

- Com isso estavmos solicitando um volume persistente de 5gb com acesso de leitura e escrita

- Para aplicar isso utilizamos o comando:

```shell
kubectl apply -f k8s/pvc.yaml
```

- Se executarmos o comando:

```shell
kubectl get pvc
```

- Vemos que está pendente, isso pq ele espera realizar o bind ou uma conexão para poder ficar "ativo".

- Para isso vamos ajustar o `k8s/deployment.yaml`:

- Debaixo de volumes adicionamos:

```yaml
volumes:
  - name: goserver-volume # Aqui pode ser qualquer nome
    persistentVolumeClaim:
      claimName: goserver-pvc # Esse nome aqui precisa ser o mesmo definido no arquivo pvc.yaml
```

- E para montar esse volume precisamos chamar esse volumes.name, no caso `goserver-volume` para montar o volume, podemos informa isso debaixo de `volumeMounts`:

```yaml
volumeMounts:
  - mountPath: "/go/pvc"
    name: "goserver-volume"
```

- Feito isso atualizamos o deployment:

```shell
kubectl apply -f k8s/deployment.yaml
```

- Agora se formos verificar como estar o `pvc` que criamos estará diferente:

```shell
kubectl get pvc
```

- Retorna isso:

```shell
NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
goserver-pvc   Bound    pvc-15a19d62-4e88-4452-bad5-2c0a9e26d79f   5Gi        RWO            standard       11m
```

- O status dele está como Bound ou seja foi feito o bind do volume

- Para realizarmos um teste vamos entrar em um pod:

```shell
kubectl get po
```

- Pegar um id de um pod e executar o comando:

```shell
kubectl exec -it ID_DO_POD -- bash  
```

- Acessar a pasta pvc `cd /go/pvc/` e criar um arquivo:

```shell
touch oi
```

- E vamos apagar esse pod:

```shell
kubectl delete pod ID_DO_POD
```

- O arquivo não será perdido podemos pegar o novo pod que foi gerado e executar os mesmos comandos o arquivo deve estar lá!

- Um ponto importante é que o tipo de acesso está como `ReadWriteOnce` ou seja outros pods que estiverem em outros nodes não irão necessariamente conseguir acessar

---

## Criando StatefulSet

- Uma introdução sobre isso pode ser encontrada aqui [StatefulSet](https://enormous-platinum-684.notion.site/StatefulSet-1289f6c3f577461f88410f840da8d967)

- Primeiro vamos ver um problema que pode ocorrer quando utiliza volumes com pods random.

- Vamos criar o arquivo `k8s/statefulset.yaml`

- Nesse primeiro momento o kind dele vamos colocar como `Deployment`

- E executamos o comando:

```shell
kubectl apply -f k8s/statefulset.yaml
```

- Executando o comando: 

```shell
kubectl get po
```

- Eu terei algo como isso:

```shell
mysql-69c59457b7-7bxc2      1/1     Running            0          21s
mysql-69c59457b7-9sr2g      1/1     Running            0          25s
mysql-69c59457b7-b5skz      1/1     Running            0          29s
```

- Porém não consigo definir quem é o master, então esse é o problema de utilizar o Deployment para lidar com volumes.

- Então vamos remove-los:

```shell
kubectl delete deploy mysql
```

- E no arquivo `k8s/statefulset.yaml` alteramos o kind de Deployment para StatefulSet e no spec dele precisamos adicionar o serviceName nesse caso com final -h pois ele será `headless`:

```yaml
spec:
  serviceName: mysql-h
```

- Feito isso iremos executar o comando:

```shell
kubectl apply -f k8s/statefulset.yaml
```

- Dessa forma ele cria os pods com esses nomes:

```shell
NAME                        READY   STATUS              RESTARTS   AGE
mysql-0                     1/1     Running             0          19s
mysql-1                     1/1     Running             0          15s
mysql-2                     0/1     ContainerCreating   0          3s
```

- Se precisarmos escalar fica bem mais simples,

- Ou seja se aumentarmos a quantidade de replica ele irá pegar sempre do último e criando os seguintes um de cada vez

- E realizando um downsize ou seja diminuir o número de replicas, ele irá remover uma por uma na ordem contrária

- Porém caso seja necessário criar os pods de forma paralela só adicionar no spec do arquivo yaml a propriedade `podManagementPolicy` com valor `Parallel` e remover o número de replicas

- Antes de aplicar é necessário remover o anterior:

```shell
kubectl delete statefulset mysql
```

- Depois só aplicar:

```shell
kubectl apply -f k8s/statefulset.yaml
```

- Para escalar também podemos utilizar via cli:

```shell
kubectl scale statefulset  mysql --replicas=5
```

- E será criado todos em paralelo

---

## Criando headless service

- Uma introdução sobre isso pode ser encontrada aqui [StatefulSet](https://enormous-platinum-684.notion.site/StatefulSet-1289f6c3f577461f88410f840da8d967)

- Iremos criar o arquivo `mysql-service-h.yaml` o `h` é de headless:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-h
spec:
  selector:
    app: mysql
  ports:
    - port: 3306
  clusterIP: None
```

- a prop `clusterIP: None` informamos que ele não irá utilizar IP mas sim o service name, ou seja será via DNS, E como ele irá resolver ou saber esse nome?
O name que utilizamos no `k8s/statefulset.yaml`, que no caso foi `serviceName: mysql-h` precisa ser o mesmo em `k8s/mysql-service-h.yaml` em `metadata.name`,

- Quando eu criar o `statefulset.yaml`, será criado o `mysql-h` e quando eu for criar o `mysql-service-h.yaml` que será `clusterIP: None`, o kubernetes irá resolver daí via dns utilizando o `metadata.name`.

- No arquivo `statefulset.yaml` vamos adicionar 3 replicas: 

```yaml
spec:
  replicas: 4
```

- Para garantir vamos remover o statefulset criado anteriormente executando o comando:

```shell
kubectl delete statefulset mysql
```

- E executar o comando:

```shell
kubectl apply -f k8s/statefulset.yaml
```

E o comando:

```shell
kubectl apply -f k8s/mysql-service-h.yaml
```

- Executando o comando:

```shell
kubectl get po
```

- Obtemos o seguinte resultado:

```shell
mysql-0                     1/1     Running   0          21s
mysql-1                     1/1     Running   0          21s
mysql-2                     1/1     Running   0          21s
mysql-3                     1/1     Running   0          20s
```

- Então quando alguém precisar gravar alguma informação direcionamos para o mysql-0

- Executando o comanod:

```shell
kubectl get svc
```

- Obtemos o seguinte:

```shell
mysql              ClusterIP      None            <none>        3306/TCP       9m24s
```


- Para verificar se o mysql está acessivel via nome podemos realizar um ping dentro de um pod:

```shell
kubectl exec -it NOME_DO_POD --bash
```

- E damos um ping no mysql-h:

```shell
ping mysql-h
```

- E para pingar um especifico utilizamos o comando:

```shell
ping mysql-0.mysql-h
```

- Ou seja dentro de um serviço conseguimos chamar o mysql com base no nome.


---

## Criar volumes dinamicamente com statefulset

- A ideia é criar um volume persistente para cada replica para isso iremos ajustar o arquivo `statefulset.yaml`, toda vez que subir uma replica ele irá aplicar o template:

```yaml
# statefulset.yaml
  volumeClaimTemplates:
  - metadata:
      name: mysql-volume
    spec:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 5Gi
```

Para o nosso mysql utilizar o volume adicionamos essa prop:

```yaml
# statefulset.yaml
  volumeMounts:
    - mountPath: /var/lib/mysql
      name: mysql-volume
```

- Para testar inicialmente removeremos o statefulset criado anteriormente com o comando:

```shell
kubectl delete statefulset mysql
```
- E vamos cria-lo novamente:

```shell
kubectl apply -f k8s/statefulset.yaml
```

- Para verificar se foram criados os volumes podemos utilizar o seguinte comando:

```shell
kubectl get pvc
```

- pvc = persistent volume claim

- Com o comando acima devemos obter algo assim:

```shell
NAME                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
mysql-volume-mysql-0   Bound    pvc-44ab1f5f-c500-4252-b860-b407c6423400   5Gi        RWO            standard       109s
mysql-volume-mysql-1   Bound    pvc-270a7133-ca64-4118-b391-8df6ea418db3   5Gi        RWO            standard       66s
mysql-volume-mysql-2   Bound    pvc-f229edfe-0d0f-425e-a016-15449786c9f4   5Gi        RWO            standard       66s
mysql-volume-mysql-3   Bound    pvc-4ed04809-068e-4c5e-b810-b0673730c6b1   5Gi        RWO            standard       66s
```

- E caso removamos um pod, o kubernets irá recria-lo e atachar o volume! Pois o volume não será removido!

---

## Banco de dados no Kubernetes?

- Isso é complexo...

- Para uma aplicação critica que precisa tunar e melhorar o banco de dados

- Então se for uma aplicação pequena que não irá crescer, sim até tudo bem!

- Mas para aplicação criticas, escolha serviços gerenciaveis RDS da AWS por exemplo, ou outros etc...

- Um exemplo que pode utilizar a base no kubernetes pode ser o wordpress um blog um pouco mais simples, podemos utilizar o kubernetes