import * as core from "@actions/core";
import * as exec from "@actions/exec";
import sourceMapSupport from "source-map-support";

import * as cache from "./cache";
import * as inputs from "./inputs";

async function main() {
	sourceMapSupport.install();

	await cache.restoreCache();
	const child_env = { ...process.env, ...{"GITHUB_TOKEN": inputs.get().token} }
	await exec.exec("go", ["run", "."], {env: child_env, cwd: inputs.get().dir});
	await cache.saveCache();
}

main().catch(error => core.setFailed(error.stack || error));
