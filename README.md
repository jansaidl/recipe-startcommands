# Zerops x multiple startCommands

```yaml
zerops:
  - setup: app
    build:
      os: alpine
      base: go@latest
      buildCommands:
        - go build -o main main.go
      deployFiles: main
    deploy:
      readinessCheck:
        httpGet:
          host: localhost
          port: 8080
          path: /
          scheme: http
    run:
      base: alpine@latest
      ports:
        - httpSupport: true
          port: 8080
        - httpSupport: true
          port: 8081
        - httpSupport: true
          port: 8082
        - httpSupport: true
          port: 8083

      # this is backward compatibility: unit name zerops@zerops
      start: ./main --port 8080
      initCommands:
        - "echo 'INIT COMMANDS: service zerops is listening on port 8080'"
      startCommands:

        # this is unit zerops@foo
        - command: ./main --port 8081
          name: foo

        # this is unit zerops@bar
        - command: ./main --port 8082
          name: bar
          workingDir: /var/www
          initCommands:
            - "echo 'INIT COMMANDS: service bar is listening on port 8082'"

        # this is unit zerops@3
        - command: ./main --port 8083
          initCommands:
            - "echo 'INIT COMMANDS: service 3 is listening on port 8083'"

```