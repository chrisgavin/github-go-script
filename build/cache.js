"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.saveCache = exports.restoreCache = void 0;
const cache = __importStar(require("@actions/cache"));
const fs = __importStar(require("fs"));
const core = __importStar(require("@actions/core"));
const exec = __importStar(require("@actions/exec"));
async function cachePaths(existingOnly) {
    let path = "";
    await exec.exec("go", ["env", "GOMODCACHE"], {
        silent: true,
        listeners: {
            stdout: (data) => {
                path += data.toString();
            }
        },
    });
    const paths = [path.trim()];
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
async function restoreCache() {
    const paths = await cachePaths(false);
    if (paths.length === 0) {
        return;
    }
    const key = await cacheKey();
    core.info(`Restoring cache with key ${key}.`);
    await cache.restoreCache(paths, key);
}
exports.restoreCache = restoreCache;
async function saveCache() {
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
exports.saveCache = saveCache;
//# sourceMappingURL=cache.js.map