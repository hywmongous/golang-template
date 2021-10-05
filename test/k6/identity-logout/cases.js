import { check, sleep } from 'k6';
import http from 'k6/http';
import {
    successfull_login
} from '../identity-login/cases.js'
import {
    BASE_URL,
    CSRF_HEADER_KEY,
    JWT_ACCESS_TOKEN_COOKIE_NAME,
    JWT_REFRESH_TOKEN_COOKIE_NAME,
    buildUrl,
    initCases,
} from '../index.js'

export function successfull_logout(
    username = 'some1@email',
    password = 'P@ssw0rd',
) {
    const login_response = successfull_login(username, password)

    const logout_url = `${BASE_URL}/api/v1/authentication/logout`
    const logout_body = null
    const logout_headers = {
        "Csrf": login_response.headers[CSRF_HEADER_KEY]
    }

    const logout_response = http.post(logout_url, logout_body, {
        headers: logout_headers,
    })

    check(logout_response, {
        'successfull_logout: Status code is 200': (response) => response.status === 200,
    })

    return logout_response
}

export const weightedCases = initCases([
    { weight: 100, case: successfull_logout, },
]);
