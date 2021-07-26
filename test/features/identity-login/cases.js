import { check } from 'k6';
import http from 'k6/http';
import {
    BASE_URL,
    CSRF_HEADER_KEY,
    JWT_ACCESS_TOKEN_COOKIE_NAME,
    JWT_REFRESH_TOKEN_COOKIE_NAME,
    buildUrl,
    initCases,
} from '../index.js'

export function successfull_login() {
    const username = 'username';
    const password = 'password';
    const base = buildUrl(`${username}:${password}`);
    const url = `${base}/api/v1/authentication/login`;

    const loginResponse = http.post(url);

    check(loginResponse, {
        'Status code is 200': (response) => response.status === 200,
        'Contains CSRF header': (response) => response.headers[CSRF_HEADER_KEY] !== undefined,
        'Contains cookie acces token': (response) => response.cookies[JWT_ACCESS_TOKEN_COOKIE_NAME] !== undefined,
        'Contains cookie refresh token': (response) => response.cookies[JWT_REFRESH_TOKEN_COOKIE_NAME] !== undefined,
    });
};

export function invalid_credentials_login() {
    const username = '';
    const password = '';
    const base = buildUrl(`${username}:${password}`);
    const url = `${base}/api/v1/authentication/login`;

    const loginResponse = http.post(url);

    check(loginResponse, {
        'Status code is 401': (response) => response.status === 401,
        'Does not contain CSRF header': (response) => response.headers[CSRF_HEADER_KEY] === undefined,
        'Does not contain cookie acces token': (response) => response.cookies[JWT_ACCESS_TOKEN_COOKIE_NAME] === undefined,
        'Does not contain cookie refresh token': (response) => response.cookies[JWT_REFRESH_TOKEN_COOKIE_NAME] === undefined,
    });
}

export function missing_credentials_login() {
    const url = `${BASE_URL}/api/v1/authentication/login`;

    const loginResponse = http.post(url);

    check(loginResponse, {
        'Status code is 401': (response) => response.status === 401,
        'Does not contain CSRF header': (response) => response.headers[CSRF_HEADER_KEY] === undefined,
        'Does not contain cookie acces token': (response) => response.cookies[JWT_ACCESS_TOKEN_COOKIE_NAME] === undefined,
        'Does not contain cookie refresh token': (response) => response.cookies[JWT_REFRESH_TOKEN_COOKIE_NAME] === undefined,
    });
}

export const weightedCases = initCases([
    { weight: 85, case: successfull_login, },
    { weight: 10, case: invalid_credentials_login, },
    { weight:  5, case: missing_credentials_login, },
]);
