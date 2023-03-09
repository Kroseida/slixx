import isScreen from '@/core/screenHelper';
import Vue from "vue";

export default {
    namespaced: true,
    state: {
        subscribeId: -1,
        permissionSubscribeId: -1,
        sidebarClose: true,
        sidebarStatic: false,
        sidebarActiveElement: null,
        localUser: {
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
        permissions: [],
        isPermitted(permission) {
            if (!this.localUser) {
                return false;
            }
            return this.localUser.permissions.includes(permission);
        }
    },
    mutations: {
        subscribePermissions(state, {callback}) {
            Vue.prototype.$graphql.unsubscribe(state.permissionSubscribeId);
            state.permissionSubscribeId = Vue.prototype.$graphql.subscribeTrackedObject(`
            query {
                data: getPermissions {
                    value
                    name
                }
            }
            `, (data) => {
                state.permissions = data;
                callback(data);
            });
        },
        subscribeLocalUser(state, {callback}) {
            Vue.prototype.$graphql.unsubscribe(state.subscribeId);
            state.subscribeId = Vue.prototype.$graphql.subscribeTrackedObject(`
            query {
                data: getLocalUser {
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
                state.localUser = data;
                callback(data);
            });
        },
        toggleSidebar(state) {
            const nextState = !state.sidebarStatic;

            localStorage.sidebarStatic = nextState;
            state.sidebarStatic = nextState;

            if (!nextState && (isScreen('lg') || isScreen('xl'))) {
                state.sidebarClose = true;
            }
        },
        switchSidebar(state, value) {
            if (value) {
                state.sidebarClose = value;
            } else {
                state.sidebarClose = !state.sidebarClose;
            }
        },
        handleSwipe(state, e) {
            if ('ontouchstart' in window) {
                if (e.direction === 4) {
                    state.sidebarClose = false;
                }

                if (e.direction === 2 && !state.sidebarClose) {
                    state.sidebarClose = true;
                }
            }
        },
        changeSidebarActive(state, index) {
            state.sidebarActiveElement = index;
        },
        afterAuthentication(state, {token}) {
            localStorage.setItem('token', token);
            Vue.prototype.$graphql.reconnect(token);
        }
    },
    actions: {
        toggleSidebar({commit}) {
            commit('toggleSidebar');
        },
        switchSidebar({commit}, value) {
            commit('switchSidebar', value);
        },
        handleSwipe({commit}, e) {
            commit('handleSwipe', e);
        },
        changeSidebarActive({commit}, index) {
            commit('changeSidebarActive', index);
        },
    },
};
