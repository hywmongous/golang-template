import { check } from 'k6';
import http from 'k6/http';
import {
    BASE_URL,
    CSRF_HEADER_KEY,
    JWT_ACCESS_TOKEN_COOKIE_NAME,
    JWT_REFRESH_TOKEN_COOKIE_NAME,
    buildUrl,
    initCases,
    REGISTERED_USERS,
    getUser,
} from '../index.js'

export function random_successfull_login() {
    const user = getUser()
    return successfull_login(user.username, user.password)
}

export function successfull_login(
    username = 'some@email',
    password = 'P@ssw0rd',
) {
    const base = buildUrl(`${username}:${password}`);
    const url = `${base}/api/v1/authentication/login`;

    const login_response = http.post(url);

    check(login_response, {
        'successfull_login: Status code is 200': (response) => response.status === 200,
        'successfull_login: Contains CSRF header': (response) => response.headers[CSRF_HEADER_KEY] !== undefined,
        'successfull_login: Contains cookie acces token': (response) => response.cookies[JWT_ACCESS_TOKEN_COOKIE_NAME] !== undefined,
        'successfull_login: Contains cookie refresh token': (response) => response.cookies[JWT_REFRESH_TOKEN_COOKIE_NAME] !== undefined,
    });

    return login_response
};

export function invalid_credentials_login() {
    const username = '';
    const password = '';
    const base = buildUrl(`${username}:${password}`);
    const url = `${base}/api/v1/authentication/login`;

    const login_response = http.post(url);

    check(login_response, {
        'invalid_credentials_login: Status code is 401': (response) => response.status === 401,
        'invalid_credentials_login: Does not contain CSRF header': (response) => response.headers[CSRF_HEADER_KEY] === undefined,
        'invalid_credentials_login: Does not contain cookie acces token': (response) => response.cookies[JWT_ACCESS_TOKEN_COOKIE_NAME] === undefined,
        'invalid_credentials_login: Does not contain cookie refresh token': (response) => response.cookies[JWT_REFRESH_TOKEN_COOKIE_NAME] === undefined,
    });

    return login_response
}

export function missing_credentials_login() {
    const url = `${BASE_URL}/api/v1/authentication/login`;

    const login_response = http.post(url);

    check(login_response, {
        'missing_credentials_login: Status code is 401': (response) => response.status === 401,
        'missing_credentials_login: Does not contain CSRF header': (response) => response.headers[CSRF_HEADER_KEY] === undefined,
        'missing_credentials_login: Does not contain cookie acces token': (response) => response.cookies[JWT_ACCESS_TOKEN_COOKIE_NAME] === undefined,
        'missing_credentials_login: Does not contain cookie refresh token': (response) => response.cookies[JWT_REFRESH_TOKEN_COOKIE_NAME] === undefined,
    });

    return login_response
}

export const weightedCases = initCases([
    { weight: 85, case: random_successfull_login, },
    { weight: 10, case: invalid_credentials_login, },
    { weight:  5, case: missing_credentials_login, },
]);
