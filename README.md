# Notes API

This is the backend for `Notes` - Online notes application.

## How to run?

### Prerequisite

- Make sure that your machine have `GoLang 1.13.8` installed. If it is not installed, please visit below link to install GoLang.
https://golang.org/doc/install
- Make sure that your machine `Docker - 19+` , `Docker Machine` and `Docker compose` installed. Follow the links below if not installed.
https://docs.docker.com/install/
https://docs.docker.com/compose/install/
https://docs.docker.com/v17.12/machine/get-started/#create-a-machine

### Run

- Clone this repository to your local directory.
```
$git clone https://github.com/jeroldleslie/my-notes-backend
$cd my-notes-backend
```
- Run `run_docker_api.sh`

```
$./run_docker_api.sh
```

Your Notes API is up now. For front-end you can access it through the Base URL - `http://<docker machine ip>:8000`.

#### To get docker machine IP
- run command `docker-machine ip default` for default docker machine.

##### Example:

```
#docker-machine ip default
```
you will get

```
#docker-machine ip default
#192.168.xxx.xxx
```



