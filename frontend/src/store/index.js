import Vue from 'vue';
import Vuex from 'vuex';

import layout from './layout';
import users from './users';
import user from './user';
import login from './login';

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        layout,
        users,
        user,
        login,
    },
});
