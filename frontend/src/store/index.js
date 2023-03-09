import Vue from 'vue';
import Vuex from 'vuex';

import layout from './layout';
import users from './users';
import storages from "@/store/storages";
import user from './user';
import login from './login';
import storage from "./storage";

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        layout,
        users,
        storages,
        user,
        login,
        storage
    },
});
