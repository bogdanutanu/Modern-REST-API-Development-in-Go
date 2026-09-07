# Chapter 11 - Documenting with OpenAPI

Generating the documentation from comments in the code.

Install the `swag` tool:
```
go install github.com/swaggo/swag/cmd/swag@latest
```

Create the Swagger 2.0 documentation:
```
~/go/bin/swag init
```

Now, having Node installed, we'll use a Node package to generate 
the OpenAPI 3.0 docs from the OpenAPI 2.0:
```
npx -p swagger2openapi swagger2openapi --yaml --outfile docs/openapi.yaml docs/swagger.yaml
```

We can now format the docs in human readable HTML, with tools such as Redoc:
```
npx @redocly/cli build-docs docs/openapi.yaml
```