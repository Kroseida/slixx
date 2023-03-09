import Vue from "vue";

const defaultState = () => ({
    subscribeId: -1,
    users: [],
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
        subscribeUsers(state, {callback, error}) {
            state.users = [];
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(`
            query {
                data: getUsers(limit: 10, search: "${state.table.search}", page: ${state.table.page}) {
                    rows  {
                      id
                      name
                      firstName
                      lastName
                      email
                      active
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
                state.users = data.rows;
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
        unsubscribeUsers(state) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
        }
    }
};
