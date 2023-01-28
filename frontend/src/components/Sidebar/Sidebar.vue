<template>
  <div class="sidebar-wrapper">
    <nav
        :class="{sidebar: true, sidebarStatic, sidebarOpened}"
        style="min-height: 100vh"
        @mouseenter="sidebarMouseEnter"
        @mouseleave="sidebarMouseLeave"
    >
      <header class="logo">
        <router-link to="/app/dashboard"><span class="primary-word">SLIXX</span><span class="secondary-word"> Dashboard</span></router-link>
      </header>
      <ul class="nav">
        <NavLink
            :activeItem="activeItem"
            header="Dashboard"
            link="/app/dashboard"
            iconName="fa fa-dashboard"
            index="dashboard"
            id="sidebar__dashboard"
            isHeader
        />
        <NavLink
            :activeItem="activeItem"
            header="Typography"
            link="/app/typography"
            iconName="fa fa-font"
            index="typography"
            isHeader
        />
        <NavLink
            :activeItem="activeItem"
            header="Tables Basic"
            link="/app/tables"
            iconName="fa fa-table"
            index="tables"
            isHeader
        />
        <NavLink
            :activeItem="activeItem"
            header="Notifications"
            link="/app/notifications"
            iconName="fa fa-bell"
            index="notifications"
            isHeader
        />
        <NavLink
            :activeItem="activeItem"
            header="Components"
            link="/app/components"
            iconName="fa fa-code"
            index="components"
            :childrenLinks="[
              { header: 'Charts', link: '/app/components/charts' },
              { header: 'Icons', link: '/app/components/icons' },
              { header: 'Maps', link: '/app/components/maps' },
            ]"
        />
        <NavLink
            :activeItem="activeItem"
            v-if="isPermitted('user.view')"
            header="Users"
            link="/app/user"
            iconName="fa fa-user"
            index="notifications"
            isHeader
            id="sidebar__user"
        />
      </ul>
    </nav>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import isScreen from '@/core/screenHelper';
import NavLink from './NavLink/NavLink';

export default {
  name: 'Sidebar',
  components: { NavLink },
  data() {
    return {
      alerts: [
        {
          id: 0,
          title: 'Sales Report',
          value: 15,
          footer: 'Calculating x-axis bias... 65%',
          color: 'danger',
        },
        {
          id: 1,
          title: 'Personal Responsibility',
          value: 20,
          footer: 'Provide required notes',
          color: 'primary',
        },
      ],
    };
  },
  methods: {
    ...mapActions('layout', ['changeSidebarActive', 'switchSidebar']),
    setActiveByRoute() {
      const paths = this.$route.fullPath.split('/');
      paths.pop();
      this.changeSidebarActive(paths.join('/'));
    },
    sidebarMouseEnter() {
      if (!this.sidebarStatic && (isScreen('lg') || isScreen('xl'))) {
        this.switchSidebar(false);
        this.setActiveByRoute();
      }
    },
    sidebarMouseLeave() {
      if (!this.sidebarStatic && (isScreen('lg') || isScreen('xl'))) {
        this.switchSidebar(true);
        this.changeSidebarActive(null);
      }
    },
  },
  created() {
    this.setActiveByRoute();
  },
  computed: {
    ...mapState('layout', {
      sidebarStatic: state => state.sidebarStatic,
      sidebarOpened: state => !state.sidebarClose,
      activeItem: state => state.sidebarActiveElement
    }),
    ...mapState('layout', ['localUser', 'isPermitted']),
  },
};
</script>

<!-- Sidebar styles should be scoped -->
<style src="./Sidebar.scss" lang="scss" scoped/>
