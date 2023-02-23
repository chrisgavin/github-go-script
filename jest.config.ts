import type {Config} from "@jest/types";

const config: Config.InitialOptions = {
	preset: "ts-jest",
	testEnvironment: "node",
	testMatch: [
		"<rootDir>/tests/**/*.test.ts"
	],
	clearMocks: true,
};
export default config;
