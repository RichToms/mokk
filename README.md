# Mokk

Mokk is a CLI library that allows you to create mock APIs from a config file.
This application was built to improve the process of integrating with third-party APIs during development.

## Installation

You can install this CLI tool if you have Go installed and your GOPATH correctly configured.

```shell
go install github.com/richtoms/mokk@latest
```

## How to use

```shell
mokk generate example
mokk start --port=8080 --config=example.mokk.yml
```
These commands should provide an example config file for you to get started.
If you open another terminal you should now be able to reach the server via cURL:

```shell
curl -i -X "GET" http://localhost:8080/
```

### From Docker

If you are using Docker, you can start your Mokk server without even needing to install the CLI.

It is recommended to create your own config file and provide it to the container, otherwise it will start with its own example Users API.

```shell
docker run \
  -p 8080:80 \
  --volume my-app.mokk.yml:/app/mokk.yml \
  richtoms/mokk:0.1
```

Alternatively, you can use one of the built-in example APIs found in the `./examples` directory.

```shell
docker run \                                                                                            5m 29s
  -p 8080:80 \
  richtoms/mokk:0.1 \
  ./mokk start --config=./examples/aws.apigateway-ws.yml
```

## Configuration

Mokk is config-driven by design, using YAML to be developer-friendly. Below is an example of a Mokk config 
file to help get you started:

```yaml
name: Mokk Example Server
routes:
  - path: "users"
    method: "GET"
    statusCode: 200
    response: '{"status":"Success","users":[{"name":"MockZilla"}]}'

  - path: "users/:id"
    method: "GET"
    statusCode: 200
    response: '{"status":"Success","user":{"name":"MockZilla"}}'
```

### Path

The path field follows [GoFiber's routing patterns](https://docs.gofiber.io/guide/routing#paths) therefore you can utilise wildcards in your paths.

**Note:** it is important when defining a mix of route params and static values to define the static paths first. Due to Fiber's routing system it will struggle to match your static path when a dynamic is defined first.

### Method

The method field *should* be one of the major HTTP verbs:

- `GET`
- `POST`
- `PUT`
- `PATCH`
- `DELETE` 
- `HEAD`
- `OPTIONS`

### Response

Mokk currently only supports JSON APIs, therefore this property must contain some form of JSON array/object.