import { Counter } from 'k6/metrics'

export const PORT = '8080';
export const PROTOCOL = 'http://';
export const SUBDOMAIN = ''; // eg. 'www.'
export const ROOT_DOMAIN = 'localhost';
export const TLD = ''; // eg. '.com'
export const BASE_URL = `${PROTOCOL}${SUBDOMAIN}${ROOT_DOMAIN}${TLD}:${PORT}`;
export function buildUrl(basic) { return `${PROTOCOL}${SUBDOMAIN}${basic}@${ROOT_DOMAIN}${TLD}:${PORT}`; }

export const CSRF_HEADER_KEY = 'Csrf'
export const JWT_ACCESS_TOKEN_COOKIE_NAME = 'JWT-ACCESS-TOKEN'
export const JWT_REFRESH_TOKEN_COOKIE_NAME = 'JWT-REFRESH-TOKEN'

export function getWeightedElement(weightedArray) {
    const totalWeight = weightedArray.reduce((prev, curr) => prev + curr.weight, 0);
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
        const prefix = 'iterations-';
        // counter has prefix to ensure they are grouped in the summary
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

let roundRobinIndex = 0;
export function runCaseRoundRobin(weightedCases) {
    if (roundRobinIndex >= weightedCases.length) {
        roundRobinIndex = 0;
    }
    runTest(weightedCases[roundRobinIndex++]);
}

export function runTest(testCase) {
    testCase.case();
    testCase.counter.add(1);
}
