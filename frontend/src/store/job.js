import Vue from "vue";

const defaultState = () => ({
    subscribeId: -1,
    updatedSubscriptionId: -1,
    job: {
        id: '',
        name: "",
        strategy: "",
        description: "",
        createdAt: "",
        updatedAt: "",
        configuration: {},
        originStorageId: "",
        destinationStorageId: "",
    },
    originalJob: {
        id: '',
        name: "",
        strategy: "",
        description: "",
        createdAt: "",
        updatedAt: "",
        configuration: {},
        originStorageId: "",
        destinationStorageId: "",
    },
    strategies: {},
    storages: [],
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
        subscribeStorages(state, {callback, error, filter}) {
            state.storages = [];
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(`
            query {
                data: getStorages(limit: 10, search: "${filter}", page: ${state.table.page}) {
                    rows  {
                      id
                      name
                    }
                    page {
                      totalRows
                      totalPages
                    } 
                }
            }
            `, (data) => {
                state.storages = data.rows;
                state.table.totalPages = data.page.totalPages;
                if (state.table.totalPages === 0) {
                    state.table.totalPages = 1;
                }
                state.table.totalRows = data.page.totalRows;
                if (callback) {
                    callback(data);
                }
            }, error);
        },
        subscribeStrategies(state, {error, callback}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "getJobStrategies",
                args: [],
                fields: [
                    "name",
                    {
                        name: "configuration",
                        sub: [
                            "name",
                            "kind",
                            "default"
                        ]
                    }
                ],
                isMutation: false
            });
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(fullQuery, (data) => {
                state.strategies = {};
                data.forEach(strategy => {
                    let configuration = {};
                    strategy.configuration.forEach(config => {
                        configuration[config.name] = {
                            kind: config.kind,
                            default: config.default
                        };
                    });
                    state.strategies[strategy.name] = configuration;
                });
                callback(state.strategies);
            }, error);
        },
        subscribeJob(state, {jobId, error, callback}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "getJob",
                args: {
                    id: jobId
                },
                fields: [
                    "id",
                    "name",
                    "strategy",
                    "configuration",
                    "originStorageId",
                    "destinationStorageId",
                ],
                isMutation: false
            });
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(fullQuery, (data) => {
                state.job = data;
                state.job.configuration = JSON.parse(state.job.configuration);
                state.originalJob = JSON.parse(JSON.stringify(state.job));
                if (callback) {
                    callback(state.job);
                }
            }, error);
        },
        deleteJob(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "deleteJob",
                args: {
                    id: state.job.id
                },
                fields: [
                    "id",
                ],
                isMutation: true
            });

            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                state.job = data.message[0].data;
                callback(state.job);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        createJob(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "createJob",
                args: {
                    name: state.job.name,
                    description: state.job.description,
                    strategy: state.job.strategy,
                    configuration: JSON.stringify(state.job.configuration)
                        .replaceAll("\\", "\\\\")
                        .replaceAll('"', '\\"'),
                    originStorageId: state.job.originStorageId,
                    destinationStorageId: state.job.destinationStorageId,
                },
                fields: [
                    "id",
                    "name",
                    "strategy",
                ],
                isMutation: true
            });

            state.updatedSubscriptionId = Vue.prototype.$graphql.subscribe(fullQuery, (data) => {
                state.job = data.message[0].data;
                callback(state.job);
                Vue.prototype.$graphql.unsubscribe(state.updatedSubscriptionId);
            }, error);
        },
        saveConfiguration(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "updateJob",
                args: {
                    id: state.job.id,
                    configuration: JSON.stringify(state.job.configuration)
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
                this.commit('job/subscribeJob', {
                    jobId: state.job.id,
                });
            }, error);
        },
        saveJob(state, {callback, error}) {
            let fullQuery = Vue.prototype.$graphql.buildQuery({
                method: "updateJob",
                args: {
                    id: state.job.id,
                    name: state.job.name === state.originalJob.name ? undefined : state.job.name,
                    description: state.job.description === state.originalJob.description ? undefined : state.job.description,
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
                this.commit('job/subscribeJob', {
                    jobId: state.job.id,
                });
            }, error);
        },
        unsubscribeJob(state) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
        },
    }
};
