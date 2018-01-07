# Set-up
Start MQTT in edge-docker

# Start reader
```
go run rw_msg.go -M R
```

# Start writer
```
go run rw_msg.go -M W -N 100000
```
