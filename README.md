# makako-api

This is part of Canapads.ca.  This is a website being created to provide access to rental properties in Canada. 

this microservice is the backend API used to provide all data related to listings. It uses gRPC to talk back to makako-gateway but also gets information from Elastic and Postgres

Frontend is being built using ReactJS and the backend is composed of microservices written in Go and NodeJS

We also use ElasticSearch for searching, Postgres as a Database, Mongo DB for user storage, ORY Hydra as authorization server 

This is a work in progress in early stages and I am currently working in what you see in the diagram below, however more micro services will be added in the future.

![network](/uploads/6f6abd016e87fba283bc32402e202bd7/network.png)
