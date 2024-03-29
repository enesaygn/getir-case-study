# Getir-Case

## Usage

## Running on Heroku

To access the deployed application on Heroku, you can use the provided API endpoints directly in Postman. Below are the steps to do so:

1. **Open Postman**: If you haven't already, download and install Postman from [here](https://www.postman.com/downloads/).

2. **Import Collection**: To access the API endpoints, you can directly use the provided Postman collection file included in the project. Follow these steps:

   a. Inside the project directory, you'll find a file named `Getir-Cases.postman_collection.json`.

   b. Open Postman and navigate to the "Collections" tab in the sidebar.

   c. Click on the "Import" button and select the `Getir-Cases.postman_collection.json` file from the project directory.

   d. Once imported, you will see the collection named "Getir Cases" in your Postman collections.

### Running Locally

To run the application on your local machine, follow these steps (applicable for Mac):

1. Navigate to the root directory of the project in your terminal.

2. Set the port:
    ```
    export PORT=8080
    ```

3. Run the application:
    ```
    go run main.go
    ```

4. Access the application by visiting `http://localhost:8080` in your browser or using an HTTP client.

### Local Debugging with VS Code

If you want to run and debug your application locally using VS Code, you can configure your launch settings as follows:

1. Open the `.vscode/launch.json` file in your VS Code project.

2. Add a new configuration object inside the `configurations` array, like this:

```json
{
    "name": "Launch Package",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceFolder}",
    "env": {
        "PORT": "8080"
    }
}
```

3. Now you can debugging.

### Testing

To run tests for the application, follow these steps:

1. Navigate to the root directory of the project in your terminal.

2. Run tests for all packages:
    ```
    go test -v ./...
    ```

By following these steps, you can successfully run and test the application on Heroku or your local machine.


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
