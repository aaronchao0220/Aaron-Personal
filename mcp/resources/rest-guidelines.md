# REST Guidelines

## **MUST:** Use standardized response body

Navigation links are a way to simplify traversation of an API resource dataset for clients. These links are vital when cursor-based pagination is used since it is impossible for the client to calculate the cursors.

These links **MUST** be available on a `links` field on the root object in `GET` operations that returns a list of resources.

The actual list of resources **MUST** be available on a `data` field on the root object.

### HTTP Response Example

```json
{
 "data": [
   ...
 ],
 "links": {
   "prev": { "href": "..." },
   "next": { "href": "..." }
 }
}
```

## **MUST:** Use standardized format for links

If present, `links` field **MUST** be an object. The members of the object must be "link objects," references to "link objects," or `null` if a link doesn't exist.

A link object (a member of the `links` field) **MUST** have `href` field which value is a string pointing to the link's target.

These guidelines are following the [JSON:API](https://jsonapi.org/format/#document-links) specification in regard to links object.

### HTTP Response Example:

```json
{
  "links": {
    "prev": {
      "href": "https://example.org/api/v1/apps?page=<previous-cursor>"
    },
    "next": {
      "href": "https://example.org/api/v1/apps?page=<next-cursor>"
    }
  }
}
```

## **SHOULD:** Allow sorting on metadata fields

An API **SHOULD** have support for sorting on metadata fields. This makes it easier for a client to fetch sorted resources based on for example which entities that were updated most recently.

When supported, the terminology `created`, `updated`, `deleted` **MUST** be used in the sorting names.

### Example

```bash
# sort on last updated in descending order:
https://example.org/api/v1/apps?sort=-updatedAt
```

## **MUST:** Allow sorting by ascending or descending order

The default order **MUST** be ascending. This can either be declared by simply including a field name as sort condition, or by explicitly prefix it using a plus (+) character.

To use descending order, the field in the sort condition **MUST** be prefixed with a minus (-) character.

### API Specification

```yaml
paths:
  /apps:
    get:
      summary: List apps
      description: Returns a list of apps.
      parameters:
        ...
        - name: sort
          in: query
          required: false
          description: The field to sort by, with +/- prefix indicating sort order
          schema:
            type: string
            enum:
              - name
              - +name
              - -name
              - size
              - +size
              - -size
            default: "+name"
```

### HTTP Request Example

```bash
# sort on the name field using ascending order:
https://example.org/api/v1/apps?sort=name
# or:
https://example.org/api/v1/apps?sort=+name

# sort on the name field using descending order:
https://example.org/api/v1/apps?sort=-name

# sort on the name field in descending order,
# and secondary by the size in ascending order:
https://example.org/api/v1/apps?sort=-name,+size
```

## **MUST:** Support sorting

An API that supports listing resources **MUST** at the very least support rudimentary sorting using the `sort` query parameter.

## **MUST:** Use standardized query parameters

The query parameter `sort` **MUST** be used when supporting sorting. It **MAY** support multiple sort conditions by seperating fields using commas.

### API Specification

```yaml
paths:
  /users:
    get:
      operationId: getUsers
      summary: List users
      description: Returns a list of users using cursor-based pagination.
      parameters:
        ...
        - name: sort
          in: query
          required: false
          description: The field to sort by, with +/- prefix indicating sort order
          schema:
            type: string
            enum:
              - name
              - +name
              - -name
              - id
              - +id
              - -id
              - tenantId
              - +tenantId
              - -tenantId
              - status
              - +status
              - -status
              - subject
              - +subject
              - -subject
              - createdAt
              - +createdAt
              - -createdAt
            default: "+name"
```

### HTTP Request Example

```bash
# sort on the name field:
https://example.org/api/v1/users?sort=name

# sort on the createdAt field, then on status:
https://example.org/api/v1/users?sort=createdAt,status
```

## **SHOULD:** Use `page` query parameter for navigation

The reserved query parameter `page` **MAY** be used to identify the targeted page inside the resource dataset. Since pagination relies on links in the API response, the exact query parameters used to solve pagination **MAY** vary between APIs.

If it the pagination query parameters are missing from a request, you **MUST** return the first page of the dataset.

### HTTP Request Example

```http
> GET /api/v1/apps?page=a435jfb2 HTTP/1.1
> host: example.org
> authorization: bearer xyz
> accept: application/json

< HTTP/1.1 200
< content-type: application/json
< content-length: 123

{
  "data": ["list", "of", "resources"],
  "links": {
    "prev": { "href": "https://example.org/api/v1/apps?page=<previous-cursor>" },
    "next": { "href": "https://example.org/api/v1/apps?page=<next-cursor>" },
  }
}
```

## **MUST:** Use `limit` query parameter for page result size

The `limit` query parameter allows the client to define the page result size, and **MUST** be supported when pagination is supported. If the query parameter is omitted, the default value **MUST** be used.

The max page size value **MUST** be used when the `limit` query parameter value is larger than what you support in your API.

### HTTP Request Example

```http
> GET /api/v1/apps?page=a435jfb2&limit=3 HTTP/1.1
> host: example.org
> authorization: bearer xyz
> accept: application/json

< HTTP/1.1 200
< content-type: application/json
< content-length: 123

{
  "data": ["list", "of", "resources"],
  "links": {
    "prev": { "href": "https://example.org/api/v1/apps?page=<previous-cursor>" },
    "next": { "href": "https://example.org/api/v1/apps?page=<next-cursor>" },
  }
}
```

## **MUST:** Support pagination

An API that supports listing resources **MUST** at the very least support pagination.

There are many reasons to support pagination, some of them are:

- Avoiding introduction of breaking changes later on in the API lifecycle
- Reduce processing needed to serve a request
- Reduce bandwidth used
- Improve client performance

A resource-listing method **MUST** support pagination regardless of the expected number of resources accessible to a user.

These guidelines are following the [JSON:API specification](https://jsonapi.org/format/#fetching-pagination) in regards to pagination.

### HTTP Request Example

```http
> GET /api/v1/apps?page=a435jfb2&limit=3 HTTP/1.1
> host: example.org
> authorization: bearer xyz
> accept: application/json

< HTTP/1.1 200
< content-type: application/json
< content-length: 123

{
  "data": ["list", "of", "resources"],
  "links": {
    "prev": { "href": "https://example.org/api/v1/apps?page=<previous-cursor>" },
    "next": { "href": "https://example.org/api/v1/apps?page=<next-cursor>" },
  }
}
```

## **MUST:** Have a default and max page result size

There **MUST** be default and max page sizes clearly defined in your API.

The exact values of the default and max page sizes varies between APIs due to for example performance or scalability concerns, intended use, or values before guidelines were in place.

If there are no concerns, the **RECOMMENDED** default value is **20** items and the **RECOMMENDED** max value is **100** items.

### API Specification

```yaml
openapi: '3.0.0'
paths:
  /databases:
    get:
      parameters:
        - name: limit
          in: query
          description: The preferred number of entries returned
          required: false
          schema:
            type: integer
            format: int32
            default: 20
            minimum: 1
            maximum: 100
```

## **MUST:** Use a standardized error format

When an error occur during a request, the API MUST include an errors field on the root-level object in the response. It MAY be the only member of the root-level object being returned.

The errors field MUST be an array, containing at least one error.

Each error returned **MUST** be JSON, and contain information that's concrete or
unique enough to indicate what went wrong. The following fields **MUST** be
defined in the error object.

| Field    | Description                                                                                                                                              |
| -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `code`   | The unique code for the error, it should have a unique uppercase, alphabetic prefix, followed by a hyphen, and a numeric error code. Example: `FOO-123`. |
| `title`  | A summary in english explaining what went wrong. Example: `You don't have permission to view this resource`.                                             |
| `detail` | Optional. **MAY** be used to provide more concrete details. Example: `This resource requires app read privileges`.                                       |
| `meta`   | Optional. **MAY** be used to provide more arbitrary data to the client. Note that this data needs to be described in your API specification.             |

When an error occurs during a request, the API **MAY** include a traceId field on the root-level object in the response. This is an OpenTracing or OpenTelemetry Trace ID, and can be used to simplify debugging of the failed request.

### Response body

```json
{
  "errors": [
    {
      "code": "FOO-123",
      "title": "You don't have permission to view this resource.",
      "detail": "This resource requires app read privileges.",
      "meta": {
        "locale": "en-GB",
        "type": "mismatch"
      }
    }
  ],
  "traceId": "4bf92f3577b34da6a3ce929d0e0e4736"
}
```

### More information

More information about the error format in the [JSON:API](https://jsonapi.org/format/#document-top-level) specification.

## **MUST NOT:** Use version in path

With the switch to date and header based versioning the version is handled there, for that reason no version should be included in the path since it has no semantic value anymore.
This will enable us to evolve the platform API without changing the paths.

### Do

```text
/analytics/reload-tasks...
```

### Don't

```text
/v1/analytics/reload-tasks/...
```

## **MUST:** Pluralize resource names

Usually, a collection of resource instances is provided. The special case of a
resource singleton is a collection with cardinality 1.

Resource names **MUST** be pluralized.

**Exception**: if the term doesn't have a suitable plural form (uncountable nouns),
such as "evidence" and "feedback," the singular form **SHOULD** be used.

## **MUST:** Use `kebab-case` for path segments

To ensure consistency all path segments part of a URL **MUST** use kebab-case
(lowercase separate words with hyphens).

### Example

```html
/reload-tasks/{reload-task-id}
```

An exception is template segments in specification documents. In these cases,
`kebab-case`, `camelCase`, or `snake_case` **MAY** be used. For example
`{reload_task_id}`, `{reloadTaskId}` are all acceptable in this case since they
aren't part of the URL an end-user would send requests to.

## **MUST:** Use noun for resources names

The API **MUST** be resource based, so the
only place where actions **MUST** appear is in the HTTP methods.

Using nouns (or "things") is a more natural way of exposing APIs than for
example verbs (or "actions").

> Note: All new root resource names must be approved by PM and QAC during the feature design review.

## **MUST:** Use approved resource in the URL

An API **MUST** use approved resources in the url.
It **MUST** reside all the way to the left in the URL to give it the highest scope.
All new APIs **MUST** use a namespaced resource name to evolve the platform in a consistent way.

If a new resource/namespace needs to be introduced it **MUST** be approved by QAC and PM and added [approved resource namespaces](https://github.com/qlik-trial/qac/blob/main/decision_logs/api2.0/allowed-resources-namespaces/rest2_restAllowedResourcesNamespaces.json).
The list of currently approved resources can be found [approved resources](https://github.com/qlik-trial/qac/blob/main/decision_logs/api2.0/allowed-resources-namespaces/rest1_restAllowedResources.json).

_NOTE: During a transition period both the flat platform API and the namespaced API will be available in parallel._

### Do

```text
/analytics/reload-tasks...
```

### Don't

```text
/v1/reload-tasks/...
```

## **MUST:** Use standardized naming when possible

As listed under the [general terminology section](/general/api-strategy/guidelines/2025-rest/#terminology), `settings` is
the preferred term for this capability.

In addition, where applicable (paths, fields, query parameters, etc.), the following
terms **MUST** be used.

Naming settings fields `SHOULD` be done in a way where they reflect the current state from a reader’s perspective, not the one modifying it. For example, enabled instead of enable.
A user reading would verify if something is enabled, a user modifying is seeking to enable.

In addition, this should be suffixed to the capability’s name - for example `automationsEnabled`, not `enabledAutomation`.

| Keywords                                   | Preferred | Description                                      |
| ------------------------------------------ | --------- | ------------------------------------------------ |
| `Enabled`, `On`, `Off`, `Toggle`, `Active` | `Enabled` | Whether a capability is `enabled` or `disabled`. |

## **MUST:** Expose settings using `/settings` sub-resource

When your API supports settings, it **MUST** be configured under a sub-resource
called `/settings`.

```http
/api/v1/foos/settings
```

## **SHOULD:** Only support `GET`, `PATCH`, and `PUT` operations

_Usually_, settings exist regardless of whether a user creates them or not.

You **MUST** support `GET` (fetch), `PATCH` / `PUT` (update) operations.

You **SHOULD NOT** support `POST` (create) or `DELETE` (delete) operations in
the vast majority of cases.
