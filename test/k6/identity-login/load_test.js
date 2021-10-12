import { runCaseWeighted } from "../index.js";
import { weightedCases } from "./cases.js";
import { NORMAL_LOAD } from "./loads.js";

export const options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    minIterationDuration: "10s", // I think this can be used to simulate the session duration
    stages: [
        { duration: "5m", target: NORMAL_LOAD },
        { duration: "10m", target: NORMAL_LOAD },
        { duration: "5m", target: 0 },
    ],
    thresholds: {
        http_req_duration: ["p(99)<1000"],
    },
};

export default () => {
    runCaseWeighted(weightedCases);
};
