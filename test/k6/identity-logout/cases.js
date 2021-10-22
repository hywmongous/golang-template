import { check, sleep } from "k6";
import http from "k6/http";
import { successful_login } from "../identity-login/cases.js";
import {
    BASE_URL,
    CSRF_HEADER_KEY,
    initCases,
    getRandomEmail,
    getRandomPassword,
} from "../index.js";

export function successful_logout(
    username = getRandomEmail(),
    password = getRandomPassword()
) {
    const login_response = successful_login(username, password);

    // Mimics a user session before logging out
    sleep(0.1);

    const logout_url = `${BASE_URL}/api/v1/authentication/logout`;
    const logout_body = null;
    const logout_headers = {
        Csrf: login_response.headers[CSRF_HEADER_KEY],
    };

    const logout_response = http.post(logout_url, logout_body, {
        headers: logout_headers,
    });

    check(logout_response, {
        "successful_logout: Status code is 200": (response) =>
            response.status === 200,
    });

    return logout_response;
}

export const weightedCases = initCases([
    { weight: 100, case: successful_logout },
]);
