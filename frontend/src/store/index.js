import Vue from 'vue';
import Vuex from 'vuex';

import layout from './layout';
import users from './users';
import storages from "@/store/storages";
import user from './user';
import login from './login';
import storage from "./storage";
import jobs from "./jobs";
import job from "./job";
import satellites from "./satellites";
import satellite from "./satellite";

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        layout,
        users,
        storages,
        user,
        login,
        storage,
        jobs,
        job,
        satellites,
        satellite
    },
});
