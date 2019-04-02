# ElasticSearch Feeder
This repository contains the needed code to feed data from the blockchain to ElasticSearch. The project is written in [Go](https://golang.org/).
To run this component correctly, you should be familiar with the [Data marketplace](https://github.com/lgsvl/data-marketplace) components because there is a particular dependency between the components.

# Build prerequisites
  * Install [golang](https://golang.org/).
  * Install [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git).
  * Configure go. GOPATH environment variable must be set correctly before starting the build process.

### Download and build source code

```bash
mkdir -p $HOME/workspace
export GOPATH=$HOME/workspace
mkdir -p $GOPATH/src/github.com/lgsvl
cd $GOPATH/src/github.com/lgsvl
git clone git@github.com:lgsvl/data-marketplace-esearchFeeder.git
cd data-marketplace-esearchFeeder
./scripts/build
```

### Docker Build

You can use the dockerfile to build a docker image:
```
docker build -t esearchfeeder .
docker run data-stream-delivery
```

### Kubernetes deployment

The [deployment](./deployment) folder contains the deployment and persistent volume claim manifests to deploy this component.
We assume that there is an ElasticSearch cluster already deployed in Kubernetes. We also assume that [Data marketplace Chaincode REST](https://github.com/lgsvl/data-marketplace-chaincode-rest) is already deployed.


# Testing data-stream-delivery

Run the tests:
```bash
./scripts/run_glide_up
./scripts/run_units.sh
```