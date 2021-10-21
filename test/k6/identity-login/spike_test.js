import { runCaseWeighted } from "../index.js";
import { weightedCases } from "./cases.js";
import { BELOW_NORMAL_LOAD, SPIKE_LOAD } from "./loads.js";

export const options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    minIterationDuration: "10s", // I think this can be used to simulate the session duration
    stages: [
        { duration: "10s", target: BELOW_NORMAL_LOAD }, // below normal load
        { duration: "1m", target: BELOW_NORMAL_LOAD },
        { duration: "10s", target: SPIKE_LOAD }, // spike
        { duration: "3m", target: SPIKE_LOAD }, //
        { duration: "10s", target: BELOW_NORMAL_LOAD }, // scale down. Recovery stage.
        { duration: "3m", target: BELOW_NORMAL_LOAD },
        { duration: "10s", target: 0 },
    ],
};

export default () => {
    runCaseWeighted(weightedCases);
};
