# Getir-Case

## Fetch Data from MongoDB

### Request Payload

The request payload of the first endpoint includes a JSON object with the following fields:

- `"startDate"` (string): Start date in the format "YYYY-MM-DD".
- `"endDate"` (string): End date in the format "YYYY-MM-DD".
- `"minCount"` (number): Minimum count value for filtering.
- `"maxCount"` (number): Maximum count value for filtering.

### Sample Request Payload

```json
{
    "startDate": "2016-01-26",
    "endDate": "2018-02-02",
    "minCount": 2700,
    "maxCount": 3000
}
```

### Response Payload

The response payload of the first endpoint has the following structure:

- `"code"` (number): Status code. 0 indicates success.
- `"msg"` (string): Description of the code. Set to "success" for successful requests.
- `"records"` (array): Array of objects containing filtered items.

Each object in the "records" array has the following fields:

- `"key"` (string): Key value.
- `"createdAt"` (string): Date and time of creation in ISO 8601 format.
- `"totalCount"` (number): Sum of the "counts" array in the document.

### Sample Response Payload

```json
{
    "code": 0,
    "msg": "Success",
    "records": [
        {
            "key": "TAKwGc6Jr4i8Z487",
            "createdAt": "2017-01-28T01:22:14.398Z",
            "totalCount": 2800
        },
        {
            "key": "NAeQ8eX7e5TEg7oH",
            "createdAt": "2017-01-27T08:19:14.135Z",
            "totalCount": 2900
        }
    ]
}
```

## In-memory Endpoints

### POST Endpoint

#### Request Payload

The request payload of the POST endpoint includes a JSON object with the following fields:

- `"key"` (string): Key value.
- `"value"` (string): Value associated with the key.

#### Sample Request Payload
```json
{
    "key": "example-key",
    "value": "example-value"
}
```

##### Response Payload

The response payload of the POST endpoint returns an echo of the request or an error (if any).

### GET Endpoint

#### Request Payload

The request payload of the GET endpoint includes a query parameter:

- `"key"` (string): Key value.

#### Sample Request

GET http://localhost:8080/in-memory?key=active-tabs

#### Response Payload

The response payload of the GET endpoint returns a JSON object with the following fields:

- `"key"` (string): Key value.
- `"value"` (string): Value associated with the key.

#### Sample Response Payload
```json
{
    "key": "active-tabs",
    "value": "getir"
}
```
