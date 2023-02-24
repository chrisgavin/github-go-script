import * as cache from "@actions/cache";
import * as fs from "fs";
import * as core from "@actions/core";
import * as exec from "@actions/exec";

async function cachePaths(existingOnly:boolean):Promise<string[]> {
	let path = "";
	exec.exec("go", ["env", "GOMODCACHE"]).then((stdout) => {
		path += stdout;
	});

	const paths:string[] = [path.trim()];
	if (!existingOnly) {
		return paths;
	}

	return paths.filter(async (path) => {
		const exists = await fs.promises.access(path).then(() => true).catch(() => false);
		if (!exists) {
			core.warning(`Path ${path} does not exist so it will not be cached.`);
		}
		return exists;
	});
}

async function cacheKey() {
	return "github-go-script-go-module-cache";
}

export async function restoreCache():Promise<void> {
	const paths = await cachePaths(false);
	if (paths.length === 0) {
		return;
	}
	const key = await cacheKey();
	core.info(`Restoring cache with key ${key}.`);
	await cache.restoreCache(paths, key);
}

export async function saveCache():Promise<void> {
	const paths = await cachePaths(true);
	if (paths.length === 0) {
		return;
	}
	const key = await cacheKey();
	core.info(`Saving cache with key ${key}.`);
	try {
		await cache.saveCache(paths, key);
	}
	catch (e) {
		if (e instanceof cache.ReserveCacheError) {
			core.warning(e);
		}
		else {
			throw e;
		}
	}
}
