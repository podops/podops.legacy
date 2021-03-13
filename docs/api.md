# Rest API

Based on simple REST principles, the PodOps API endpoints return JSON metadata about shows and episodes, directly from the PodOps backend.

## Requests

The PodOps API is based on REST principles. Data resources are accessed via standard HTTPS requests in UTF-8 format to an API endpoint. Where possible, the API uses appropriate HTTP verbs for each action:

| METHOD | ACTION                                          |
|--------|-------------------------------------------------|
| GET | Retrieves resources                                |
| POST | Creates resources                                 |
| PUT  | Changes and/or replaces resources or collections  |
| DELETE | Deletes resources                               |


## Responses
The API returns all response data as a JSON object. See the **Object Model** for a description of all the retrievable objects.


## URIs and IDs
In requests to the API and responses from it, you will frequently encounter the following parameters:

| PARAMETER	    | DESCRIPTION                                | EXAMPLE       |
|---------------|--------------------------------------------|---------------|
| Name          | 'name' identifies a resource. The following rules apply:<br>'name' must contain only lowercase letters, numbers, dashes (-), underscores (_).<br>'name' must contain 8-44 characters.<br>Spaces and dots (.) are not allowed.| simple_podcast<br>my-first-podcast|
| GUID          | A unique ID identifying a resource. | 91804b93b56a  |
| ParentGUID    | Same as GUID, references the parent resource for a given resource, e.g. the show an episode belongs to. | c7b5414a9c02  |
| Kind | Identifies the type of resource. | show<br>episode<br>asset |


## Timestamps
Timestamps are 64-bit integers, representing the number of seconds elapsed since January 1, 1970 UTC. In other words, UNIX time.


## Response Status Codes
The API uses the following response status codes, as defined [here](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes):

|STATUS CODE | DESCRIPTION                                   |
|------------|-----------------------------------------------|
| 200        | OK - The request has succeeded. The client can read the result of the request in the body and the headers of the response. |
| 201	     | Created - The request has been fulfilled and resulted in a new resource being created. |
| 202	     | Accepted - The request has been accepted for processing, but the processing has not been completed. |
| 204	     | No Content - The request has succeeded but returns no message body. |
| 206        | Partial Content - The server is delivering only part of the resource (byte serving) due to a range header sent by the client. |
| 307        | Temporary Redirect |
| 400	     | Bad Request - The request could not be understood by the server due to malformed syntax. |
| 401	     | Unauthorized - The request requires user authentication or, if the request included authorization credentials, authorization has been refused for those credentials. |
| 403	     | Forbidden - The server understood the request, but is refusing to fulfill it. |
| 404	     | Not Found - The requested resource could not be found. This error can be due to a temporary or permanent condition. |
| 500	     | Internal Server Error. Nothing you can do. |


## Authentication
All requests to API require authentication. This is achieved by sending a valid JWT access token in the request header (Bearer token). 
**details tbd**


## API Endpoint Reference
The API provides a set of endpoints, each with its own unique path. The endpoints enable external applications to access and manipulate PodOps data.

The base address of the PodOps API is https://api.podops.dev.

| METHOD | ENDPOINT                       | USAGE                              | 
|--------|--------------------------------|------------------------------------|
| GET    | /                              | Returns the version                |
| GET    | /_a/token                      | Verifies a token                   | 
| GET    | /a/v1/productions              | Get several productions            |
| POST   | /a/v1/production               | Create a new production            |
| GET    | /a/v1/resource/:prod/:kind/:id | Get a resource                     |
| GET    | /a/v1/resource/:prod/:kind     | Retrieve several resources         | 
| POST   | /a/v1/resource/:prod/:kind/:id | Create a resource                  |
| PUT    | /a/v1/resource/:prod/:kind/:id | Update a resource                  |
| DELETE | /a/v1/resource/:prod/:kind/:id | Delete a resource                  | 
| POST   | /a/v1/build                    | Start the build                    | 
| POST   | /a/v1/upload/:prod             | Upload asset                       | 

The following list of endpoints is either for internal use only or serves an administrative purpose:

| METHOD | ENDPOINT                       | USAGE                        | 
|--------|--------------------------------|------------------------------|
| POST   | /_a/token                      | Create an authorization      |
| POST   | /_t/import                     | Background asset import task |
