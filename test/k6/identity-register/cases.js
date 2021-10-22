import { check } from "k6";
import http from "k6/http";
import {
    buildUrl,
    initCases,
    REGISTERED_USERS,
    getRandomEmail,
    getRandomPassword,
} from "../index.js";

export function successful_registration(
    username = getRandomEmail(),
    password = getRandomPassword()
) {
    const registration_url = buildUrl(
        `${username}:${password}`,
        "/api/v1/identities"
    );

    console.log("Registration");
    console.log(username);
    console.log(password);

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
    { weight: 100, case: successful_registration },
]);
