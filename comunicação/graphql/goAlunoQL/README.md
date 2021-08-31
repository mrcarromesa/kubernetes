# GraphQL com GO

- Para esse projeto utilizaremos essa dependencia: [gqlgen](https://github.com/99designs/gqlgen)
- Documentação dela pode ser encontrada aqui [gqlgen](https://gqlgen.com)

### Iniciando o projeto

- Executar o seguinte comando:

```shell
go mod init github.com/USUARIO/REPOSITORIO
```

- Isso ajudará o go a gerenciar as dependencias
- Irá criar o arquivo go.mod na pasta

---

### Adicionar dependência gqlgen

- Executar o comando:

```shell
go get github.com/99designs/gqlgen
```

- Ele irá baixar, e podemos verificar no arquivo `go.mod` que ele já exibirá a dependência do gqlgen:

```mod
require github.com/99designs/gqlgen v0.13.0 // indirect
```

### Criando o esqueleto do projeto

- Para criar uma estrutura basica especifica para o o `gqlgen` utilize o comando:

```shell
go run github.com/99designs/gqlgen init
```

- Ele gera um esqueleto de modelo

---

### Criando schema

- Vamos criar um arquivo `graph/schema.graphqls`

- Adicionamos os types, que são as estruturas, os input que são as estrutura para o input de dados, as querys para realizar a consulta dos dados e os mutations para inserir/alterar os dados

- Por fim executamos o comando:

```shell
go run github.com/99designs/gqlgen init
```

- Caso não tivesse nada ele geraria o esqueleto

- Feito isso ele criará um arquivo `server.go` para já subirmos o servidor graphql,
- Irá criar os models com base no nosso schema

- Os arquivos que vamos trabalhar serão `graph/resolver.go` e `schema.resolvers.go`

- Em qualquer projeto graphql, nós temos os resolvers, que são conjuntos de forma geral que vai ajudar a resolver o que está sendo pedido pelo client, eles tem a implementação, para buscar os dados seja em uma tabela ou em outra api, para trazer ou inserir o dado, e por fim retornamos apenas as informações solicitadas pelo client.

- O `schema.resolvers.go` utilizamos para implementar a busca dos dados seja ela de uma base de dados, repositories ou api

- Nele também temos os:

```go
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
```

Eu recebo todos os resolvers, e coloco tudo que o client precisa 

---

### Implementando o Resolver

- Ajustamos o arquivo `graph/resolver.go` para montar a estrutura de retornos
- E implementamos a função Categories de `schema.resolvers.go`:

```go
func (r *queryResolver) Categories(ctx context.Context) ([]*model.Category, error) {
	return r.Resolver.Categories, nil
}

func (r *queryResolver) Courses(ctx context.Context) ([]*model.Course, error) {
	return r.Resolver.Courses, nil
}

func (r *queryResolver) Chapters(ctx context.Context) ([]*model.Chapter, error) {
	return r.Resolver.Chapters, nil
}
```

- Para testar executamos o seguinte comando: 

```shell
go run server.go
```

- Depois acessar no browser o endereço que é exibido, no meu caso é o http://localhost:8080/

- Para consultar uma categorie por exemplo podemos utilizar o seguinte comando no playground:

```
query findCategories {
  categories {
    id
    name
    description
    courses {
      name
    }
  }
}

```

- Por enquanto não virá nada pois não temos nada cadastrado\

---

### Mutation

- Vamos realizar alguns ajustes em `schema.resolvers.go`:

```go
func (r *mutationResolver) CreateCategory(ctx context.Context, input model.NewCategory) (*model.Category, error) {
	category := model.Category{
		ID:          fmt.Sprintf("T%d", rand.Int()),
		Name:        input.Name,
		Description: &input.Description,
	}
	r.Categories = append(r.Categories, &category)

	return &category, nil
}
```

- Na função recebemos o paramentro de entrada `input` que utilizamos para definir os valores de `Name` e `Description`, nesse caso da Category

- Para testar executamos o seguinte comando:

```shell
go run server.go
```

- Depois acessar no browser o endereço que é exibido, no meu caso é o http://localhost:8080/

- Ali no playground executamos a seguinte estrutura:

```
mutation createCategory {
  createCategory(input: {name: "PHP", description: "PHP is awsome"}) {
    id
    name
    description
  }
}
```

- Clicar em executar e ele irá retornar o resultado!

- Após isso podemos executar o findCategories:



```
query findCategories {
  categories {
    id
    name
    description
  }
}

```

- E executar que ele trará o resultado!

----

- No caso do model Course, ele precisa receber também a Category, não apenas o id da category mas a Category inteira, para resolver isso precisamos buscar a category antes de inserir juntamente com o course:

```go
func (r *mutationResolver) CreateCourse(ctx context.Context, input model.NewCourse) (*model.Course, error) {

	var category *model.Category

	for _, v := range r.Categories {
		if v.ID == input.CategoryID {
			category = v
		}
	}

	course := model.Course{
		ID:          fmt.Sprintf("T%d", rand.Int()),
		Name:        input.Name,
		Description: &input.Description,
		Category:    category,
	}

	r.Courses = append(r.Courses, &course)

	return &course, nil
}
```

- No caso obtenho o input.CategoryID, e então preciso pecorrer todas as categories em busca da category, e ao encontrar eu posso inserir

---

- Na hora de consultar itens que tem relação não está sendo retornado esses itens relacionados, pois só estamos realizando a relação apenas quando criamos
- Não na hora que buscamos os dados
- E temos que cuidar para que ao fazer isso, não seja feita muitas e muitas consultas na base de dados...

- Para isso vamso criar o arquivo `graph/model/category.go`, `graph/model/chapter.go`, `graph/model/course.go`
- E vamos remover as structs correspondentes de `graph/model/models_gen.go` e adicionar ao model correspondente.

- Vamos realizar os ajustes para não ficar realizando muito aninhamento dos registros
- Em `graph/model/course` vamos remover o `Chapters` pois não vamos precisar disso aqui agora:

```go
Chapters    []*Chapter `json:"chapters"`
```

- Em `graph/model/category` podemos remover o `Courses     []*Course `json:"courses"``, pois só iremos utiliza-lo quando precisar

- Pegar o nome do module, podemos pegar essa info em `go.mod`:

```mod
module github.com/mrcarromesa/graphql
```

- Por fim precisamos adicionar o endereço desses models em `gqlgen.yml`, debaixo de `models:`:

```yml
models:
  Category:
    model: github.com/mrcarromesa/graphql/graph/model.Category
  Course:
    model: github.com/mrcarromesa/graphql/graph/model.Course
  Chapter:
    model: github.com/mrcarromesa/graphql/graph/model.Chapter
```

- Basicamente definimos cada um dos modulos e tiramos o relacionamento automatico, se precisarmos fazer uma relação de listagem precisamos implementar essa relação.

- Pois eu posso precisar pegar uma informação de um banco de dados e outra informação de outra fonte, e depois relacionar elas.

- Agora vamos executar o comando:

```shell
gqlgen generate
```

- ou

```shell
go run github.com/99designs/gqlgen generate
```

- E daí ele deverá gerar os resolvers em `graph/schema.resolvers.go`:

```go
type categoryResolver struct{ *Resolver }
type courseResolver struct{ *Resolver }
```

- Também ele adiciona os metodos:

```go
func (r *categoryResolver) Courses(ctx context.Context, obj *model.Category) ([]*model.Course, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *courseResolver) Chapters(ctx context.Context, obj *model.Course) ([]*model.Chapter, error) {
	panic(fmt.Errorf("not implemented"))
}
```

- Ele recebe a category mais retorna o Course
- recebe o course mas retorna os Chapters

- Para obter os cursos de uma categoria especifica utilizamos isso:

```go
func (r *categoryResolver) Courses(ctx context.Context, obj *model.Category) ([]*model.Course, error) {
	var courses []*model.Course

	for _, v := range r.Resolver.Courses {
		if v.Category.ID == obj.ID {
			courses = append(courses, v)
		}
	}

	return courses, nil
}
```

- Obviamente se fosse para obter de um banco de dados, seria algo como `SELECT * FROM courses WHERE id = obj.ID`

- O mesmo fazemos com o `Chapter`

- Agora com tudo isso implementado podemos executar o servidor:

```shell
go run server.go
```

- E acessamos o playground novamente

- Criamos uma nova categoria:

```
mutation createCategory {
  createCategory(input: {name: "PHP", description: "PHP is awsome"}) {
    id
    name
    description
  }
}
```

- Criamos o curso com base nos comandos em `schema.graphqls`:

```
mutation createCourse {
  createCourse(input: {name: "Evolving with PHP", description: "Mega PHP is awsome", categoryId: ""}) {
    id
    name
    description
    category {
      id
      name
    }
  }
}
```

- E agora podemos consultar as categorias juntamente com os cursos:

```
query findCategories {
  categories {
    id
    name
    description
    courses {
      name
    }
  }
}
```

- E ele deve retornar:

```
{
  "data": {
    "categories": [
      {
        "id": "T5577006791947779410",
        "name": "PHP",
        "description": "PHP is awsome",
        "courses": [
          {
            "name": "Evolving with PHP"
          }
        ]
      }
    ]
  }
}
```

- Isso só foi possivel pois implementamos isso em `schema.resolvers.go`:

```go
func (r *categoryResolver) Courses(ctx context.Context, obj *model.Category) ([]*model.Course, error) {
	var courses []*model.Course

	for _, v := range r.Resolver.Courses {
		if v.Category.ID == obj.ID {
			courses = append(courses, v)
		}
	}

	return courses, nil
}
```

- Que isso busca todos os cursos de cada categoria.

- E se eu executar o seguinte sem o `courses`:

```
query findCategories {
  categories {
    id
    name
    description
  }
}
```

- Ele nem irá executar a função `Courses`.

- Podemos criar um capitulo também:

```
mutation createChapter {
  createChapter(input: {name: "Evolving with PHP",  courseId: ""}) {
    id
    name
    category {
      id
      name
    }
    course {
      id
      name
    }
  }
}
```

- E podemos utilizar o find dessa forma:

```
query findCategories {
  categories {
    id
    name
    description
    courses {
      name
      chapters {
        name
      }
    }
  }
}
```

- Com GraphQL eu posso escolher os dados que eu quero retornar!

----


### IMPORTANTE

**N + 1**

- Pode acontecer de ao realizar uma dada consulta, sejam feitas outras consultas para cada um dos registros retornados na consulta,

- Imagina que eu quero buscar todos os cursos, se eu tivesse por exemplo 2 cursos, ele iria fazer a consulta desse curso 1, iria realizar as consultas internas desse registro, e iria refazer esse processo com o curso 2, 3, 4, ele irá sair criando um monte de sql, pq para cada registro ele irá sair fazendo um monte de consulta.

- Isso pode ser um grande problema de desempenho e peformance, temos que tomar muito cuidado com isso.

- É necessário aplicar debugs, profiles para testar a aplicação para resolver o problema de N + 1, se não a aplicação não irá aguentar, o banco de dados não aguenta muitos usuários executandos consultas e consultas profundas no banco de dados ao mesmo tempo

- A ferramenta utilizada possuí a solução para isso [Dataloaders](https://gqlgen.com/reference/dataloaders/)

- Mesmo assim é interessante utilizar limitador, e não trazer o banco de dados inteiro

- O mesmo se dá com o ORM