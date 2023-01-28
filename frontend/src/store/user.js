import Vue from "vue";

const defaultState = () => ({
    subscribeId: -1,
    updatedSubscriptionId: -1,
    user: {
        id: '',
        name: "",
        firstName: "",
        lastName: "",
        email: "",
        active: false,
        createdAt: "",
        updatedAt: "",
        description: "",
    },
    originalUser: {
        id: '',
        name: "",
        firstName: "",
        lastName: "",
        email: "",
        active: false,
        createdAt: "",
        updatedAt: "",
        description: "",
    },
    authentication: {
        password: {
            value: "",
            repeat: "",
        }
    }
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
        subscribeUser(state, {userId, error}) {
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(`
            query {
                user: getUser(id: "${userId}") {
                    id
                    name
                    firstName
                    lastName
                    email
                    active
                    createdAt
                    updatedAt
                    description
                    permissions
                }
            }
            `, (data) => {
                state.user = data;
                state.originalUser = JSON.parse(JSON.stringify(data));
            }, error);
        },
        createUser(state, {callback, error}) {
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                user: createUser(description: "${state.user.description}", email: "${state.user.email}", firstName: "${state.user.firstName}", lastName: "${state.user.lastName}", name: "${state.user.name}", active: ${state.user.active}) {
                    id
                    name
                    firstName
                    lastName
                    email
                    active
                    createdAt
                    updatedAt
                    description
                    permissions
                }
            }
            `, (data) => {
                state.user = data.message[0].user;
                callback(state.user);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        saveUser(state, {callback, error}) {
            let fullQuery = 'mutation {';
            if (state.user.name !== state.originalUser.name) {
                fullQuery += `name: updateUserName(id: "${state.user.id}", name: "${state.user.name}") {id name}`;
            }
            if (state.user.firstName !== state.originalUser.firstName) {
                fullQuery += `firstName: updateUserFirstName(id: "${state.user.id}", firstName: "${state.user.firstName}") {id firstName}`;
            }
            if (state.user.lastName !== state.originalUser.lastName) {
                fullQuery += `lastName: updateUserLastName(id: "${state.user.id}", lastName: "${state.user.lastName}") {id lastName}`;
            }
            if (state.user.email !== state.originalUser.email) {
                fullQuery += `email: updateUserEmail(id: "${state.user.id}", email: "${state.user.email}") {id email}`;
            }
            if (state.user.active !== state.originalUser.active) {
                fullQuery += `active: updateUserActive(id: "${state.user.id}", active: ${state.user.active}) {id active}`;
            }
            if (state.user.description !== state.originalUser.description) {
                fullQuery += `description: updateUserDescription(id: "${state.user.id}", description: "${state.user.description}") {id description}`;
            }
            fullQuery += '}';

            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = -1;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                callback(data);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                this.commit('user/subscribeUser', {
                    userId: state.user.id,
                });
            }, error);
        },
        addPermission(state, {permission, callback, error}) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = -1;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                permission: addUserPermission(id: "${state.user.id}", permissions: ["${permission}"]) {
                    id
                }
            }
            `, (data) => {
                callback(data);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                this.commit('user/subscribeUser', {
                    userId: state.user.id,
                });
            }, error);
        },
        removePermission(state, {permission, callback, error}) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = -1;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                permission: removeUserPermission(id: "${state.user.id}", permissions: ["${permission}"]) {
                    id
                }
            }
            `, (data) => {
                callback(data);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                this.commit('user/subscribeUser', {
                    userId: state.user.id,
                });
            }, error);
        },
        unsubscribeUser(state) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
        },
        updatePassword(state, {callback, error}) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = -1;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                password: createPasswordAuthentication(id: "${state.user.id}", password: "${state.authentication.password.value}") {
                    id
                }
            }
            `, (data) => {
                callback(data);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                this.commit('user/subscribeUser', {
                    userId: state.user.id,
                });
            }, error);
        }
    }
};
