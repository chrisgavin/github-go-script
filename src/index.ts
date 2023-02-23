import * as core from "@actions/core";
import * as exec from "@actions/exec";
import {promises as fsPromises} from "fs";
import * as path from "path";
import sourceMapSupport from "source-map-support";

import * as inputs from "./inputs";

async function main() {
	sourceMapSupport.install();

	const child_env = { ...process.env, ...{"GITHUB_TOKEN": inputs.get().token} }
	await exec.exec("go", ["run", "."], {env: child_env, cwd: inputs.get().dir});
}

main().catch(error => core.setFailed(error.stack || error));
