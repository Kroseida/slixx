import Vue from "vue";

const defaultState = () => ({
    subscribeId: -1,
    users: [],
})

export default {
    namespaced: true,
    state: defaultState(),
    mutations: {
        reset(state) {
            const initial = defaultState()
            Object.keys(initial).forEach(key => { state[key] = initial[key] })
        },
        subscribeUsers(state, {callback, error}) {
            state.users = [];
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedArray(`
            query {
                users: getUsers {
                    id
                    name
                    firstName
                    lastName
                    email
                    active
                    createdAt
                    updatedAt
                }
            }
            `, (data) => {
                state.users = data;
                if (callback) {
                    callback();
                }
            }, error);
        },
        unsubscribeUsers(state) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
        }
    }
};
