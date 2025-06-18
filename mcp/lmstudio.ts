import { LMStudioClient, tool } from "@lmstudio/sdk";
import { z } from "zod";
import lint from "@qlik-trial/api-guidelines-lint";
import { readFile } from "fs/promises";

const client = new LMStudioClient();

const lintTool = tool({
    name: "lint",
    description: "lint a async api specification",
    parameters: {
        spec: z.string().describe("The async api specification to lint"),
    },
    implementation: async ({ spec }) => {
        let res = await lint(spec, 'SYSTEM_EVENTS_V2');
        return JSON.stringify(res);
    },
});

const readFileContent = async (path: string) => {
    return readFile(path, "utf8");
  };

const x = async () => {
    const model = await client.llm.model("mistralai/magistral-small");
    
    let spec = await readFileContent("./specs/access-controls.yaml");
    let result = await model.act("lint the following async api specification: \n" + spec, [lintTool], {
        onMessage: (message) => console.info(message.toString()),
    });

}

x().then(() => {
    console.info("done");
});