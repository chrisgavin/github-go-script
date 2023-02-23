import * as core from "@actions/core";

export class Inputs {
	"token" = core.getInput("token", {required: true});
	"dir" = core.getInput("dir", {required: true});
}

export function get():Inputs {
	return new Inputs();
}
