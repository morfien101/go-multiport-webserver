# What is this

This is a demo server that listens on 2 ports. It is crudly written but is only for demo use.

# What does it do

This service simply starts and listens on 2 ports 8080 and 8443.
I might add certificates later.

# Show me

```sh
go run *.go
8080 - Got reuest on 0.0.0.0:8080.
8443 - Got reuest on 0.0.0.0:8443.
```

```sh
for port in 8080 8443; do curl http://localhost:$port/; done
Hello World! - port 8080
Hello World! - port 8443
```

# WARNING

DEMO CODE. Use at your peril.