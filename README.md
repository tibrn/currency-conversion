# Currency-conversion
Simple application that takes currencies rates from external third-party API and makes conversions.

# API
1. `/create` 
    - Create new project
    - returns API-Key for project

2. `/convert`
    - Convert value for symbol
    - authorized with project API-KEY in Authorization header
    - accepts json(e.g. {"symbol":"EUR/USD", value:5.65})
    - returns value converted

# CLI

## Build
`make cli`
##  convert

```
 convert [flags]
```

### Options

```
  -k, --api-key string   API Key used for authentication
  -h, --help             help for convert
  -s, --symbol string    Symbol for conversion (e.g. EUR/USD)
  -v, --value string     Value to be converted (e.g 1.7658)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.converter.yaml)
      --host string     Host for converter API (default "http://127.0.0.1:8081")
```


##  create

```
 create [flags]
```

### Options

```
  -h, --help   help for create
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.converter.yaml)
      --host string     Host for converter API (default "http://127.0.0.1:8081")
```
