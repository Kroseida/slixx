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
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "getUser",
                args: {
                    id: userId,
                },
                fields: [
                    "id",
                    "name",
                    "firstName",
                    "lastName",
                    "email",
                    "active",
                    "createdAt",
                    "updatedAt",
                    "description",
                    "permissions"
                ],
                isMutation: false
            });

            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(fullQuery, (data) => {
                state.user = data;
                state.originalUser = JSON.parse(JSON.stringify(data));
            }, error);
        },
        deleteUser(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "deleteUser",
                args: {
                    id: state.user.id,
                },
                fields: [
                    "id",
                ],
                isMutation: true
            });

            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                state.user = data.message[0].data;
                callback(state.user);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        createUser(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "createUser",
                args: {
                    description: state.user.description.replaceAll('"', '\\"'),
                    email: state.user.email,
                    firstName: state.user.firstName,
                    lastName: state.user.lastName,
                    name: state.user.name,
                    active: state.user.active,
                },
                fields: [
                    "id",
                    "name",
                    "firstName",
                    "lastName",
                    "email",
                    "active",
                    "createdAt",
                    "updatedAt",
                    "description",
                    "permissions"
                ],
                isMutation: true
            });

            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                state.user = data.message[0].data;
                callback(state.user);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        saveUser(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "updateUser",
                args: {
                    id: state.user.id,
                    name: state.user.name !== state.originalUser.name ? state.user.name : undefined,
                    firstName: state.user.firstName !== state.originalUser.firstName ? state.user.firstName : undefined,
                    lastName: state.user.lastName !== state.originalUser.lastName ? state.user.lastName : undefined,
                    email: state.user.email !== state.originalUser.email ? state.user.email : undefined,
                    active: state.user.active !== state.originalUser.active ? (state.user.active === 'true') : undefined,
                    description: state.user.description !== state.originalUser.description ?
                        state.user.description.replaceAll("\\", "\\\\").replaceAll('"', '\\"') :
                        undefined,
                },
                fields: [
                    "id",
                    "name",
                    "firstName",
                    "lastName",
                    "email",
                    "active",
                    "createdAt",
                    "updatedAt",
                    "description",
                    "permissions"
                ],
                isMutation: true
            });
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = -1;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                callback(data);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                this.commit('user/subscribeUser', {
                    userId: state.user.id,
                    error
                });
            }, error);
        },
        addPermission(state, {permission, callback, error}) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = -1;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                data: addUserPermission(id: "${state.user.id}", permissions: ["${permission}"]) {
                    id
                }
            }
            `, (data) => {
                callback(data);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                this.commit('user/subscribeUser', {
                    userId: state.user.id,
                    error
                });
            }, error);
        },
        removePermission(state, {permission, callback, error}) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = -1;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                data: removeUserPermission(id: "${state.user.id}", permissions: ["${permission}"]) {
                    id
                }
            }
            `, (data) => {
                callback(data);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                this.commit('user/subscribeUser', {
                    userId: state.user.id,
                    error
                });
            }, error);
        },
        unsubscribeUser(state) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
        },
        updatePassword(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "createPasswordAuthentication",
                args: {
                    id: state.user.id,
                    password: state.authentication.password.value,
                },
                fields: [
                    "id",
                ],
                isMutation: true
            });

            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = -1;
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                callback(data);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
                this.commit('user/subscribeUser', {
                    userId: state.user.id,
                    error
                });
            }, error);
        }
    }
};
