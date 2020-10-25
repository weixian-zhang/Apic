#### What is APIC

Apic is cli tool that mocks Restful APIs with custom Path, QueryString, Response, Headers, Cookies, and serves Swagger docs for mocked APIs.  

#### Installation

* Windows  
[Download apic.exe](shorturl.at/asGN5)

* Linux

* Docker  
  [docker pull wxzd/apic:alpine-0.9](https://hub.docker.com/repository/registry-1.docker.io/wxzd/apic/tags?page=1)


* Kubernetes  
  Deloyment yaml coming soon
  
 #### Usage
  
 * Basic usage: `apic rest`
 * Custom Api and Swagger port: `apic rest -p 8071 --swaggerport 8072`
 * Custom response: `apic rest -p 8071 --swaggerport 8072 -r {\"userName\":\"weixian\"}`  
 
 * Run in Docker: `docker run --rm -p 8071:8071 -p 8072:8072 wxzd/apic:alpine-0.9 /app/apic rest -p 8071 --swaggerport 8072`
