import Vue from "vue";

const defaultState = () => ({
    subscribeId: -1,
    updatedSubscriptionId: -1,
    storage: {
        id: '',
        name: "",
        kind: "",
        description: "",
        createdAt: "",
        updatedAt: "",
        configuration: {},
    },
    originalStorage: {
        id: '',
        name: "",
        kind: "",
        description: "",
        createdAt: "",
        updatedAt: "",
        configuration: {},
    },
    kinds: {}
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
        subscribeKinds(state, {error, callback}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "getStorageKinds",
                args: [],
                fields: [
                    "name",
                    {
                        name: "configuration",
                        sub: [
                            "name",
                            "kind"
                        ]
                    }
                ],
                isMutation: false
            });
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(fullQuery, (data) => {
                state.kinds = {};
                data.forEach(kind => {
                    let configuration = {};
                    kind.configuration.forEach(config => {
                        configuration[config.name] = config.kind;
                    });
                    state.kinds[kind.name] = configuration;
                });
                callback(state.kinds);
            }, error);
        },
        subscribeStorage(state, {storageId, error, callback}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "getStorage",
                args: {
                    id: storageId
                },
                fields: [
                    "id",
                    "name",
                    "kind",
                    "configuration",
                ],
                isMutation: false
            });
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(fullQuery, (data) => {
                state.storage = data;
                state.storage.configuration = JSON.parse(state.storage.configuration);
                state.originalStorage = JSON.parse(JSON.stringify(state.storage));
                if (callback) {
                    callback(state.storage);
                }
            }, error);
        },
        deleteStorage(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "deleteStorage",
                args: {
                    id: state.storage.id
                },
                fields: [
                    "id",
                ],
                isMutation: true
            });

            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                state.storage = data.message[0].data;
                callback(state.storage);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        createStorage(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "createStorage",
                args: {
                    name: state.storage.name,
                    description: state.storage.description,
                    kind: state.storage.kind,
                    configuration: JSON.stringify(state.storage.configuration)
                        .replaceAll("\\", "\\\\")
                        .replaceAll('"', '\\"')
                },
                fields: [
                    "id",
                    "name",
                    "kind",
                ],
                isMutation: true
            });

            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                state.storage = data.message[0].data;
                callback(state.storage);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        saveConfiguration(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "updateStorage",
                args: {
                    id: state.storage.id,
                    configuration: JSON.stringify(state.storage.configuration)
                        .replaceAll("\\", "\\\\")
                        .replaceAll('"', '\\"')
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
                this.commit('storage/subscribeStorage', {
                    storageId: state.storage.id,
                });
            }, error);
        },
        saveStorage(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "updateStorage",
                args: {
                    id: state.storage.id,
                    name: state.storage.name === state.originalStorage.name ? undefined : state.storage.name,
                    description: state.storage.description === state.originalStorage.description ? undefined : state.storage.description,
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
                this.commit('storage/subscribeStorage', {
                    storageId: state.storage.id,
                });
            }, error);
        },
        unsubscribeStorage(state) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
        },
    }
};
