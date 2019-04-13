# Fusio server

Fusio is an IoT-platform written in Go


## Usage


Compile cmd/fusio: 
```
go build . 
```

Copy config file example-config.yaml-> config.yaml and fill it where needed
as a minimum set the database correctly

Load demo data into database:
```
./fusio -c config.yaml -d
```

Finally, run server: 
```
./fusio -c config.yaml
```

For help come visit matrix room [#fusio:icy.name](https://matrix.to/#/#fusio:icy.name)

