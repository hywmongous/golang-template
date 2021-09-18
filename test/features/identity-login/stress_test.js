import { runCaseWeighted } from '../index.js';
import { weightedCases } from './cases.js';
import {
    BELOW_NORMAL_LOAD,
    NORMAL_LOAD,
    BREAKING_POINT_LOAD,
    ABOVE_BREAKING_POINT_LOAD,
} from "./loads.js";

export const options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    minIterationDuration: '10s', // I think this can be used to simulate the session duration
    stages: [
        { duration: '2m', target: BELOW_NORMAL_LOAD }, // below normal load
        { duration: '5m', target: BELOW_NORMAL_LOAD },
        { duration: '2m', target: NORMAL_LOAD }, // normal load
        { duration: '5m', target: NORMAL_LOAD },
        { duration: '2m', target: BREAKING_POINT_LOAD }, // around the breaking point
        { duration: '5m', target: BREAKING_POINT_LOAD },
        { duration: '2m', target: ABOVE_BREAKING_POINT_LOAD }, // above the breaking point
        { duration: '5m', target: ABOVE_BREAKING_POINT_LOAD },
        { duration: '10m', target: 0 }, // scale down. Recovery stage.
    ],
};

export default () => { runCaseWeighted(weightedCases); }
