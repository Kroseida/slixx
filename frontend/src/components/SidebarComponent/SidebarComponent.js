import SidebarLinkComponent from "components/SidebarLinkComponent/SidebarLinkComponent.vue";
import {useGlobalStore} from "stores/global"

export default {
  name: "SidebarComponent",
  components: {
    SidebarLinkComponent
  },
  setup() {
    return {
      globalStore: useGlobalStore(),
      leftDrawerOpen: false,
      links: [
        {
          title: 'Dashboard',
          icon: 'fas fa-dashboard',
          link: '/#/'
        },
        {
          title: 'Satellites',
          icon: 'fas fa-satellite',
          link: '/#/satellite',
          permission: 'satellite.view'
        },
        {
          title: 'Storages',
          icon: 'fas fa-hard-drive',
          link: '/#/storage',
          permission: 'storage.view'
        },
        {
          title: 'Jobs',
          icon: 'fas fa-table-list',
          link: '/#/job',
          permission: 'job.view'
        },
        {
          title: 'Users',
          icon: 'fas fa-users',
          link: '/#/user',
          permission: 'user.view'
        },
      ]
    };
  }
};
