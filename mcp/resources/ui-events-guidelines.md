# UI event guidelines

## Version

This is version **1.0.0** of the guidelines. It's a live document, and you
are very welcome to leave feedback or raise questions as either issues or
pull requests on [GitHub](https://github.com/qlik-trial/api-guidelines), or on [Slack](https://qlikdev.slack.com/archives/CU68E4YAG).

## Introduction

The Qlik UI event guidelines seek to guarantee quality, uniformity,
consistency, and usefulness of UI events provided both internally and publicly.

They follow the design principles provided by the
[Design Atlas framework](https://internal.qlik.dev/design/design-atlas/).

These guidelines are in many ways a contract, and it's assumed that any
given event in the platform adheres to it.

### Terminology

This section describes common terminology (and alternative keywords), and which
one that's preferred when you are referring to it.

Except for the requirement keywords "MUST," "MUST NOT," "REQUIRED," "SHALL,"
"SHALL NOT," "SHOULD," "SHOULD NOT," "RECOMMENDED," "MAY,"
and "OPTIONAL" used in these guidelines are to be interpreted as described in
[RFC 2119](https://www.ietf.org/rfc/rfc2119.txt), there are also some other common
keywords that **SHOULD** be followed:

| Keywords                                     | Preferred   | Description                                                                                          |
| -------------------------------------------- | ----------- | ---------------------------------------------------------------------------------------------------- |
| `Tools`, `Clients`, `UIs`, `CLIs`            | `UI`        | These guidelines apply to events, sent from client interfaces jobs-to-be-done.                       |
| `Content`, `Item`, `App`, `Chart`            | `Content`   | Content types are entities that the user can create and manage in our system.                        |
| `Property`, `Field`, `Key`                   | `Field`     | Used to describe a member of a data model or structure, usually in JSON formats with a defined type. |
| `Parameter`, `Argument`, `Variable`          | `Parameter` | Parameters are options you can pass to the tools and definitions to influence the outcome.           |
| `Resource`, `(Domain) entity`                | `Resource`  | Describes unique resources that clients may interact with in the platform.                           |
| `Actions`, `User Actions`, `Content Actions` | `Action`    | Common actions on content that users can do in our platform through tools                            |

## How to use the guidelines

### Updating these guidelines

Here are some guidelines to updating the guidelines, so there can be more
guidelines:

- See the [Google developer documentation style guide](https://developers.google.com/style/highlights)
  for things to think about when contributing to these guidelines.
- Be careful when changing headings since links from other places might break,
  VS Code with the "Markdown All In One" extension does a good job of updating
  the table of contents automatically.
- You may use the VS Code extension "markdownlint" to get linting applied in your
  IDE, alternatively run the `./tools/lint.sh` script locally.

### Application of the guidelines

These guidelines are applicable to any UI event exposed publicly or internally
by Qlik. Internal services tend to eventually be exposed publicly and
consistency is valuable to both external and internal system event consumers.

Internal events aren't required to follow the UI events guidelines except
for [**MUST:** Use standardized channel names](#must-use-standardized-channel-names)
which applies to all events in the platform.

### Existing UI events

It is highly recommended applying new rules from these guidelines on already existing
UI events defined and used in the platform tools.

All new UI events **MUST** follow the guidelines.

## Event payload format

### **MUST:** Use CloudEvents fields in a standardized way

Qlik uses CloudEvents to ensure all events are described in a common way. The
events **MUST** follow version `v1` of the specification.
[See the CloudEvents specification documentation](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md)

The JSON Event Format defines how each field in the event **MUST** be defined.
[See the JSON Event format documentation](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/formats/json-format.md)

While CloudEvents provide common fields and recommendations, Qlik need to
standardize on a way of using these fields. All fields in this list
**MUST** be included.

| Field             | Value                                         | Constraints | Description                                                                                                                 |
| ----------------- | --------------------------------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------- |
| `data`            | Object                                        | REQUIRED    | See [Naming and formatting](#naming-and-formatting)                                                                         |
| `specversion`     | `1.0`                                         | REQUIRED    | All system events **MUST** follow version `1.0`                                                                             |
| `datacontenttype` | `application/json`                            | OPTIONAL    | The contentType `application/json`                                                                                          |
| `id`              | `<uuid>`                                      | REQUIRED    | Unique identifier for the event instance                                                                                    |
| `time`            | `String`                                      | REQUIRED    | [ISO 8601](https://www.iso.org/iso/home/standards/iso8601.htm) format with Zulu time. (For example 2020-06-06T00:01:02.123Z) |
| `type`            | `com.qlik.v<version>.<eventContext>.<action>` | REQUIRED    | [**MUST:** Use Qlik pattern for event type](#must-use-qlik-pattern-for-event-type)                                          |
| `source`          | `com.qlik/<tool-name>`                        | REQUIRED    | **MUST** be defined as `com.qlik/<tool-name>`                                                                               |

### **MUST:** Use Qlik metadata attributes in a standardized way

The following Qlik metadata attributes are part of the CloudEvent extension,
and follow the same [naming convention](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md#attribute-naming-convention) and use the same type system as standard attributes

| Field       | Value    | Constraints | Description                                                                    |
| ----------- | -------- | ----------- | ------------------------------------------------------------------------------ |
| `tenantid`  | `<uuid>` | REQUIRED    | Tenant Identifier from where the Event has been triggered                      |
| `userid`    | `<uuid>` | REQUIRED    | Internal User Identifier of the Event that represents who triggered the action |
| `sessionid` | `<uuid>` | OPTIONAL    | Session identifier of the Event                                                |

Example event payload:

```json
{
   "specversion":"1.0",
   "source":"com.qlik/analytics-app-client",
   "type":"com.qlik.v1.analytics.analytics-app-client.sheet-view.opened",
   "id":"<eventId>",
   "time":"2024-01-30T07:06:22Z",
   "tenantid": "<tenantId>",
   "userid": "<userId>",
   "sessionid": "<sessionId>",
   "data":{...}
}
```

### **MUST:** Use Qlik pattern for event type

To have predictable and uniform `type` the field **MUST** be
defined as `com.qlik.v<version>.<context>.<action>`.
The `type` **MUST** be unique in the system.

- `version`: **MUST** be the major part of eventTypeVersion and **MUST** follow
  the tool versioning.
  For example if the resource or tool the event relates to is major version `v1`
  then the event **MUST** also have major version `v1`.
- `context`: **MUST** be the API tool name, and sub-tool if the
  event applies to one.
  The `context` **MUST** use `kebab-case`.
- `action`: **MUST** be the action performed on the `context` and follow
  [**MUST:** Use standardized action names](#must-use-standardized-action-names)

### **MUST:** Use standardized action names

All actions **MUST** be named as a past-tense verb as they have already happened.
Examples: `opened`, `shared`, `enabled`

To keep the eventing uniform within the platform the following standards
**MUST** be followed:

| Type of event  | Action       | Description                                    |
| -------------- | ------------ | ---------------------------------------------- |
| User Action    | `browsed`    | When a user views and compare a list of items  |
| User Action    | `opened`     | When a user opens a UI                         |
| User Action    | `shared`     | When a shares something with others            |
| ...            | `...`        |                                                |
| Content Action | `duplicated` | Content gets duplicated                        |
| Content Action | `imported`   | Content get imported                           |
| Content Action | `disabled`   | Resource disabled to prevent it from executing |
| ...            | `...`        |                                                |

Examples:

- `com.qlik.v1.analytics.analytics-app-client.sheet-view.opened`
- `...`

### **MUST:** Use standardized channel names

The channel name **MUST** be consistently named to avoid confusion and conflicts.
The channel name **MUST** use the form `ui-events.<ui>` and follow the structure
defined in the [Design Atlas](https://internal.qlik.dev/design/design-atlas/ui/).

If the events need to be separated further the `uis` can be
suffixed to the channel name forming `system-events.<ui>.<ui>`.

Channel names **MUST** be kebab-cased.

Examples:

- `ui-events.analytics`
- `ui-events.automations`
- `ui-events.data-integration`
- `ui-events.management-console`
- `ui-events.analytics.insight-chat`

## Naming and formatting

### **MUST:** Use `camelCase` for parameters and fields

There is no official specification on how multi-word query parameters should be
separated, but the industry praxis is to either use `snake_case` or `camelCase`.
To avoid inconsistencies in naming between query parameters and response body
fields, `camelCase` **MUST** be used.

Examples:

```json
{
  "myAwesomeProp": ""
}
```

### **MUST:** Use standardized field names

This section describes a set of standard field names that **MUST** be used when similar concepts are needed.

| Name           | Type       | Description                                                                                               | Examples that **MUST NOT** to be used                  |
| -------------- | ---------- | --------------------------------------------------------------------------------------------------------- | ------------------------------------------------------ |
| `createdAt`    | `DateTime` | DateTime when the resource was created                                                                    | `created`, `createdDate`                               |
| `createdBy`    | `String`   | Id of the user who created the resource, **MUST** reference a userId                                      | `creator`, `createdByUser`                             |
| `description`  | `String`   | Longer description on a resource                                                                          | `title`, `comment`, `metadata`                         |
| `name`         | `String`   | Short name assigned to the resource                                                                       | `title`, `description`                                 |
| `ownerId`      | `String`   | Owner id of the resource, **MUST** reference a userId                                                     | `owner`, `creator`, `userId`                           |
| `<resource>Id` | `String`   | Field that holds the id of another resource, e.g. tenantId, spaceId. The resource **MUST** be in singular | `tid`, `tenantsId`, `userIdentifier`                   |
| `updatedAt`    | `DateTime` | DateTime when the resource was created                                                                    | `modified`, `updatedDate`, `modifiedAt`, `lastUpdated` |
| `updatedBy`    | `String`   | Id of the user who last updated the resource, **MUST** reference a userId                                 | `update`, `updatedByUser`, `updatedByUserId`           |
| `...`          | `String`   | ....                                                                                                      | `...`                                                  |

## API Documentation

### **MUST:** Follow general AsyncAPI guidelines

- Each UI event generated by the tool **MUST** be documented.
- The documentation **MUST** be done in [AsyncAPI 2.0.0](https://v2.asyncapi.com/docs/reference)
- The documentation **MUST** include the payload data, with the JSON structure described
  using [JSON-Schema](https://json-schema.org).

### **MUST:** Include general AsyncAPI documentation fields

- For the general event description, the documentation **MUST** include:
  - `info`, which in turn **MUST** include `title`, `version` and `description`.
    The `description` field is used by tooling, and by providing proper
    metadata the user confidence in the provided events can be increased.
    Example of good description: The audit events allow you to gain insight of
    what goes on in your tenant. Example of bad description: Events for the management
    service.

Example Async API document snippet:

```json
{
  "asyncapi": "2.0.0",
  "info": {
    "title": "analytics app",
    "description": "API specification for the events that are sent the analytics apps (sense client).",
    "version": "1.0.0",
    "x-qlik-guidelines-version": "1.0.0"
  },
  "channels": {}
}
```

### **MUST:** Use Qlik AsyncAPI extension for `Visibility`

The specification **MUST** include `x-qlik-visibility` to specify the
intended audience for the events. Allowed values:

- `public` - Event is intended for public use
- `private` - Event is intended for Qlik internal use only.

Example Async API document snippet:

```json
{
  ...,
  "messages": {
      "sheetOpened": {
        "name": "com.qlik.v1.analytics.analytics-app-client.sheet-view.opened",
        "title": "User open a sheet",
        "x-qlik-visibility": "public",
        ...
}
```

### **MUST:** Use Qlik AsyncAPI extension for `Stability`

The specification **MUST** include `x-qlik-stability` to specify what
to expect from the events in terms of changes,
it also denotes the deprecation period. Allowed values:

- `stable`
  - Highly reliable
  - Breaking changes are extremely unlikely
- `experimental`
  - Unreliable
  - Experimental features aren't subject to the versioning policy
  - These events can be changed at any time without violating API governance

### **MUST:** Use Qlik AsyncAPI extension for `Deprecation`

The specification **MUST** include `x-qlik-deprecated` if an endpoint is
deprecated. Allowed values are `true` and `false`.
The deprecation policy is described in more detail in [apiculturist](https://apiculturist.qlikdev.com/governance/policy).

Default value is `false`. The deprecation period is further described in the
[Stability extension](#must-use-qlik-asyncapi-extension-for-stability).

### **MUST:** Use Qlik AsyncAPI extension for `Guidelines Version`

The specification **MUST** include `x-qlik-guidelines-version` in the `info` object
of the AsyncAPI specification. Allowed values are released semantic versions of
these API guidelines.

### **MUST:** Use Qlik AsyncAPI extension for Personally Identifiable Information fields

Any field that has been identified as a field containing Personally
Identifiable Information (PII) **MUST** be annotated with `x-qlik-pii: true`.
The purpose of `x-qlik-pii: true` is to be able to highlight PII fields in API
documentation and to aid in PII related tasks.

Any uncertainties regarding fields being PII or not should be brought up with
the Software Security Office (SSO).

Example:

```json
{ ...
  "User": {
    "type": "object",
      "properties": {
          "name": {
            "type": "string",
            "x-qlik-pii": "true"
          },
}
```

### **MUST:** Use Qlik AsyncAPI extension for Customer Data fields

Any field that has been identified as a field containing Customer Assets Data or
Customer Content Data **MUST** be annotated with `x-qlik-customer-data: true`.

Example:

```json
{ ...
  "Sheet": {
    "type": "object",
      "properties": {
          "title": {
            "type": "string",
            "x-qlik-customer-data": "true"
          },
}
```

### **MUST:** Include common fields for events

- For each event, the documentation **MUST** include:
  - `name`: name of the event that **MUST** match `type` for consistent documentation
  - `title`: human-friendly name of the event
  - `description`: a shorter description of the event
  - `tags`: list of the tags. Tags are used to group the paths and operations.

### **MUST:** Use the API governance machinery

To secure the contract, live up to the [API governance policy](https://apiculturist.qlikdev.com/governance/policy)
and highlight API changes within and outside R&D each API **MUST** use the
API governance machinery. Consumers of an API expect it to be backwards compatible.

Details on setup can be found in [Apiculturist](https://apiculturist.qlikdev.com/api-tooling-setup).

### **SHOULD:** Use the style guide for the metadata

To ensure that the documentation for the APIs are consistent, it's encouraged
to follow the [Google developer documentation style guide](https://developers.google.com/style/).
Some highlights can be found [here](https://developers.google.com/style/highlights).
Description fields are supporting Markdown syntax as described by [CommonMark](https://spec.commonmark.org/0.27/).
