<template>
  <b-navbar class="header d-print-none app-header">
    <b-nav>
      <b-nav-item>
        <a class="px-1" href="#" @click="toggleSidebarMethod" id="barsTooltip">
          <i class='fa fa-bars' style="font-size: 21px"/>
        </a>
      </b-nav-item>
    </b-nav>
    <b-nav class="ml-auto">
      <b-nav-item-dropdown
          class="notificationsMenu d-sm-down-none mr-2"
          menu-class="notificationsWrapper py-0 animate__animated animate__animated-fast animate__fadeIn"
          right>
        <template slot="button-content">
          <span class="avatar rounded-circle mr-2">
            <i class="fa fa-user" style="margin-left: 2px"/>
          </span>
          <span class="px-2" id="localUser__name">{{ localUser.name }}</span>
          <i class='fi flaticon-arrow-down px-2'/>
        </template>
        <Notifications/>
      </b-nav-item-dropdown>
    </b-nav>
  </b-navbar>
</template>

<script>
import {mapState, mapActions} from 'vuex';
import Notifications from '@/components/Notifications/Notifications';

export default {
  name: 'Header',
  components: {Notifications},
  computed: {
    ...mapState('layout', ['sidebarClose', 'sidebarStatic', 'localUser']),
  },
  methods: {
    ...mapActions('layout', ['toggleSidebar', 'switchSidebar', 'changeSidebarActive']),
    switchSidebarMethod() {
      if (!this.sidebarClose) {
        this.switchSidebar(true);
        this.changeSidebarActive(null);
      } else {
        this.switchSidebar(false);
        const paths = this.$route.fullPath.split('/');
        paths.pop();
        this.changeSidebarActive(paths.join('/'));
      }
    },
    toggleSidebarMethod() {
      if (this.sidebarStatic) {
        this.toggleSidebar();
        this.changeSidebarActive(null);
      } else {
        this.toggleSidebar();
        const paths = this.$route.fullPath.split('/');
        paths.pop();
        this.changeSidebarActive(paths.join('/'));
      }
    },
    logout() {
      window.localStorage.setItem('authenticated', false);
      this.$router.push('/login');
    },
  }
};
</script>

<style src="./Header.scss" lang="scss"></style>
