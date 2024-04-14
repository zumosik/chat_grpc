# GRPC Chat
> Chat app using React, Go and GRPC.

<!-- TOC -->
* [Product Name](#product-name)
  * [Installation](#installation)
  * [Usage example](#usage-example)
  * [Release History](#release-history)
  * [Contributing](#contributing)
<!-- TOC -->

![GitHub License](https://img.shields.io/github/license/zumosik/online_game)
[![Go](https://img.shields.io/badge/Go-1.22.2-00ADD8.svg)](https://golang.org/)
<!-- [![React](https://img.shields.io/badge/React-16.13.1-61DAFB.svg)](https://reactjs.org/) -->


Realtime chat application with groups, commands and file uploads. Uses grpc and Go for backend and React for client. 
_(In future)_

## Installation
1. Install [goose](https://github.com/pressly/goose) and [Docker](https://www.docker.com/)
2. Run docker-compose file
```shell
cd server 
docker-compose up
```
3. Run migrations for postgres
```shell
cd server
make postgres_up 
```

## Usage example

In this version you can only send grpc requests (proto files can be found [here](https://github.com/zumosik/grpc_chat_protos)) to server, but in future this will be good chat app.


## Release History

* 0.0.1
    * Work in progress


## Contributing

1. Fork it (<https://github.com/zumosik/chat_grpc/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

