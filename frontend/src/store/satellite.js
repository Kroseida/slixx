import Vue from "vue";

const defaultState = () => ({
    subscribeId: -1,
    updatedSubscriptionId: -1,
    satellite: {
        id: '',
        name: "",
        address: "",
        token: "",
        description: "",
        createdAt: "",
        updatedAt: "",
    },
    originalSatellite: {
        id: '',
        name: "",
        address: "",
        token: "",
        description: "",
        createdAt: "",
        updatedAt: "",
    },
    table: {
        page: 1,
        totalPages: 1,
        totalRows: 1
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
        subscribeSatellite(state, {satelliteId, error, callback}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "getSatellite",
                args: {
                    id: satelliteId
                },
                fields: [
                    "id",
                    "name",
                    "address",
                ],
                isMutation: false
            });
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(fullQuery, (data) => {
                state.satellite = data;
                state.originalSatellite = JSON.parse(JSON.stringify(state.satellite));
                if (callback) {
                    callback(state.satellite);
                }
            }, error);
        },
        deleteSatellite(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "deleteSatellite",
                args: {
                    id: state.satellite.id
                },
                fields: [
                    "id",
                ],
                isMutation: true
            });

            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                state.satellite = data.message[0].data;
                callback(state.satellite);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        createSatellite(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "createSatellite",
                args: {
                    name: state.satellite.name,
                    description: state.satellite.description,
                    address: state.satellite.address,
                    token: state.satellite.token,
                },
                fields: [
                    "id",
                    "name",
                    "address",
                ],
                isMutation: true
            });

            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                state.satellite = data.message[0].data;
                callback(state.satellite);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        saveSatellite(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "updateSatellite",
                args: {
                    id: state.satellite.id,
                    name: state.satellite.name === state.originalSatellite.name ? undefined : state.satellite.name,
                    description: state.satellite.description === state.originalSatellite.description ? undefined : state.satellite.description,
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
                this.commit('satellite/subscribeSatellite', {
                    satelliteId: state.satellite.id,
                });
            }, error);
        },
        unsubscribeSatellite(state) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
        },
    }
};
