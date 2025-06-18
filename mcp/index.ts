import { FastMCP } from "fastmcp";
import { z } from "zod"; // Or any validation library that supports Standard Schema
// @ts-ignore
import lint from "@qlik-trial/api-guidelines-lint";
import { readFile } from 'fs/promises';

const server = new FastMCP({
  name: "MCP POC server for R&D dev-teams",
  version: "1.0.0",
});


server.addTool({
  name: "api-lint",
  description: "lint a async & openapi API specifications",
  parameters: z.object({
    spec: z.string().describe("The async or openapi API specification to lint"),
  }),
  execute: async (args) => {
    let res = await lint(args.spec, args.spec.search('asyncapi') !== -1 ? 'SYSTEM_EVENTS_V2' : 'REST_V2');
    return JSON.stringify(res);
  },
});


const readFileContent = async (path: string) => {
  return readFile(path, "utf8");
};

const docs = {
  "system-event-guidelines.md": {
    name: "System Event Guidelines",
    description: "Async API guidelines for system events following AsyncAPI 3.0 specification",
  },
  "rest-guidelines.md": {
    name: "REST API Guidelines",
    description: "REST API guidelines following OpenAPI specification",
  },
  "cloud-event-migration-guide.md": {
    name: "Cloud Events Migration Guide",
    description: "A step-by-step guide to migrate from system events to the latest Cloud Events spec format & System Events guidelines",
  },
  "tenant_purged.md": {
    name: "Purge Tenant Recipe",
    description: "This page outlines what a tenant purge means from the perspective of your micro-service and the resources it persists.",
  },
  "ui-event-guidelines.md": {
    name: "UI Event Guidelines",
    description: "Async API guidelines for UI events",
  },
};

// create a loop for each document in the docs object and make is registered as a resource
Object.entries(docs).forEach(([fileName, { name, description }]) => {

  server.addResource({
    uri: `file://resources/${fileName}`,
    name,
    description,
    mimeType: "text/markdown",
    async load() {
      var text = await readFileContent(`./resources/${fileName}`);
      return {text};
    },
  });
});

server.start({
  transportType: "stdio",
});