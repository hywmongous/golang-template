import { check } from "k6";
import { randomString } from "https://jslib.k6.io/k6-utils/1.1.0/index.js";
import http from "k6/http";
import { buildUrl, initCases, REGISTERED_USERS } from "../index.js";

export function successful_random_registration() {
    return successful_registration(
        `${randomString(15)}@${randomString(5)}.${randomString(5)}`,
        randomString(15)
    );
}

export function successful_registration(
    username = "some1@email",
    password = "Passw0rd"
) {
    const registration_url = buildUrl(
        `${username}:${password}`,
        "/api/v1/identities"
    );

    const registration_response = http.post(registration_url);

    check(registration_response, {
        "successful_registration: Status code is 201": (response) =>
            response.status === 201,
    });

    // Only if the user was created we will add it to the shared state
    if (registration_response.status === 201) {
        REGISTERED_USERS.push({
            username,
            password,
        });
    }

    return registration_response;
}

export const weightedCases = initCases([
    { weight: 100, case: successful_random_registration },
]);
