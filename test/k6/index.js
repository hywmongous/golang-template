import { Counter } from "k6/metrics";
import { randomString } from "https://jslib.k6.io/k6-utils/1.1.0/index.js";

export const PORT = "80";
export const PROTOCOL = "http";
export const SUBDOMAIN = ""; // eg. 'www.'
export const ROOT_DOMAIN = "localhost";
export const TLD = ""; // eg. '.com'
export const BASE_URL = `${PROTOCOL}://${SUBDOMAIN}${ROOT_DOMAIN}${TLD}:${PORT}`;
export function buildUrl(basic, path) {
    if (path === undefined) {
        path = "";
    }
    return encodeURI(
        `${PROTOCOL}://${basic}@${SUBDOMAIN}${ROOT_DOMAIN}${TLD}:${PORT}${path}`
    );
}

export const CSRF_HEADER_KEY = "Csrf";
export const JWT_ACCESS_TOKEN_COOKIE_NAME = "JWT-ACCESS-TOKEN";
export const JWT_REFRESH_TOKEN_COOKIE_NAME = "JWT-REFRESH-TOKEN";

export const REGISTERED_USERS = [];

export function getRandomRegisteredUser() {
    return REGISTERED_USERS[
        Math.floor(Math.random() * REGISTERED_USERS.length)
    ];
}

export function getRandomEmail() {
    return `${randomString(15)}@${randomString(5)}.${randomString(5)}`;
}

export function getRandomPassword() {
    return randomString(15);
}

export function getWeightedElement(weightedArray) {
    const totalWeight = weightedArray.reduce(
        (prev, curr) => prev + curr.weight,
        0
    );
    const random = Math.random() * totalWeight;
    let sum = 0;
    for (const element of weightedArray) {
        sum += element.weight;
        if (random <= sum) {
            return element;
        }
    }
    return null;
}

export function initCases(cases) {
    for (const elemenet of cases) {
        const prefix = "iterations-";
        // counter has prefix to ensure they are grouped in the summary.
        elemenet.counter = new Counter(`${prefix}${elemenet.case.name}`);
    }
    return cases;
}

export function runCaseWeighted(weightedCases) {
    runTest(getWeightedElement(weightedCases));
}

export function runCaseByName(weightedCases, caseName) {
    for (const element of weightedCases) {
        if (element.case.name === caseName) {
            runTest(element);
        }
    }
}

let caseRoundRobinIndex = 0;
export function runCaseRoundRobin(weightedCases) {
    runTest(weightedCases[caseRoundRobinIndex++ % weightedCases.length]);
}

export function runTest(testCase) {
    testCase.case();
    testCase.counter.add(1);
}
