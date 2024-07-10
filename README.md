# Canapads-api

This is part of Canapads.ca.  This is a website being created to provide access to rental properties in Canada. 

![canapads_polygon](https://github.com/jebo87/canapads-api/assets/7273200/5c2dd29e-4cbc-4f26-91c5-602615f753cb)

this microservice is the backend API used to provide all data related to listings. It uses gRPC to talk back to makako-gateway but also gets information from Elastic and Postgres

Frontend is being built using ReactJS and the backend is composed of microservices written in Go and NodeJS

We also use ElasticSearch for searching, Postgres as a Database, Mongo DB for user storage, ORY Hydra as authorization server 

This is a work in progress in early stages and I am currently working in what you see in the diagram below, however more micro services will be added in the future.

![canapads architecture](https://github.com/jebo87/canapads-api/assets/7273200/a46be66d-4b87-46e1-a32b-9f59387ba96a)

