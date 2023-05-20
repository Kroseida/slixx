import Vue from 'vue';
import Router from 'vue-router';

import Layout from '@/components/Layout/Layout';
import Login from '@/pages/Authentication/Login/Login';
import ErrorPage from '@/pages/Error/Error';

import Dashboard from "@/pages/Dashboard/Dashboard.vue";
import UserList from '@/pages/User/UserList/UserList';
import UserDetails from '@/pages/User/UserDetails/UserDetails';
import StorageList from "@/pages/Storage/StorageList/StorageList.vue";
import StorageDetails from "@/pages/Storage/StorageDetails/StorageDetails.vue";
import JobList from "@/pages/Job/JobList/JobList.vue";
import JobDetails from "@/pages/Job/JobDetails/JobDetails.vue";

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/login',
            name: 'Login',
            component: Login,
        },
        {
            path: '/error',
            name: 'Error',
            component: ErrorPage,
        },
        {
            path: '/app',
            name: 'Layout',
            component: Layout,
            children: [
                {
                    path: 'dashboard',
                    name: 'Dashboard',
                    component: Dashboard,
                },
                {
                    path: 'storage',
                    name: 'StorageList',
                    component: StorageList,
                },
                {
                    path: 'storage/:id',
                    name: 'StorageDetails',
                    component: StorageDetails,
                },
                {
                    path: 'user',
                    name: 'UserList',
                    component: UserList,
                },
                {
                    path: 'user/:id',
                    name: 'UserDetails',
                    component: UserDetails,
                },
                {
                    path: 'job',
                    name: 'JobList',
                    component: JobList,
                },
                {
                    path: 'job/:id',
                    name: 'JobDetails',
                    component: JobDetails,
                },
            ],
        },
    ],
});
