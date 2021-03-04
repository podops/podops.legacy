
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

#### References

* https://gqlgen.com

* https://github.com/99designs/gqlgen
* https://www.freecodecamp.org/news/deep-dive-into-graphql-with-golang-d3e02a429ac3/

