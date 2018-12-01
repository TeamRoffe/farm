# Api specification
This document defines the endpoint responses and format

## /metrics

Prometheus metrics endpoint https://github.com/prometheus/client_golang

## Api endpoints

The different endpoints explained for the different components

### Client

#### /healthz
Health endpoint, is the client running

#### /v1/status
Shows the F.A.R.M client status

{
    "status:":"idle|pouring",
    "pumps": [
        {
            "no": 0
            "status": "idle|pouring",
            "liquid": "jægermeister",
            "rate": "100",
        },
        {
            "no": 1
            "status": "idle|pouring",
            "liquid": "jægermeister",
            "rate": "100",
        }
        ...
    ]
}

#### /v1/ingredients
{
    "liquid": [
        {
            "name": "jægermesiter",
            "type": "soda|spiritus|other
            "description": "Kittlar dödsskönt i kistan",
            "pump": 0,
            "rate": 1400, # Pour rate, milliseconds per centiliter
        },
        {
            "name": "cola",
            "type": "soda|spiritus|other
            "description": "Kittlar dödsskönt i kistan",
            "pump": 1,
            "rate": 100, # Pour rate, milliseconds per centiliter
        }
        ...
    ]

}


### Server
