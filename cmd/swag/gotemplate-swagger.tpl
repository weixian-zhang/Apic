swagger: "2.0"
info:
  title: Welcome to Apic
  description: 
  version: 1.0.0
host: http://{{ .HostName}}:{{. Port}}
basePath: /
schemes:
  - http
paths:
  {{ with .ApiCmds }}
   {{ range . }}
   {{ .Path}}:
      get:
        summary: {{ .Description}}
        description: {{ .Description}}
        produces:
          - application/json
        responses:
          200:
           description: OK
      post:
        summary: {{ .Description}}
        requestBody:
          description: {{ .Description}}
          required: true
          content:
            application/json:
            schema:
              type: string
            text/plain:
              schema:
                type: string
        responses:
          '201':
            description: Created
      put:
        summary: {{ .Description}}
        requestBody:
          description: {{ .Description}}
          required: true
          content:
            application/json:
            schema:
              type: string
            text/plain:
              schema:
                type: string
        responses:
          '201':
            description: Created
      delete:
        summary: {{ .Description}}
        description: {{ .Description}}
        produces:
          - application/json
        responses:
          '201':
            description: Deleted
   {{end}}
  {{end}}
