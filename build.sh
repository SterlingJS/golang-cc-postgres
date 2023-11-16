#!/bin/bash
set -e

SUBSCRIPTION="085ce977-b6b1-44ad-ac4d-3d3ec9ec82c7"
AKS_CLUSTER="fmz-c-x-app-aks-01"
ACR=fmzexacr01.azurecr.io
NAMESPACE=keda-demo
CONSUMER_IMAGE=$ACR/keda-demo/consumer
PRODUCER_IMAGE=$ACR/keda-demo/producer
POSTGRES_IMAGE=$ACR/keda-demo/postgres
CONSUMER_TAG=1.0.12
PRODUCER_TAG=1.0.10
POSTGRES_TAG=1.0.3
PRODUCER_SPA_TAG=1.0.0
CONSUMER_CHART_NAME=keda-consumer
PRODUCER_CHART_NAME=keda-producer
POSTGRES_CHART_NAME=keda-postgres
CONSUMER_CHART_VERSION=0.1.36
PRODUCER_CHART_VERSION=0.1.20
POSTGRES_CHART_VERSION=0.1.16


az account set --subscription $SUBSCRIPTION
az acr login -n $ACR
kubectl config use-context $AKS_CLUSTER

docker build -t $CONSUMER_IMAGE:$CONSUMER_TAG ./consumer/
docker build -t $PRODUCER_IMAGE:$PRODUCER_TAG ./producer/
docker build -t $POSTGRES_IMAGE:$POSTGRES_TAG ./postgres/

docker push $CONSUMER_IMAGE:$CONSUMER_TAG
docker push $PRODUCER_IMAGE:$PRODUCER_TAG
docker push $POSTGRES_IMAGE:$POSTGRES_TAG


if [[ $(helm list -n $NAMESPACE | grep $CONSUMER_CHART_NAME) ]]; then
    echo "The $CONSUMER_CHART_NAME chart already exists"
    helm upgrade -f ./helm/consumer/values.yaml --version $CONSUMER_CHART_VERSION --set image.repository=$CONSUMER_IMAGE --set image.tag=$CONSUMER_TAG $CONSUMER_CHART_NAME ./helm/consumer -n $NAMESPACE
else
    echo "Installing $CONSUMER_CHART_NAME chart now..."
    helm install -f ./helm/consumer/values.yaml --version $CONSUMER_CHART_VERSION --set image.repository=$CONSUMER_IMAGE --set image.tag=$CONSUMER_TAG $CONSUMER_CHART_NAME ./helm/consumer -n $NAMESPACE
fi

if [[ $(helm list -n $NAMESPACE | grep $PRODUCER_CHART_NAME) ]]; then
    echo "The $PRODUCER_CHART_NAME chart already exists"
    helm upgrade -f ./helm/producer/values.yaml --version $PRODUCER_CHART_VERSION --set image.repository=$PRODUCER_IMAGE --set image.tag=$PRODUCER_TAG $PRODUCER_CHART_NAME ./helm/producer -n $NAMESPACE
else
    echo "Installing $PRODUCER_CHART_NAME chart now..."
    helm install -f ./helm/producer/values.yaml --version $PRODUCER_CHART_VERSION --set image.repository=$PRODUCER_IMAGE --set image.tag=$PRODUCER_TAG $PRODUCER_CHART_NAME ./helm/producer -n $NAMESPACE
fi

if [[ $(helm list -n $NAMESPACE | grep $POSTGRES_CHART_NAME) ]]; then
    echo "The $POSTGRES_CHART_NAME chart already exists"
    helm upgrade -f ./helm/postgres/values.yaml --version $POSTGRES_CHART_VERSION --set image.repository=$POSTGRES_IMAGE --set image.tag=$POSTGRES_TAG $POSTGRES_CHART_NAME ./helm/postgres -n $NAMESPACE
else
    echo "Installing $POSTGRES_CHART_NAME chart now..."
    helm install -f ./helm/postgres/values.yaml --version $POSTGRES_CHART_VERSION --set image.repository=$POSTGRES_IMAGE --set image.tag=$POSTGRES_TAG $POSTGRES_CHART_NAME ./helm/postgres -n $NAMESPACE
fi
