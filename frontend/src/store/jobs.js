import Vue from "vue";

const defaultState = () => ({
    subscribeId: -1,
    jobs: [],
    table: {
        search: "",
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
            Object.keys(initial).forEach(key => { state[key] = initial[key] })
        },
        subscribeJobs(state, {callback, error}) {
            state.jobs = [];
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(`
            query {
                data: getJobs(limit: 10, search: "${state.table.search}", page: ${state.table.page}) {
                    rows  {
                      id
                      name
                      strategy
                      createdAt
                      updatedAt
                    }
                    page {
                      totalRows
                      totalPages
                    } 
                }
            }
            `, (data) => {
                state.jobs = data.rows;
                state.table.totalPages = data.page.totalPages;
                if (state.table.totalPages === 0) {
                    state.table.totalPages = 1;
                }
                state.table.totalRows = data.page.totalRows;
                if (callback) {
                    callback();
                }
            }, error);
        },
        unsubscribeJobs(state) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
        }
    }
};