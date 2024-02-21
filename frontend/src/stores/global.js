import {defineStore} from 'pinia';

export const useGlobalStore = defineStore('global', {
  state: () => ({
    leftDrawerOpen: false,
    darkMode: false,
    counter: 0,
    session: '',
    connectionState: 'CONNECTING',
    connectedBefore: false,
    localUserSubscriptionId: -1,
    localUser: null,
    permissions: [],
    userLoaded: false,
  }),
  getters: {},
  actions: {
    isPermitted(permission) {
      if (!this.localUser) {
        return false;
      }
      return this.localUser.permissions.includes(permission);
    },
    setLeftDrawerOpen(leftDrawerOpen) {
      this.leftDrawerOpen = leftDrawerOpen;
    },
    setDarkMode(dark, darkMode) {
      this.darkMode = darkMode;
      localStorage.setItem('_darkMode', darkMode);
      dark.set(darkMode)
    },
    subscribeLocalUser(callback, error) {
      this.localUserSubscriptionId = this.$controller.user.subscribeLocalUser(
        this.localUserSubscriptionId,
        (response) => {
          this.userLoaded = true;
          this.localUser = response;
          if (!this.localUser) {
            callback(response);
          }
        },
        (data) => {
          if (!this.userLoaded) {
            error(data);
          }
        }
      );
    }
  },
});
