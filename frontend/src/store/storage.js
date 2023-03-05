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
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(`
            query {
                data: getStorageKinds {
                    name
                    configuration {
                        name
                        kind
                    }
                }
            }
            `, (data) => {
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
        updateConfiguration(state, {fieldName, value}) {
            console.log(fieldName, value)
            state.storage.configuration[fieldName] = value;
        },
        subscribeStorage(state, {storageId, error, callback}) {
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(`
            query {
                data: getStorage(id: "${storageId}") {
                    id
                    name
                    kind
                    configuration
                }
            }
            `, (data) => {
                state.storage = data;
                state.storage.configuration = JSON.parse(state.storage.configuration);
                state.originalStorage = JSON.parse(JSON.stringify(state.storage));
                callback(state.storage);
            }, error);
        },
        deleteStorage(state, {callback, error}) {
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                data: deleteStorage(id: "${state.storage.id}") {
                    id
                }
            }
            `, (data) => {
                state.storage = data.message[0].data;
                callback(state.storage);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        createStorage(state, {callback, error}) {
            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(`
            mutation {
                data: createStorage(name: "${state.storage.name}", 
                                    description: "${state.storage.description}", 
                                    kind: "${state.storage.kind}", 
                                    configuration: "${JSON.stringify(state.storage.configuration).replaceAll("\\", "\\\\").replaceAll('"', '\\"')}") {
                    id
                    name
                    kind
                }
            }
            `, (data) => {
                state.storage = data.message[0].data;
                callback(state.storage);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        saveStorage(state, {callback, error}) {
            let fullQuery = 'mutation {';
            if (state.storage.name !== state.originalStorage.name) {
                fullQuery += `name: updateStorageName(id: "${state.storage.id}", name: "${state.storage.name}") {id name}`;
            }
            if (state.storage.kind !== state.originalStorage.kind) {
                fullQuery += `name: updateStorageKind(id: "${state.storage.id}", kind: "${state.storage.kind}") {id kind}`;
            }
            fullQuery += '}';

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
