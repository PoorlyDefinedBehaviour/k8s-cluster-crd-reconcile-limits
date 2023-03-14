# About

Tested using minikube

## Running the controller

Generate the CRD files
```
make manifests
```

Install the CRD in the cluster
```
make install
```

Run the controller
```
make run
```

## Create resources

Create resources
```
go run ./scripts/cmd/loadtest.go
```

