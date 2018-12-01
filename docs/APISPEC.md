# Api specification

This document defines the endpoint responses and format

## /v1/categories GET

```json
[
  {
    "id": 1,
    "name": "drink"
  },
  {
    "id": 2,
    "name": "beer"
  }
]
```

## /v1/drinks GET

```json
[
  {
    "id": 1,
    "drink_name": "Rom&Rolf",
    "category": 1,
    "description": "För en mer traditionell fylla",
    "url": "http://qty.se"
  },
  {
    "id": 2,
    "drink_name": "GT",
    "category": 1,
    "description": "För britter, med krisp",
    "url": "http://qty.se"
  }
]
```

## /v1/drink/:id GET

```json
{
  "id": 1,
  "drink_name": "Rom&Rolf",
  "description": "För en mer traditionell fylla",
  "url": "http://qty.se",
  "Ingredients": [
    {
      "id": 1,
      "liquid_name": "Rolf",
      "liquid_id": 2,
      "volume": 33
    },
    {
      "id": 2,
      "liquid_name": "Rom",
      "liquid_id": 1,
      "volume": 3
    }
  ]
}
```

## /v1/liquid/:id GET

```json
{
  "id": 1,
  "name": "Rom"
}
```

## /v1/liquids GET

```json
[
  {
    "id": 1,
    "name": "Rom",
    "url": "http://qty.se"
  },
  {
    "id": 2,
    "name": "Rolf",
    "url": "http://qty.se"
  },
  {
    "id": 3,
    "name": "Gin",
    "url": "http://qty.se"
  },
  {
    "id": 4,
    "name": "Tonic",
    "url": "http://qty.se"
  }
]
```