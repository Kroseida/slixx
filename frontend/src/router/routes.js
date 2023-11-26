
const routes = [
  {
    path: '/auth',
    component: () => import('pages/AuthPage/AuthPage.vue')
  },
  {
    path: '/',
    component: () => import('layouts/MainLayout/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/IndexPage.vue') },
      { path: 'user', component: () => import('pages/UserListPage/UserListPage.vue') },
      { path: 'user/:id', component: () => import('pages/UserDetailsPage/UserDetailsPage.vue') },
      { path: 'storage', component: () => import('pages/StorageListPage/StorageListPage.vue') },
      { path: 'storage/:id', component: () => import('pages/StorageDetailsPage/StorageDetailsPage.vue') },
      { path: 'job', component: () => import('pages/JobListPage/JobListPage.vue') },
      { path: 'job/:id', component: () => import('pages/JobDetailsPage/JobDetailsPage.vue') },
      { path: 'satellite', component: () => import('pages/SatelliteListPage/SatelliteListPage.vue') },
      { path: 'satellite/:id', component: () => import('pages/SatelliteDetailsPage/SatelliteDetailsPage.vue') },
    ]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue')
  }
]

export default routes
