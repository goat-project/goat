<img src="https://github.com/goat-project/goat/blob/master/img/goat.png" width="100">

# Goat - GO Accounting Tool

The Goat is a service that is running in the background and waiting for a connection from a compatible client. Once the server receives accounting data, it is transformed into the configured format and writes them to the destination file.
Multiple clients can use the server at once.

![goat-diagram](https://github.com/goat-project/goat/blob/master/img/goat-diagram.png)

See [wiki](https://github.com/goat-project/goat/wiki) for description.

## Requirements
* Go 1.11 or newer.
* The accounting client that connects to a cloud, extracts data (about virtual machines, virtual networks, and storages), filters them accordingly and then sends them to a server for further processing. At the moment, Goat project implemented [Goat-one](https://github.com/goat-project/goat-one) client for [OpenNebula cloud computing platform](https://opennebula.org/).
* The accounting tool that collects data into a central accounting database where it is processed to generate statistical summaries.

## Installation
The recommended way to install this tool is using `go get`:
```
go get -u github.com/onego-project/goat
```

## Configuration
Usage of goat:
```
  -cert-file string
        server certificate file (default "server.pem")
  -debug
        True for debug mode, false otherwise
  -ip-per-file uint
        number of IPs per json file (default 500)
  -key-file string
        server key file (default "server.key")
  -listen-ip string
        IP address to bind to (default "127.0.0.1")
  -log-path string
        log path
  -out-dir string
        output directory
  -port uint
        port to bind to (default 9623)
  -storage-per-file uint
        number of storages per xml file (default 500)
  -templates-dir string
        templates directory
  -tls
        True uses TLS, false uses plaintext TCP
  -vm-per-file uint
        number of VMs per template file (default 500)
```

## Example
Run goat using a template from `templates` directory and writing an output to the `out` directory.
```
go run goat.go -out-dir=out -templates-dir=templates
```

## Container
The goat should run into the container described in [Dockerfile](https://github.com/goat-project/goat/blob/master/Dockerfile). 
Build and run commands:
```
docker image build -t goat-image .
docker run --rm --publish 9623:9623 --network host --name goat --volume goat:/var/goat goat-image
```
See [wiki](https://github.com/goat-project/goat/wiki) for more info.

## Contributing
1. Fork [goat](https://github.com/goat-project/goat/fork)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
