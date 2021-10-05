import { runCaseRoundRobin } from "../index.js";
import { weightedCases } from "./cases.js";

export const options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    minIterationDuration: '1s',
    vus: 1,
    duration: '1m',
}

export default () => { runCaseRoundRobin(weightedCases); }
