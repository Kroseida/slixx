<template>
  <div>
    <router-view v-if="connectionState === 'CONNECTED' && userLoaded && permissionLoaded"/>
    <div v-if="connectionState === 'CONNECTING'">
      <div class="text-center">
        <div class="center" role="status">
          <div class="spinner-border text-primary">
            <span class="sr-only">Loading...</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "App",
  computed: {},
  mounted() {
    this.$graphql.onConnected = () => {
      this.connectionState = "CONNECTED";
      this.$store.commit('layout/subscribeLocalUser', {
        callback: (user) => {
          const currentPath = this.$router.history.current.path;
          if (currentPath === '/login' && user) {
            this.$router.push('/app/dashboard');
          } else if (currentPath !== '/login' && !user) {
            this.$router.push('/login');
          }
          this.userLoaded = true;
        }
      });
      this.$store.commit('layout/subscribePermissions', {
        callback: () => {
          this.permissionLoaded = true;
        }
      });
    };
    this.$graphql.onReset = () => {
      this.connectionState = "CONNECTING";
    };
    this.$graphql.onClose = () => {
      this.connectionState = "CLOSED";
      this.userLoaded = false;
      this.$toasted.error('Lost connection to backend server', {
        duration: 5000,
        position: 'top-right',
        fullWidth: true,
        fitToScreen: true,
      });
    };
    setInterval(() => {
      if(this.connectionState === 'CONNECTING') {
        this.$graphql.reconnect(localStorage.getItem('token'));
      }
    }, 5000);
  },
  data() {
    return {
      connectionState: "CONNECTING",
      connectedBefore: false,
      userLoaded: false,
      permissionLoaded: false,
    };
  },
};
</script>

<style src="./styles/theme.scss" lang="scss"/>
