# Consul

- Primeiro criar o arquivo `docker-compose.yaml`

- Executar o comando:

```shell
docker-compose up -d
```

- Após isso podemos executar o comando:

```shell
docker-compose ps
```

- Para verificar o docker que subimos

- Para acessar o container executamos o comando:

```shell
docker exec -it consul01 sh
```

**Importante**

- Toda vez que subimos o consul precisamos informar se ele deve subir em modo dev, e server ou client, nesse caso vamos utilizar o modo dev:

```shell
consul agent -dev
```

- Dessa forma ele subirá o agent consul

- Podemos consulta-lo por abrir outro terminal acessar o container novamente e executar o comando:

```shell
consul members
```

- Esse comando é importante pois com ele sabemos quem faz parte do cluster,

- Em produção é recomendado rodar com pelo menos 3 máquinas servers, ou a partir daí sempre números impares

- As máquinas entram em um consenso de qual máquina é a lider, e se o lider cair, as máquinas elegem outra para assumir esse papel.

- De qualquer instancia consul eu consigo verificar os membros do cluster

---

## Acessos

- O Consul possuí API REST, e tem um servidor de DNS imbutido, se eu precisar realizar uma requisição para o Consul, eu utilizo HTTP, ou trabalho com DNS
- Para acessar o catalogo...
- Com o Consul em execução
- Abrimos outro terminal e acessamos o container e executamos o comando:

```shell
curl localhost:8500/v1/catalog/nodes
```

- O servidor DNS do Consul trabalha na porta 8600

- Para podermos utilizar consulta de servidor DNS no Linux Alpine vamos baixar o seguinte apk:

```shell
apk -U add bind-tools
```

- Agora podemos executar o comando `dig` linha de comando que executa a pesquisa de DNS:

```shell
dig @localhost -p 8600
```

- Ele conseguiu bater nesse servidor...
- Porém precisamos informar exatamente o que estamos procurando, que no caso são todos os nodes e o IP dos nodes do Consul que estão registrados:

```shell
dig @localhost -p 8600 consul01.node.consul
```

- Sempre irá terminar .consul no final, faz parte do padrão

- E para pegar apenas o IP utilizamos o comando:

```shell
dig @localhost -p 8600 consul01.node.consul +short
```

---

- Antes de prosseguir realizar o seguinte comando: 

```shell
docker-compose down
```

---

## Clusters

- Ajustar o arquivo docker-compose.yaml, para criar 3 servidores consul,
- Depois subimos os nosso containers:

```shell
docker-compose up -d
```

- Agora vamos subir os 3 servers:

```shell
docker exec -it consulserver01 sh
```

- Se executarmos o comando:

```shell
consul members
```

- Podemos ver que ainda não subiu nada!

- Vamos subir o agente em modo servidor
- Precisamos pegar o ip do nosso container, executando o comando:

```shell
ifconfig
```

- Precisamos criar o diretorio de configuração:

```shell
mkdir /etc/consul.d
```

```shell
mkdir /var/lib/consul
```

- Execute o seguinte comando:

```shell
consul agent -server -bootstrap-expect=3 -node=consulserver01 -bind=IP_OBTIDO_NO_ifconfig -data-dir=/var/lib/consul -config-dir=/etc/consul.d
```

- o `bootstrap-expect=3` informamos a quantidade de servidores esperados para ele se preparar
- `node` o nome do nó, não é obrigatório, se não informar ele pega o hostname
- `bind` informar o ip da máquina em que está execuntando
- `data-dir` local onde ele irá salvar os arquivos
- `config-dir` local onde tem os arquivos de configurações extensão .json ou .hcl(extensão da Hashcorpi)

- E agora se eu executar o `consul members` em outro terminal do mesmo container, ele irá retornar um membro

- Vamos acessar o consulserver02:

```shell
docker exec -it consulserver02 sh
```

- Obtemos o ip e criamos as pastas conforme feito no consulserver01


- Feito isso, mesmo executando o comando `consul members` não será exibido os dois, mas apenas 1 precisamos fazer com que eles se conheçam primeiro
- Agora precisamos fazer um join de um para o outro:

```shell
consul join IP_DO_OUTRO_CONTAINER
```

- Agora sim se eu executar o `consul members` ele exibirá os dois servidores

- Vamos por ultimo subir a terceira máquina:

```shell
docker exec -it consulserver03 sh
```

- Executar os mesmos comandos...

- E podemos realizar o join novamente!
- E ao executar o `consul members` será exibido os 3 servers!

---

### Client

- Adicionamos o client no arquivo `docker-compose.yaml`.
- Para organizar mais criamos a pasta: `clients/consul01` e vamos compartilhar essa pasta para ficar no diretorio de config para o consul poder trabalhar.

- E executamos em outro terminal:

```shell
docker-compose up -d
```

- Depois de subir o client executamos o comando:

```shell
docker exec -it consulclient01 sh
```

- Precisamos ver qual o ip do nosso client utilizando o ifconfig.

- Precisamos criar as pastas:

```shell
mkdir /etc/consul.d
```

```shell
mkdir /var/lib/consul
```
- Dentro desse container iremos executar um comando semelhante aos demais porém iremos criar utilizar para criar o client, quando omitimos o -server do agent ele entende que se trata de client:

```shell
consul agent -bind=IP_OBTIDO -data-dir=/var/lib/consul -config-dir=/etc/consul.d
```

- Após executar o comando ele irá reclamar que não encontrou nenhum consul...
- Abrindo uma nova aba do terminal e acessando o container do client e executando o comando:

```shell
consul members
```

- Verificamos que ele está sozinho
- para juntar aos cluster executamos o comando:

```shell
consul join IP_DE_UM_DOS_SERVERS
```

- E se executarmos o comando:

```shell
consul members
```

- Verificamos que aparece os nossos servers e o client!!!

- E voltando na aba anterior onde executamos o coamndo para subir o client ele irá informar que conseguiu sincronizar!!

----

### Criar o serviço

- Na pasta `clients/consul01` vamos criar um arquivo `qualquer_coisa.json`, no nosso caso chamaremos de `service.json`

- Temos que ter em mente o seguinte precisamos sempre ter o agente na mesma máquina do serviço

- Depois que eu faço isso eu preciso acessar o client ou qualquer servidor e pedir para o consul dar um reload com o comando:

```shell
consul reload
```

- Agora todos os consuls sabem que eu tenho o nginx rodando, agora posso executar o comando:

```shell
apk -U add bind-tools
```

- Agora posso executar:

```shell
dig @localhost -p 8600 SRV
```

- Percebemos que no dns não está trazendo nada...

- E executo o comando:

```shell
dig @localhost -p 8600 nginx.service.consul
```

- Com isso eu obtenho algo semelhante a isso como resultado do comando anterior:

```shell
;; ANSWER SECTION:
nginx.service.consul.   0       IN      A       172.31.0.4
```

- Todos os consuls sabem que tenho esse serviço rodando

- Eu posso executar o comando em outro servidor dentro dos nossos containers:

```shell
apk -U add bind-tools
```

- para instalar o dig

- E por fim executar o comando:

```shell
dig @localhost -p 8600 nginx.service.consul
```

- Veremos que ele traz o serviço do nginx também!!

---

### Realizar consulta no catalogo

- Podemos buscar pelos serviços através do catalogo de serviços, para isso podemos utilizar o comando, em qualquer container do nosso serviço consul:

```shell
curl localhost:8500/v1/catalog/services
```

- Ou também posso utilizar o comando:

```shell
consul catalog nodes -service nginx
```

- E será exibido em qual node que o serviço está registrado

- E para obter detalhes executo o comando:

```shell
consul catalog nodes -detailed
```

- Na documentação do consul temos mais detalhes de buscas que podemos realizar no catalog

- Podemos consultar utilizando a tag também que foi a que adicionamos no `clients/consul01/service.json` `web`:

```shell
dig @localhost -p 8600 web.nginx.service.consul
```

---

### Criar outro client

- Vamos adicionar no `docker-compose` o consulclient02
- Vamos duplicar a pasta `clients/consul01` para `clients/consul02` e ajustar apenas o service.json do consul02 o id para `nginx2`

- Executamos o comando up do docker para subir esse novo container:

```shell
docker-compose up -d
```

- Acessamos o container consul02:

```shell
docker exec -it consulclient02 sh
```

- E basicamente realizamos o mesmo procedimento feito no consulclient01

- E executamos o comando:

```shell
consul agent -bind=$(hostname -i) -data-dir=/var/lib/consul -config-dir=/etc/consul.d -retry-join=IP_DE_UM_SERVER -retry-join=OUTRO_IP
```

- Precisamos abrir mais um terminal para acessar o consulclient02 para executar o seguinte comando:

```shell
apk -U add bind-tools
```

- E por fim podemos fazer a consulta dos nossos serviços:

```shell
dig @localhost -p 8600 nginx.service.consul
```

- Resultand em algo assim:

```shell
;; ANSWER SECTION:
nginx.service.consul.   0       IN      A       172.31.0.4
nginx.service.consul.   0       IN      A       172.31.0.6
```

- Isso é service discovery!!!

- Conforme as máquinas do nginx forem escalando ele vai entrando no consul  e por dns nesse caso já sabemos as nossas máquinas!
- E posso realizar load balance com elas

- Eu não preciso acessar o consul para conhecer as máquinas, eu posso fazer comunicação entre os clients!


---

## Healtcheks

- O Consul irá listar todos os serviços e precisamos saber se eles estão funcionando ou não, para que o consul quando consultado quais serviços disponiveis não exiba aqueles que estão fora!

- Mais detalhes de como o consul funciona com os checks: [Checks](https://www.consul.io/docs/discovery/checks), há diversas formas de verificar se o serviço está funcionando ou não, http, grpc, tcp, ttl...

- Para começar vamos ajustar o arquivo `clients/consul01/service.js`

- Agora acessamos o `consulclient01`:

```shell
docker exec -it consulclient01 sh
```

- E damos um relaod:

```shell
consul reload
```

- O que temos? o um client com a configuração de check e outro não,
- O que não tem check, sempre estará ativo

- E se eu executar o dig:

```shell
dig @localhost -p 8600 nginx.service.consul
```

- Será exibido apenas o serviço sem check!

- Para testar vamos subir um nginx no consulclient01:

```shell
apk add nginx
```

- Precisamos criar a pasta:

```shell
mkdir /run/nginx/
```

- E podemos subir o nginx:

```shell
nginx
```

- Em seguida podemos verificar que ele subiu executando o comando: 

```shell
ps
```

- Vamos criar uma página para o nginx não retornar 404.. criamos a pasta:

```shell
mkdir /usr/share/nginx/html -p
```

- o `-p` é para criar recursivamente caso não exista

- Vamos editar o arquivo `/etc/nginx/conf.d/default.conf`

- caso o vim não esteja instalado só executar o comando:

```shell
apk update
apk add vim
```

```shell
vim /etc/nginx/conf.d/default.conf
```

- e ajustamos para:

```shell
server {
  listen 80 default_server;
  listen [::]:80 default_server;

  root /usr/share/nginx/html;

  location = /404.html {
    internal;
  }
}
```

- E criamos o arquivo:

```shell
vim /usr/share/nginx/html/index.html
```

- Reinciamos o nginx:

```shell
nginx -s reload
```

- E executamos o comando:

```shell
curl localhost
```

- Para verificar que deu certinho a alteração que fizemos no nginx.

- Por fim se executarmos o comando do dig:

```shell
dig @localhost -p 8600 nginx.service.consul
```

- Vemos que o serviço voltou! pois agora o nginx está rodando!!!

- Isso é muito bom pois só serão exibidos os serviços no catalogo de serviços os que estão executando!

---

### Utilizando menos comandos e código!

- Para subirmos os servers e deixar tudo funcionando foi necessário utilizar vários comandos,

- Para deixarmos o processo mais automático faremos o seguinte:

- Criar o arquivo `servers/server01/server.json`:

```shell
{
  "server": true,
  "bind_addr": "172.31.0.3",
  "bootstrap_expect": 3,
  "data_dir": "/tmp",
  "retry_join": ["172.31.0.5", "172.31.0.2"]
}
```

- O `"server": true` - é para informar que o bind será para server, se não ele utiliza como client default
- E no arquivo `docker-compose.yaml` vamos atachar um volume para o container correspondente:

```yaml
consulserver01:
    image: consul:1.10
    container_name: consulserver01
    hostname: consulserver01
    command: ['tail', '-f', '/dev/null'] # Manter o processo rodando
    volumes:
      - ./servers/server01:/etc/consul.d
```

- E executamos o comando:

```shell
docker-compose up -d
```

- Acessamos o container consulserver01:

```shell
docker exec -it consulserver01 sh
```

- E dentro dele só precisamos executar o comando:

```shell
consul agent --config-dir=/etc/consul.d
```

- Que será utilizado o arquivo json para subir as informações e sincronizar

- Faremos o mesmo com o consulserver2 e 3, de criar a pasta server02 e server03 respectivamente e ajustar no arquivo docker-compose.yaml

- E executamos o comando:

```shell
docker-compose up -d
```

- Acessar os containers consulserver2 e 3 e executar o comando:

```shell
consul agent --config-dir=/etc/consul.d
```

---

## Segurança no Consul

- Precisamos adicionar uma forma de segurança para nossos consuls, para que não entre ninguem "estranho" na rede dos consuls

- Para utilizar cryptografia propria do consul ele possuí um parametro chamado de `encrypt`, o qual podemos inserir nos arquivos `server.json` que estão em `servers/server0NUMBER/`,

- E podemos gerar essa chave através do consul, utilizando o comando:

```shell
consul keygen
```

- Ele irá gerar uma chave e precisamos utiliza-la em todos os servers

- Antes de testar tudo podemos executar o seguinte comando em todos os containers:

```shell
rm -rf /tmp/*
```

- Para limpar quaisquer coisas que possam atrapalhar como se fosse um cache

- E subimos o consul utilizando o comando:

```shell
consul agent --config-dir=/etc/consul.d
```

- Para verificar se tem criptografia, podemos escutar determinada porta utilizando o tcp, e ver se conseguimos ver as informações sendo trafegadas para isso precisamos instalar o seguinte:

```shell
apk -U add bind-tools
apk add tcpdump
```

- Instalado utilizamos o comando:

```shell
tcpdump -i eth0 -an port 8301 -A
```

- o `eth0` é a interface de rede, podemos obte-la utilizando o comando `ifconfig` e após o port é a porta que queremos escutar

---

### Ativar a user interface

- É necessario habilitar a UI do consul na hora que ele for subir,
- Há duas formas:
  - No client, passar o agent... e adicionar o `-ui`:
    ```shell
    consul agent --config-dir=/etc/consul.d -ui
    ```
    - Porém acontece que estamos no docker e a `ui` é feita para rodar localhost e como estamos no docker e vamos tentar acessar isso de fora do docker e vamos ter que compartilhar uma porta, precisamos fazer com essa ui possa ser acessada de qualquer lugar...daí precisamos adicionar o seguinte:
    ```shell
    consul agent --config-dir=/etc/consul.d -ui -client=0.0.0
    ```
  - A outra forma e passar essa info no arquivo de configuração...
    - Podemos ajustar o arquivo `servers/server01/server.json`:
    ```json
    {
      "server": true,
      "bind_addr": "172.31.0.3",
      "bootstrap_expect": 3,
      "data_dir": "/tmp",
      "node_name": "consulserver01",
      "encrypt": "d/cjAZ8XDDY8Wa5/u3UumvQhU+Zb/pFYnNNOMwWipe8=",
      "client_addr": "0.0.0.0",
      "ui_config": {
        "enabled": true
      }
    }
    ```
    - Dessa forma adicionando a parte:
    ```json
      "client_addr": "0.0.0.0",
      "ui_config": {
        "enabled": true
      }
    ```
    - Habilitamos a ui do consul, porém ela precisa da porta: 8500, dessa forma ajustamos isso no arquivo `docker-compose.yaml`:
    ```yaml
      consulserver01:
        image: consul:1.10
        container_name: consulserver01
        hostname: consulserver01
        command: ['tail', '-f', '/dev/null'] # Manter o processo rodando
        volumes:
          - ./servers/server01:/etc/consul.d
        ports:
          - "8500:8500" # Adicionado porta para o consul ui
    ```

- Por fim precisamos subir o docker-componse novamente para atualizar a adição da porta e acessar o container do server01 subir o agent normalmente no consul:
```shell
consul agent --config-dir=/etc/consul.d
```

- Por fim podemos acessar a ui, por acessar o link: `http://localhost:8500/ui`
- Utilizar essa interface mais para gerenciar, ver os serviços...
- Pois se for mais para verificar as configurações, o melhor é fazer pelos arquivos, pela linha de comando