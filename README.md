# A reply web server for test
The replyserver is a simple web server for test. It returns the client request and its headers.

Usage:

    ./replyserver -h
    -listen string
            The address to listen on for http requests (default ":1969")
    -listentls string
            The address to listen on for https requests (default ":1968")
    -cert string
            The file name containing a TLS server certificate (default "server.crt")
    -key string
            The file name containing a TLS server key (default "server.key")
    -check string
            The address to listen on for healty probes. (default ":1936")

For example:

    ./replyserver -cert server-tls.crt -key server-tls.key

    curl http://macbook.local:1969
    Server Name = "MacBook.local"
    Client Addr = "127.0.0.1:64349"
    Host = "macbook.local:1969"
    Method = "GET"
    URL = "/"
    Protocol = "HTTP/1.1"
    +++ Request Headers: +++
    Header["Accept"] = ["*/*"]
    Header["User-Agent"] = ["curl/7.54.0"]

    curl -k https://macbook.local:1968
    Server Name = "MacBook.local"
    Client Addr = "127.0.0.1:64500"
    Host = "macbook.local:1968"
    Method = "GET"
    URL = "/"
    Protocol = "HTTP/2.0"
    +++ Request Headers: +++
    Header["Accept"] = ["*/*"]
    Header["User-Agent"] = ["curl/7.54.0"]

Build and run with Docker:

    docker build -t bsctl/replyserver:latest .
    docker run -d -p 80:1969 -p 443:1968 -p 8080:1936 \
           --name replyserver \
           bsctl/replyserver:latest -cert server-tls.crt -key server-tls.key

Deploy on OpenShift:

    oc apply -f oscp-deploy.yaml

Deploy on Kubernetes:

    kubectl apply -f k8s-deploy.yaml

That's all.
