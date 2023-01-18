# Intro
In this tutorial it is demonstrated how to use Go client for docker engine API to copy files inside a running container. 

## create and start container

```
docker run --name awesome_mgm -d -it --rm alpine
```

## tar a file

```
tar -cf file.tar.gz script.sh

```
