# URL Shortener

Simple URL shortening service written with Golang.
Helps users to get short URLs from long origin ones.
Based on RESTful API that allows:

* Create a new short URL
* Retrieve an original URL from a short URL
* Update an existing short URL
* Delete an existing short URL
* Get statistics on the short URL (e.g., number of times accessed)

## API Endpoints

### Create Short URL
`POST` /shorten
```json
{
  "url": "https://www.example.com/some/long/url"
}
```

Response Status: `201 Created`
```json
{
  "id": "1",
  "url": "https://www.example.com/some/long/url",
  "shortCode": "abc123",
  "createdAt": "2021-09-01T12:00:00Z",
  "updatedAt": "2021-09-01T12:00:00Z"
}
```
or a `400 Bad Request` status code with error messages in case of validation errors.
Short codes must be unique and should be generated randomly.

### Retrieve the original URL from a short URL 
`GET` /shorten/abc123

Response Status: `200 OK`
```json
{
  "id": "1",
  "url": "https://www.example.com/some/long/url",
  "shortCode": "abc123",
  "createdAt": "2021-09-01T12:00:00Z",
  "updatedAt": "2021-09-01T12:00:00Z"
}
```
or a `404 Not Found` status code if the short URL was not found.

# Origin roadmap link
https://roadmap.sh/projects/url-shortening-service