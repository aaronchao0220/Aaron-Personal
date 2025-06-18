# Cloud Event Migration Guide

> <ins>Important Note</ins>: This step-by-step guide assumes your service uses Service Kits for event publishing and subscription.
> If not, some steps may not be applicable.

## Considerations before starting the migration

### Impact on consumers

- **Dependencies:** identify all the services that are directly dependent on the events you service is publishing
- **Downstream effects:** understand the journey of your events, the end-2-end flows, there are maybe sub-dependencies that need attention
- **Communication:** clearly communicate your intention and plans for the migration to the affected teams

> <ins>Info:</ins> If you want to know more about the CloudEvents projects, here are some reference links:
> - [Design Documents](https://github.com/qlik-trial/qac/blob/main/design_reviews/qac-subcouncil-qpaas/PS-21905-CloudEvents/PS-21905-CloudEvents.md#architecture--design)
> - [CloudEvents](https://cloudevents.io/)
> - [AsyncAPI](https://www.asyncapi.com/en)

## 1. Prepare for the migration

### 1.1 Update your documentation

Prepare your documentation by creating a new file containing your current specification in AsyncAPIv3 format. You will make the updates below in this new AsyncAPI yml file.

```bash
cp asyncapi.yaml asyncapi-v3.yaml
```

This will also allow you to compare the changes to your new AsyncAPI yml file side by side with the old version and make the migration process easier.

**Quick guide: migrate to AsyncAPIv3:**

#### 1.1.1 [Move metadata](https://www.asyncapi.com/docs/migration/migrating-to-v3#moved-metadata) & update the AsyncAPI version

> <ins>Note:</ins>
> In AsyncAPIv2 the `tags` property were placed outside the `info` object, now the tag has been moved under `info`, however,
> in the latest version of the API-guidelines, we have decided to remove the `tags` property as mandatory in the guidelines,
> as it was not used consistently across the services.

* Update the `asyncapi` version to `3.0.0` and the `info` object
* Remove the `tags` property
* Remove the `x-qlik-guidelines-version` property

```diff
- asyncapi: "2.0.0"
+ asyncapi: "3.0.0"
x-qlik-stability: stable
x-qlik-visibility: private
info:
    title: Groups
    version: 2.1.0
    description: Groups events in the platform.
-   x-qlik-guidelines-version: 3.0.0
- tags:
-   - name: groups
-     description: Groups is a resource that represents a Group in the platform.
```

#### 1.1.2 [Operation, channel, and message decoupling](https://www.asyncapi.com/docs/migration/migrating-to-v3#operation-channel-and-message-decoupling) and [Operation keywords](https://www.asyncapi.com/docs/migration/migrating-to-v3#operation-keywords)

Since the next step requires a focused attention to the details it has been split into smaller steps:

A. [Channel address and channel key](https://www.asyncapi.com/docs/migration/migrating-to-v3#channel-address-and-channel-key)

> <ins>Note:</ins> In AsyncAPI v3 the channel name has been moved to the `address`.
> Choose a descriptive channelId name (in this example `system-events.groups` translate well to `groupChannel`) it will later be used
> as a reusable key and move the channel name under `address`

* Update the `channel` name to `address` and add a key for the channel as shown in the example below

```diff
channels:
-  system-events.groups:
+  groupChannel:
+    address: 'system-events.groups'
```

B. [Messages instead of message](https://www.asyncapi.com/docs/migration/migrating-to-v3#messages-instead-of-message) & add keys for the messages

* Remove `publish`
* Update the `message` to `messages` and add keys for the message references, use a descriptive keyName

```diff
channels:
  groupChannel:
    address: 'system-events.groups'
-   publish:
-     message:
-       - $ref: "#/components/messages/groupDeleted"
-       - $ref: "#/components/messages/groupPurged"
+     messages:
+       groupDeleted:
+         $ref: "#/components/messages/groupDeleted"
+       tenantGroupPurged:
+         $ref: "#/components/messages/tenantGroupPurged"
```

C. Split operations and channels

> <ins>Note:</ins> In v2 `publish` and `subscribe`operations were causing confusion, specifying `publish`, implied that others could publish (according to the specs).
> In Qlik we were doing it the other way round, let's rectify (with AsyncAPIv3), by removing `subscribe` and add the more explicit `operations` (same level as `channels`)

```diff
  channels:
    groupsChannel:
      messages:
        groupDeleted:
          $ref: "#/components/messages/groupDeleted"
        tenantGroupPurged:
          $ref: "#/components/messages/tenantGroupPurged"
+ operations:
+   publishEvents:
+     action: send
+     channel:
+       $ref: "#/channels/groupsChannel"
+     messages:
+       - $ref: "#/channels/groupsChannel/messages/groupDeleted"
+       - $ref: "#/channels/groupsChannel/messages/tenantGroupPurged"
```

#### 1.1.3 Add the `cloudEventsAttributes` to the component schemas

* Copy the schema below and add it to your component's schemas

<details>
    <summary> <ins>cloudEventsAttributes schema definition in YAML and JSON format (click to expand)</ins> </summary>

```yaml
cloudEventsContextAttributes:
  required:
    - id
    - source
    - specversion
    - type
  description: CloudEvents Specification JSON Schema
  type: object
  properties:
    id:
      description: Identifies the event.
      type: string
      minLength: 1
      examples:
          - "A234-1234-1234"
    source:
      description: Identifies the context in which an event happened.
      type: string
      format: uri-reference
      minLength: 1
      examples:
          - "com.qlik/my-service"
    specversion:
      description: The version of the CloudEvents specification which the event uses.
      type: string
      minLength: 1
      examples:
          - "1.0"
    type:
      description: Describes the type of event related to the originating occurrence.
      type: string
      minLength: 1
      examples:
          - "com.qlik.v1.app.created"
    datacontenttype:
      description: Content type of the data value. Must adhere to RFC 2046 format.
      type: string
      minLength: 1
    time:
      description: Timestamp of when the occurrence happened. Must adhere to RFC 3339.
      type: string
      format: date-time
      minLength: 1
```

```json
{
  "cloudEventsContextAttributes": {
    "required": [
        "id",
        "source",
        "specversion",
        "type"
    ],
    "description": "CloudEvents Specification JSON Schema",
    "type": "object",
    "properties": {
      "id": {
          "description": "Identifies the event.",
          "type": "string",
          "minLength": 1,
          "examples": [
              "A234-1234-1234"
          ]
      },
      "source": {
          "description": "Identifies the context in which an event happened.",
          "type": "string",
          "format": "uri-reference",
          "minLength": 1,
          "examples": [
              "com.qlik/my-service"
          ]
      },
      "specversion": {
          "description": "The version of the CloudEvents specification which the event uses.",
          "type": "string",
          "minLength": 1,
          "examples": [
              "1.0"
          ]
      },
      "type": {
          "description": "Describes the type of event related to the originating occurrence.",
          "type": "string",
          "minLength": 1,
          "examples": [
              "com.qlik.v1.app.created"
          ]
      },
      "datacontenttype": {
          "description": "Content type of the data value. Must adhere to RFC 2046 format.",
          "type": "string",
          "minLength": 1
      },
      "time": {
          "description": "Timestamp of when the occurrence happened. Must adhere to RFC 3339.",
          "type": "string",
          "format": "date-time",
          "minLength": 1
      }
    }
  }
}
```

</details>

* Add the reference of the cloudEventsContextAttributes to your payloads

```diff
groupDeleted:
  ....
  name: com.qlik.v1.group.deleted
  ...
  payload:
    type: object
    allOf:
+     - $ref: "#/components/schemas/cloudEventsContextAttributes"
...
components:
  ...
  schemas:
+   cloudEventsContextAttributes:
+     type: object
+     properties: ...
```

#### 1.1.4 Now add the `cloudEventsQlikExtensionsAttributes` to the component schemas and add it at the root level of you events.

The new component schema (`cloudEventsQlikExtensionsAttributes`) is a replacement for the Qlik Extensions Attribute (`extensions: { tenantId:..., userId:..., etc...}`).

> <ins>Note:</ins> In CloudEvent v1 all the extensions context attributes are at the top-level and SHOULD follow same naming convention and
> use the same type system as the standard attributes:
>
>  - names should be lowercase alphanumeric characters
>  - names should be descriptive and concise, under 20 characters.
>  - allowed types are: Boolean, Integer, Number, String, URI, Timestamp, URI-Reference, Binary

* Add the new schema and copy all the properties from "extensions"
* Make sure all fields are in lowercase ex: `tenantId` -> `tenantid`
* Remove obsolete fields: The following fields are no longer used in the new version of the Events API guidelines, and should be removed from the specification:`description`, `group`, `consumptionContexts`, `eventTypeVersion`

> <ins>Important note:</ins> The removal of the `group` field from the latest System Events guidelines affects the [Purge Tenant recipe](/backend/recipes-and-patterns/tenant_purge), please make sure
> your service do not rely on it.

```diff
schemas:
+  cloudEventsQlikExtensionsAttributes:
+    type: object
+      required:
+        - tenantid
+    properties:
+      userid:
+        type: string
+        examples:
+          - "VZhiEfgW2bLd7HgR-jjzAh6VnicipweT"
+        description: Unique identifier for the user triggering the event.
+      tenantid:
+        type: string
+        examples:
+          - "VZhiEfgW2bLd7HgR-jjzAh6VnicipweT"
+        description: Unique identifier for the tenant related to the event.
+        ....
-      eventTypeVersion:
-        type: string
-      description:
-        type: string
-      group:
-        type: string
-      consumptionContexts:
-        type: array
```

##### Documented extensions:

<ins>`header` distributed tracing headers</ins>:

In CloudEvents v1 spec. uses the [Distributed Tracing](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/extensions/distributed-tracing.md)
documented extension, if used:

* Replace `headers` with the standardized extension

```diff
schemas:
  cloudEventsQlikExtensionsAttributes:
    type: object
      required:
        - tenantId
    properties:
      ...
-     header:
-       type: object
-       description: Carrier information for distributed trace using an HTTP header format
+     traceparent:
+       type: string
+       description: Contains a version, trace ID, span ID, and trace options.
+     tracestate:
+       type: string
+       description: A comma-delimited list of key-value pairs.
```


<details>
    <summary> <ins>header distributed tracing headers payload example (click to expand)</ins> </summary>

```diff
-"extensions": {
-  "header": {
-      "B3": ["{traceId}-{spanId}-1"],
-      "Traceparent": ["00-{traceId}-{spanId}-01"],
-      "X-B3-Sampled": ["1"],
-      "X-B3-Spanid": ["{spanId}"],
-      "X-B3-Traceid": ["{traceId}"]
-  },
-  ...
+"traceparent": "00-{traceId}-{spanId}-01"
+"tracestate": "b3={traceId}-{spanId}-{(optional)1}-{(optional){parentSpanId}}"
},
```

</details>

<ins>`actor` authentication context</ins>:
The actor attribute is used for embedding the impersonation information that triggered the action.
CloudEvents v1 uses the [Auth Context](https://github.com/cloudevents/spec/blob/main/cloudevents/extensions/authcontext.md) documented extension:

* Replace `actor` with the standardized extension

```diff
schemas:
  cloudEventsQlikExtensionsAttributes:
    type: object
      required:
        - tenantId
    properties:
      ...
-     actor
-       type: object
-       description": "actor includes the specifics of the impersonating entity.
+     authtype:
+       type: string
+       description: Representing the type of principal that triggered the occurrence.
+     authclaims:
+       type: string
+       description: A JSON string representing claims of the principal that triggered the event
```
<details>
    <summary> <ins>actor authentication context payload example (click to expand)</ins> </summary>

```diff
-"extensions": {
-    "actor": {
-        "iss": "qlik.api.internal/data-engineering-exporter",
-        "sub": "data-engineering-exporter",
-        "subType": "service"
-    }
-},
+"authtype": "service_account",
+"authclaims": "{\"iss\":\"qlik.api.internal/data-engineering-exporter\",\"sub\":\"data-engineering-exporter\",\"subType\":\"service\"}"
```
</details>

> Note: if your current extensionObjects fields are not documented in a reference schema - please make sure to convert
> each property to the new format and remove the old.

* Add the reference of the cloudEventsQlikExtensionsAttributes to your payloads

```diff
groupDeleted:
  ....
  name: com.qlik.v1.group.deleted
  ...
  payload:
    type: object
    allOf:
      - $ref: "#/components/schemas/cloudEventsAttributes"
+     - $ref: "#/components/schemas/cloudEventsQlikExtensionsAttributes"
...
```

#### 1.1.5 Remove the "old" schemas references, `cloudEventsv0` and `extensions` field

* Depending on how you called them previously ex: `extensionObject`, `cloudEvents` etc... you can now remove them from the specification

```diff
  groupDeleted:
    title: Group deleted
    name: com.qlik.v1.group.deleted
...
  payload:
    type: object
    allOf:
-     - $ref: "#/components/schemas/cloudEvents"
      - $ref: "#/components/schemas/cloudEventsAttributes"
      - $ref: "#/components/schemas/cloudEventsQlikExtensionsAttributes"
      - type: object
        properties:
          source:
            type: string
            default: "com.qlik/groups"
            description: The source of this event
          eventType:
            type: string
            default: com.qlik.v1.group.deleted
            description: Unique identifier for the event type.
-         extensions:
-           allOf:
-             - $ref: "#/components/schemas/extensionObject"

components:
...
schemas:
-    extensionObject:
    ...
-    cloudEvents:
```

* Rename `eventType` to `type`

```diff
  ...
  payload:
    type: object
    allOf:
      - $ref: "#/components/schemas/cloudEventsAttributes"
      - $ref: "#/components/schemas/cloudEventsQlikExtensionsAttributes"
      - type: object
        properties:
          source:
            type: string
            default: "com.qlik/groups"
            description: The source of this event
-     eventType:
+     type:
        type: string
        default: com.qlik.v1.group.deleted
        description: Unique identifier for the event type.
```

#### 1.1.6 Merge you PR

* validate your changes and ensure the descriptions matches the properties and that your API is compliant with the new [System Events API guidelines](https://internal.qlik.dev/general/api-strategy/guidelines/2025-event/)

<ins>**Validation tools & API Linter**</ins>

- [API Analyser (UI)](https://apiculturist.qlikdev.com/api-analyzer/) - select `SYSTEM EVENTS V2` as the Guideline ruleset
- [AsyncAPI studio](https://studio.asyncapi.com/)

* Once ready, and reviewed by your peers, merge the PR to default branch of your component


> Note: If you need any help or review or have any specific question about CloudEvents/AsyncAPI or api-governance
> use the [#api-design](https://qlikdev.slack.com/archives/CU68E4YAG) slack channel

### 1.1.7 API compliance

⚠ Make sure your API specification is fully compliant with the API guidelines.

⚠ Make also sure that all your fields are properly documented in the API specification and matches what the events your service publishes.

#### 1.1.7 Add your new API specification to the API-governance

    Use the [#api-governance](https://qlikdev.slack.com/archives/C0SHCF8BF) Slack channel, ask to have your API added by using the workflow Set up new API.
    Once the API is added to the API-governance, release a new version of the components.

#### 1.1.8 Disable the old API specification for API-governance

    When the old API specification is removed it needs to be disabled in the monitoring to state that it no longer requires a governance status. This is done
    by contacting the governance team in the [#api-governance](https://qlikdev.slack.com/archives/C0SHCF8BF) slack channel explaining that the API specification has been replaced.

## 2. Publish the new events

### 2.1 Implement the code changes in your service

CloudEvents support has been added to the Go Service Kit in [version v27](https://github.com/qlik-trial/go-service-kit/releases/tag/v27.0.0).
Along with the new CloudEvents methods changes, the legacy event functions were marked as deprecated. Please make sure you upgrade to v27 or later.

> Note: the legacy methods `PublishEvent` and `PublishAsyncEvent` now converts your event to a CloudEvent format & combines the legacy format into a
> temporary format we call SuperEvent used for backward/forward compatibility purposes

The new structure of the CloudEvent may affect code changes beyond the "publishing to solace" logic, and therefore is it recommended, if possible, to use
new objects/structures and helper function, for an easy rollback if needed.

You can now use the new methods exposed in the ServerKits for sending the new events - to ensure backward compatibility, the methods will automatically
duplicate the cloudEvents v0 fields for you

```go
deleteEvent := gskEvents.NewCloudEvent{
    Type:       "com.qlik.v1.resource.deleted",
    Source:     "com.qlik/my-service",
    Tenantid:   teanantid,
    Data: myDataPlayod
}
// ....
msgClient.PublishCloudEvent(ctx, "system-events.my-service", deleteEvents)

// other available methods
// msgClient.PublishAsyncCloudEvent(ctx, "system-events.my-service", deleteEvents, callBack)
// msgClient.PublishDirectCloudEvent(ctx, "system-events.my-service", deleteEvents)
```

> **Tenant purge**: If your service has implemented the [Purge Tenant recipe](/backend/recipes-and-patterns/tenant_purge) and uses the go-service-kit make sure you are using
> the CloudEvent method for publishing the Purge event
> ```diff
> --err = purge.RecordTenantPurgeEnded(ctx, mySubject, nil, myEventPublisher, baseInfo, resultInfo)
> ++err = purge.PublishTenantPurgeEnded(ctx, mySubject, nil, myEventPublisher, baseInfo, resultInfo)
> ```

Once release, properly communicate your intention to all your dependent services; If you are unsure about who is consuming your event, you can always refer to the access list
in messaging configurator see: https://github.com/qlik-trial/messaging-configurator/blob/main/manifests/chart/messaging-configurator/files/acls.json or for conveniance you can
use the [API analysis app - sheet Events Dependency graph](https://rd-sense.eu.qlikcloud.com/sense/app/b15a1d25-f12c-43ac-a700-f6d8b7680f94/sheet/QCBEz/state/analysis)

> Important Note: If you want to future-proof your API, it is recommended to run message validation against your async API specification, while some will, recommend
> ["runtime validation"](https://www.asyncapi.com/docs/guides/message-validation#runtime-validation), the most effective way would be to implement component test with message validation

## 3. Consume new events

The dev-teams can now implement the logic for consuming the new events. Below is an over simplified example

```diff
+ eventSub := messaging.NewCloudEventSubscription("system-events.my-resource", func(ctx context.Context, event *gskEvents.CloudEvent, err error) {
+     // handle new event
+ })

- eventSub := messaging.NewEventSubscription("system-events.my-resource" ...

err := messagingClient.Subscribe(eventSub)
```
