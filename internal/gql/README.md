
## Generate the GraphQL endpoint code

#### Install dependencies

```shell
go get github.com/99designs/gqlgen     
```

#### Initialize the code generation

```shell
go run github.com/99designs/gqlgen init --schema schema.graphqls
```

#### (Re)generate the models

```shell
go run github.com/99designs/gqlgen generate
```

#### Run the playground

Set environment variables:

```shell
export PROJECT_ID=podops
export GOOGLE_APPLICATION_CREDENTIALS=/Users/turing/devel/workspace/podops/google-credentials.json
export API_ENDPOINT=http://localhost:8080
```

Run the server:

``shell
go run server.go
```

or, in just one line:

```shell
clear && PROJECT_ID=podops GOOGLE_APPLICATION_CREDENTIALS=/Users/turing/devel/workspace/podops/google-credentials.json API_ENDPOINT=http://localhost:8080 go run server.go
```

#### References

* https://gqlgen.com

* https://github.com/99designs/gqlgen
* https://www.freecodecamp.org/news/deep-dive-into-graphql-with-golang-d3e02a429ac3/


github.com/99designs/gqlgen v0.13.0
github.com/vektah/gqlparser/v2 v2.1.0

go get github.com/vektah/gqlparser/v2@v2.1.0
go get github.com/vektah/gqlparser/v2/ast@v2.1.0
