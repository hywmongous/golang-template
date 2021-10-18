export { default as identity_login } from "./identity-login/smoke_test.js";
export { default as identity_logout } from "./identity-logout/smoke_test.js";

const minDuration = 1;

export const options = {
    minIterationDuration: `${minDuration}s`,
    scenarios: {
        identity_login: createSmokeScenario("identity_login"),
        identity_logout: createSmokeScenario("identity_logout"),
    },
};

function createSmokeScenario(func) {
    return {
        exec: func,
        executor: "constant-vus",
        vus: 1,
        duration: "10s",
        startTime: `${Math.random() * minDuration}s`,
        gracefulStop: "5s",
    };
}
