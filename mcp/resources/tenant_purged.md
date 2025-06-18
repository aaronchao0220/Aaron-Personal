# Purge Tenant recipe


**Note: this recipe describes a process that is in-implementation. Final details are expected to be
in place some time at the beginning of 2022 Q1.**

## Introduction

This page outlines what a tenant purge means from the perspective of your micro-service and the resources it persists.

The purge of a tenant is implemented using an event driven approach, and is initiated by the `tenants` service.
All micro-services persisting user data **MUST** participate in the tenant purge.
Your micro-service will listen to the tenant purged event message and subsequently respond to it according to this document.

The purges will be audited on a platform level (GDPR compliance) using a Splunk Dashboard and a Qlik Sense App.
This is why services are required to follow the enforced logging and event formats, outlined in sections below.

## Micro-service purge - Overview

The generalized micro-service purge procedure and expectations are captured in the overview below.

<Image src={overview} alt="tenant purge overview"/>

| Requirement Number | Requirement                                    | Description                                                                                                                                                                        |
| ------------------ | ---------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| REQ_1              | [Subscribe to tenant purged event](#subscribe-to-tenant-purged-event) | The micro-service should listen and trigger purges based on `tenant purged` events.  |
| REQ_2              | [Feature flag](#feature-flag)                                     | Put your actual purge calls behind a toggle.                                                                                                            |
| REQ_3              | [Bulk purging process](#bulk-purging-process)                     | Resources should be purged in one go, not one by one.                                                                                                              |
| REQ_4              | [Purge logging](#purge-logging)                                  | Your micro-service will provide purge related logs according to the defined format.                                                |
| REQ_5              | [Publish purge event](#publish-purge-event)                            | Send a  `com.qlik.v1.<eventContext>.purged` event when purging a resource.                                           |

## Subscribe to tenant purged event

The `tenant.purged` event initiates the purge of all resources, related to a specific tenantId.
Notes on payload:

<TabList>
  <Tab name="CloudEvent">
  In the CloudEvents spec format - Qlik extension context attributes are in the root level and the meta property has been moved under data
  - `tenantid` denotes the tenantID to be purged
  - `_meta.purgeId` enables the tracking of an entire tenant purge across participating micro-services.
```json
{
  "specversion": "1.0",
  "type": "com.qlik.v1.tenant.purged",
  "source": "com.qlik/tenants",
  "id": "BPo6E1j1--J4LgauBybscWZY0fPph52b",
  "time": "2022-09-28T20:17:24.445Z",
  "datacontenttype": "application/json",
  "tenantid": "vUcpSAGgp2RdJLztTqT1883CYzAsnidX",
  "data": {
    "_meta": {
      "purgeId": "6334abd4bf68c98d9e8e4ab6"
    },
    "id": "vUcpSAGgp2RdJLztTqT1883CYzAsnidX",
    "name": "kf248psa73at1er",
    "hostnames": ["kf248psa73at1er.eu.qlikdev.com"]
  }
}
```
  </Tab>
  <Tab name="legacy event">
  - `extensions.tenantId` denotes the tenantID to be purged
  - `extensions.meta.purgeId` enables the tracking of an entire tenant purge across participating micro-services.
```json
{
  "cloudEventsVersion": "0.1",
  "eventType": "com.qlik.v1.tenant.purged",
  "eventTypeVersion": "1.0.0",
  "source": "com.qlik/tenants",
  "eventID": "BPo6E1j1--J4LgauBybscWZY0fPph52b",
  "eventTime": "2022-09-28T20:17:24.445Z",
  "contenttype": "application/json",
  "extensions": {
    "tenantId": "vUcpSAGgp2RdJLztTqT1883CYzAsnidX",
    "meta": {
       "purgeId": "6334abd4bf68c98d9e8e4ab6"
    }
  },
  "data": {
    "id": "vUcpSAGgp2RdJLztTqT1883CYzAsnidX",
    "name": "kf248psa73at1er",
    "hostnames": ["kf248psa73at1er.eu.qlikdev.com"]
  }
}
```
  </Tab>
</TabList>

## Bulk purging process

The bulk purging function should execute upon receiving the `com.qlik.v1.tenant.purged` event.
All user data for the given tenant id should be purged directly.

The purge process **MUST NOT** make authenticated requests to other services except to feature-flag service
(user JWT, S2S jwt).

The tenant has already been moved to `purged` state.

You will execute this purge function upon receiving the `com.qlik.v1.tenant.purged` event.

### Purging options

Your micro-service receives the `com.qlik.v1.tenant.purged` event and initiates the Bulk Delete.
Option A or B below are up to the service-owners to freely decide on.

**Option A**:

For each service resource type: (Resource-A|Resource-B|Resource-C)

1. Service logs that the purge started.
2. Service purges all entries of the resource type.
3. Service logs that the purge ended, and publishes a system event for the purged resource type.

**Option B**:

For the entire service:

1. Service logs that the purge started.
2. Service purges all (sub)resources.
3. Service logs that the purge ended, and publishes one system event in total.

## Purge logging

The actual tenant purge routine of a micro-service, **MUST** be accompanied by logging to Splunk.
See [Purging options](#purging-options) for more details on options on when/how to perform logging.

The initiation and the subsequent result of a purge are the two distinct log entries expected when performing a tenant purge.
Logging helpers are provided by the Go and Node Service kits, facilitating the adoption of the standard log formatting.

Proof of deletion will be obtained by performing Splunk logs reviews:
- [Splunk Dashboard: Purge tenant list](https://qlik.splunkcloud.com/en-US/app/search/qcs_purge_tenant_list)

Your micro-service **MUST** perform tenant purge `start` and `ended` logging.
Log entries will surface when formatted according to the contract below, in the dashboard below.

- [Splunk Dashboard: Detailed info](https://qlik.splunkcloud.com/en-GB/app/search/qcs_stage_tenant_purge_detailed_info)

### Formats

Purge start log example:

```json

{
    ...
    "action": "purge",
    "level": "info",
    "message": "purge started",
    "purgeId": "1b83bb2c-d55b-4f0d-824f-760752348812",
    "resourceType": "audits",
    "tenantId": "ca8f46b2-5bd0-42b4-b296-286d8dc7be3c",
    "timestamp": "2022-10-28T07:02:59.617486624Z",
    "traceId": "7b2f58801c3b76c7c13f809251fd81ea",
}
```

Purge ended log example:

```json

{
    ...
    "action": "purge",
    "errorMessage": "",
    "level": "info",
    "message": "purge ended",
    "purgeId": "1b83bb2c-d55b-4f0d-824f-760752348812",
    "purgedCount": 0,
    "resourceType": "audits",
    "success": true,
    "tenantId": "ca8f46b2-5bd0-42b4-b296-286d8dc7be3c",
    "timestamp": "2022-10-28T07:02:59.621913412Z",
    "traceId": "7b2f58801c3b76c7c13f809251fd81ea",
}
```

## Publish purge event

Your service **MUST NOT** trigger a snowstorm of resource events. For example if you use a MongoDB.deleteMany call
to get rid of all 2000 apps for that tenant, do not send 2000 app deleted events. Typically send one purged event
to account for the entire purge of all instances of that resource type for that tenant.

Purge event publishing helpers are provided by the Go and Node Service kits, to enforce the standard formatting across micro-services.

Proof of deletion will be obtained by reviewing the event published by your service for a particular tenant purge.
The purge events can be tracked using the Qlik Sense App [Tenant Purge Tracker](https://qlikinternal.us.qlikcloud.com/sense/app/8773b2c2-b072-4c0d-acde-7c1a98b3c9de/sheet/475ac4fd-e10f-444e-94bb-6687e70bda11/state/analysis).

### API-guidelines compliance

The `com.qlik.v<version>.<eventContext>.purged` events **MUST** follow the standardized purge format, see:
[System events guidelines](/general/api-strategy/guidelines/2025-event) and [must-use-standardized-purged-action-format](/guidelines/event/event-payload-format/must-use-standardized-purged-action-format).


### Async-API event template

Template specification to facilitate the `*.purged` system-events onboarding process.

<TabList>
  <Tab name="CloudEvent">

```yaml
asyncapi: 3.0.0
x-qlik-stability: stable
x-qlik-visibility: private
info:
  title: Example com.qlik.v1.myEventContext.purged example format
  description: string
  version: 1.0.0
channels:
  myServiceChannel:
    address: system-events.service
    messages:
      tenantResourcePurged:
        $ref: "#/components/messages/tenantResourcePurged"
operations:
  publishEvents:
    action: send
    channel:
      $ref: "#/channels/myServiceChannel"
    messages:
      - $ref: "#/channels/myServiceChannel/messages/tenantResourcePurged"
components:
  schemas:
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
    cloudEventsQlikExtensionsAttributes:
      type: object
      required:
        - tenantid
      properties:
        tenantid:
          type: string
          example: VZhiEfgW2bLd7HgR-jjzAh6VnicipweT
          description: Unique identifier for the tenant related to the event.
    tenantPurgedResult:
      type: object
      description: 'Result of the tenant purged operation, for the specific resourceType.'
      required:
        - purgedCount
        - purgeId
        - resourceType
        - success
      properties:
        errorMessage:
          type: string
          description: Detailed message about a failed purge operation.
          example: Failed to connect to database.
        purgedCount:
          type: number
          description: Number of resources that was successfully purged.
          example: 1000
        purgeId:
          type: string
          description: Unique identifier of the tenant purge request.
          example: 00000000-0000-0000-0000-000000000000
        resourceType:
          type: string
          description: Type of resource that was purged.
          example: audits
        success:
          type: boolean
          description: Status of the purge operation.
          example: true
  messages:
    tenantResourcePurged:
      title: TenantResource purged
      name: com.qlik.v1.myEventContext.purged
      description: This event will be published when the myEventContext of a tenant has been purged.
      payload:
        type: object
        allOf:
          - $ref: '#/components/schemas/cloudEventsContextAttributes'
          - $ref: '#/components/schemas/cloudEventsQlikExtensionsAttributes'
          - type: object
            required:
              - data
            properties:
              type:
                type: string
                enum:
                  - com.qlik.v1.myEventContext.purged
                description: Unique identifier for the event type.
              data:
                type: object
                description: Purge specific data related to the event.
                allOf:
                  - $ref: '#/components/schemas/tenantPurgedResult'

```

  </Tab>  
  <Tab name="legacy format">

```yaml
asyncapi: 2.0.0
x-qlik-stability: stable
x-qlik-visibility: private
info:
  title: Example com.qlik.v1.myEventContext.purged example format
  description: string
  version: 0.0.1
  x-qlik-guidelines-version: 3.2.1
tags:
  - name: string
    description: string
channels:
  system-events.<service>:
    publish:
      message:
        oneOf:
          - $ref: '#/components/messages/TenantResourcePurged'
components:
  schemas:
    cloudEvents:
      type: object
      required:
        - cloudEventsVersion
        - contentType
        - eventID
        - eventTime
        - eventTypeVersion
        - source
      properties:
        cloudEventsVersion:
          type: string
          enum:
            - '0.1'
          description: The event follows CloudEvents version 0.1.
          example: '0.1'
        contentType:
          type: string
          default: application/json
          description: Event payload will be application/json.
          example: 'application/json'
        eventID:
          type: string
          description: Unique identifier for the event.
          example: '00000000-0000-0000-0000-000000000000'
        eventTime:
          type: string
          description: Timestamp of when the event happened.
          format: date-time
          example: '2018-10-30T07:06:22Z'
        eventTypeVersion:
          type: string
          description: Indicates the version of the event.
          example: '1.0.0'
        source:
          type: string
          default: com.qlik/<service>
          description: Denotes the source from which events originates.
          example: 'com.qlik/audit'
    extensions:
      type: object
      description: Additional metadata and custom fields
      required:
        - group
        - tenantId
      properties:
        group:
          type: string
          description: Identifier by which the purge events are aggregated in Data Engineering.
          enum:
            - purged
        tenantId:
          type: string
          description: Unique identifier for the tenant related to the event.
          example: '00000000-0000-0000-0000-000000000000'
    tenantPurgedResult:
      type: object
      description: Result of the tenant purged operation, for the specific resourceType.
      required:
        - purgedCount
        - purgeId
        - resourceType
        - success
      properties:
        errorMessage:
          type: string
          description: Detailed message about a failed purge operation.
          example: 'Failed to connect to database.'
        purgedCount:
          type: number
          description: Number of resources that was successfully purged.
          example: 1000
        purgeId:
          type: string
          description: Unique identifier of the tenant purge request.
          example: '00000000-0000-0000-0000-000000000000'
        resourceType:
          type: string
          description: Type of resource that was purged.
          example: 'audits'
        success:
          type: boolean
          description: Status of the purge operation.
          example: true
  messages:
    TenantResourcePurged:
      title: TenantResource purged
      name: com.qlik.v1.myEventContext.purged
      description: This event will be published when the myEventContext of a tenant has been purged.
      tags:
        - name: string
      payload:
        type: object
        allOf:
          - $ref: '#/components/schemas/cloudEvents'
          - type: object
            required:
              - data
              - eventType
              - extensions
            properties:
              eventType:
                type: string
                enum:
                  - com.qlik.v1.myEventContext.purged
                description: Unique identifier for the event type.
              extensions:
                type: object
                description: Additional metadata and custom fields
                allOf:
                  - $ref: '#/components/schemas/extensions'
              data:
                type: object
                description: Purge specific data related to the event.
                allOf:
                  - $ref: '#/components/schemas/tenantPurgedResult'

```
  </Tab>
</TabList>

## Repeated tenant purged events

Most of the time your service will only receive one such message for a given tenantId.

Occasionally you may receive multiple messages for the same tenantId (a few minutes, an hour, a day or a month later).

The re-purging will help services failing the first purge attempt by giving them a second or third chance.

If your service was successful the first time, the second time will also claim success and a `purgedCount` of 0.
Note that the `purgedCount` of 0 will be interpreted as proof that no further work was needed by your service.

## Scalability

Tenant service is taking some precautions to send the `com.qlik.v1.tenant.purged` events at different times. However
it is still possible your service would receive a large number of purge events for large tenants with lots of resources
you need to delete. We don't want this to affect the SLO's of your service or the MongoDB / S3 QCS infrastructure.
So what to do to reduce that risk?

In short: take your time to consume `com.qlik.v1.tenant.purged` events. They are intentionally published to a dedicated
queue `system-events.tenants.purged`. There are no event types in that queue that require time-critical handling.
Let's clarify how to configure your solace subscription.

To handle the `com.qlik.v1.tenant.purged` event, subscribe to the group queue `system-events.tenants.purged`.
Only process one purge event at a time per pod, to reduce the number of MongoDB deletions running in parallel. Your pod
should not pre-fetch 100 purge events from the queue and process them in parallel, do not flood MongoDB with a large
number of delete requests in parallel.

Do you need to implement a leader election pattern so only 1 pod processes purge events in your service? No, for most
services this should be overkill. It's okay your replicas both process one event at a time. Just don't aggressively purge
many tenants concurrently inside a pod. Use the Splunk dashboard to see if your purge times are a source of concern when
we purge large tenants.

### Instructions for Solace

You want to limit the number of in flight Solace purge messages to 1. This is set within
[https://github.com/qlik-trial/messaging-configurator/blob/main/manifests/chart/messaging-configurator/files/queue-overrides.json]()

 "maxDeliveredUnackedMsgsPerFlow": 1

on a per service/queue basis. There are some real world examples already in this JSON file.

We may use a durable queue if we need to keep messages for multiple hours.

### Instructions for S3

According to [Best practices design patterns: optimizing Amazon S3 performance](https://docs.aws.amazon.com/AmazonS3/latest/userguide/optimizing-performance.html)
S3 can handle at least 3,500 requests/second per [prefix](https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-prefixes.html).

Limiting ourself to handling one event at a time should be enough to
make sure we are no where near the S3 scalability limitations. Given
an S3 latency of 10 ms a single thread can make only 100 requests per
second. Not even 3% of what S3 can "at least" do.

The API both support [Deleting a single object](https://docs.aws.amazon.com/AmazonS3/latest/userguide/delete-objects.html)
and [Deleting multiple objects](https://docs.aws.amazon.com/AmazonS3/latest/userguide/delete-multiple-objects.html)

Use what is convenient. [DELETE and CANCEL requests are free](https://aws.amazon.com/s3/pricing/)
while [LIST requests for any storage class are charged at the same rate as S3 Standard PUT, COPY, and POST requests.](https://aws.amazon.com/s3/pricing/)

### Instructions for MongoDB

Before deleting any resources, ensure you know your MongoDb client's write timeout setting.
If there is a risk of exceeding the write timeout, consider splitting the deletion process into smaller batches.

If you know the maximum number of records you can create for a tenant, use the `deleteMany` function to delete all the resources within that limit.
For instance, in the identity-provider service, the maximum number of IDPs allowed to be created is 100.
Therefore, using the `deleteMany` function to delete all the resources that fall within that limit is safe.

If you cannot predict the total number of records that need to delete, use the "bulk select and delete with an interval" method.
This involves selecting a batch of resources that need to delete and then deleting them with an interval between each batch.
This ensures the system is not overwhelmed with too many delete requests.

If you are using a shared MongoDB cluster, it is recommended to use a setting of 2000 as batch size with a 5-second wait time. However, if your microservice uses a standalone MongoDB cluster, you can adjust the setting according to your specific needs.

- [Golang Example](https://github.com/qlik-trial/identities/blob/main/internal/persistence/batchFindDelete.go)
- [NodeJS Example](https://github.com/qlik-trial/node-service-kit/blob/b20565eee60cd5de95d6e6c30c9a195b72deaa57/lib/dal/BaseDao.js#L372)

If "bulk select and delete with an interval" fits your workflow, it is recommended to set the read preference to secondary for deleting resources.
This helps ensure the load on the primary is manageable, which could impact the performance of the entire cluster. See more details [here](https://www.mongodb.com/docs/manual/core/read-preference)

During performance and scale tests, ensure the delete operation does not pressure the entire database cluster.
This can be achieved by monitoring the performance of the microservice during deletion and adjusting the interval between delete batches as necessary to avoid overloading the system.

## Test plan

Team must test in SDE, in stage (can take advantage of our purge test we will keep running for ever in future)
Verify that their Stage MongoDB is cleared for expected tenant id's, ensure you don't purge the entire region or wrong tenant.

Team should verify if purging a large number of their resources could take a long time.
Does it affect the cluster / their SLO budget / lock down the MongoDB etc.
They can work with P&S team to get help setting this up in Stage.
