# PostgreIntern

SETUP:
path - address of a directory
;commands - set of commands for execution
  - change path and commands in config.yaml

Run:
- docker-compose up -d 
- cd cmd/
- go run main.go
- change some file in your path

Check DB: 
- docker exec -it {container_id} psql -U admin -d events

  
 
