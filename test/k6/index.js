import { Counter } from "k6/metrics";

export const PORT = "80";
export const PROTOCOL = "http";
export const SUBDOMAIN = ""; // eg. 'www.'
export const ROOT_DOMAIN = "localhost";
export const TLD = ""; // eg. '.com'
export const BASE_URL = `${PROTOCOL}://${SUBDOMAIN}${ROOT_DOMAIN}${TLD}:${PORT}`;
export function buildUrl(basic) {
    return `${PROTOCOL}://${SUBDOMAIN}${basic}@${ROOT_DOMAIN}${TLD}:${PORT}`;
}

export const CSRF_HEADER_KEY = "Csrf";
export const JWT_ACCESS_TOKEN_COOKIE_NAME = "JWT-ACCESS-TOKEN";
export const JWT_REFRESH_TOKEN_COOKIE_NAME = "JWT-REFRESH-TOKEN";

export const REGISTERED_USERS = [
    { username: "some1@email", password: "P@ssw0rd" },
    { username: "some2@email", password: "P@ssw0rd" },
    { username: "some3@email", password: "P@ssw0rd" },
    { username: "some4@email", password: "P@ssw0rd" },
    { username: "some5@email", password: "P@ssw0rd" },
    { username: "some6@email", password: "P@ssw0rd" },
    { username: "some7@email", password: "P@ssw0rd" },
    { username: "some8@email", password: "P@ssw0rd" },
    { username: "some9@email", password: "P@ssw0rd" },
    { username: "some10@email", password: "P@ssw0rd" },
    { username: "some11@email", password: "P@ssw0rd" },
    { username: "some12@email", password: "P@ssw0rd" },
    { username: "some13@email", password: "P@ssw0rd" },
    { username: "some14@email", password: "P@ssw0rd" },
    { username: "some15@email", password: "P@ssw0rd" },
    { username: "some16@email", password: "P@ssw0rd" },
    { username: "some17@email", password: "P@ssw0rd" },
    { username: "some18@email", password: "P@ssw0rd" },
    { username: "some19@email", password: "P@ssw0rd" },
    { username: "some20@email", password: "P@ssw0rd" },
    { username: "some21@email", password: "P@ssw0rd" },
    { username: "some22@email", password: "P@ssw0rd" },
    { username: "some23@email", password: "P@ssw0rd" },
    { username: "some24@email", password: "P@ssw0rd" },
    { username: "some25@email", password: "P@ssw0rd" },
    { username: "some26@email", password: "P@ssw0rd" },
    { username: "some27@email", password: "P@ssw0rd" },
    { username: "some28@email", password: "P@ssw0rd" },
    { username: "some29@email", password: "P@ssw0rd" },
    { username: "some30@email", password: "P@ssw0rd" },
    { username: "some31@email", password: "P@ssw0rd" },
    { username: "some32@email", password: "P@ssw0rd" },
    { username: "some33@email", password: "P@ssw0rd" },
    { username: "some34@email", password: "P@ssw0rd" },
    { username: "some35@email", password: "P@ssw0rd" },
    { username: "some36@email", password: "P@ssw0rd" },
    { username: "some37@email", password: "P@ssw0rd" },
    { username: "some38@email", password: "P@ssw0rd" },
    { username: "some39@email", password: "P@ssw0rd" },
    { username: "some40@email", password: "P@ssw0rd" },
    { username: "some41@email", password: "P@ssw0rd" },
    { username: "some42@email", password: "P@ssw0rd" },
    { username: "some43@email", password: "P@ssw0rd" },
    { username: "some44@email", password: "P@ssw0rd" },
    { username: "some45@email", password: "P@ssw0rd" },
    { username: "some46@email", password: "P@ssw0rd" },
    { username: "some47@email", password: "P@ssw0rd" },
    { username: "some48@email", password: "P@ssw0rd" },
    { username: "some49@email", password: "P@ssw0rd" },
    { username: "some50@email", password: "P@ssw0rd" },
];

export function getRandomUser() {
    return REGISTERED_USERS[
        Math.floor(Math.random() * REGISTERED_USERS.length)
    ];
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

let caseRoundRobinIndex = 0;
export function runCaseRoundRobin(weightedCases) {
    runTest(weightedCases[caseRoundRobinIndex++ % weightedCases.length]);
}

export function runTest(testCase) {
    testCase.case();
    testCase.counter.add(1);
}
