import Vue from "vue";

const defaultState = () => ({
    updatedSubscriptionId: -1,
    authentication: {
        name: '',
        password: ''
    },
    loggingIn: false,
})

export default {
    namespaced: true,
    state: defaultState(),
    mutations: {
        reset(state) {
            const initial = defaultState()
            Object.keys(initial).forEach(key => {
                state[key] = initial[key]
            })
        },
        login(state, {callback, error}) {
            state.loggingIn = true;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                authentication: authenticate(name: "${state.authentication.name}", password: "${state.authentication.password}") {
                    id,
                    token
                }
            }
            `, (data) => {
                state.loggingIn = false;
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                callback(data);
            }, (data) => {
                state.loggingIn = false;
                error(data);
            });
        }
    }
};
