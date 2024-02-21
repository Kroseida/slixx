<template>
  <div>
    <router-view
      v-if="globalStore.connectionState === 'CONNECTED'
      && globalStore.userLoaded
      && constantsStore.permissionLoaded
      && constantsStore.storageKindsLoaded
      && constantsStore.jobStrategiesLoaded
      && constantsStore.scheduleKindsLoaded
      && constantsStore.environmentLoaded
      && constantsStore.scheduleKindsLoaded"
    />
    <div v-else>
      <div class="text-center">
        <q-spinner size="50px" color="primary" style="margin-top: 250px"/>
        <div class="text-grey-9 text-h6 text-weight-bold" style="margin-top: 50px">Connecting to backend server</div>
      </div>
    </div>
  </div>
</template>

<script>
import {defineComponent} from 'vue'
import {useGlobalStore} from "stores/global";
import {useConstantsStore} from "stores/constants";
import {useRouter} from 'vue-router'
import {Notify} from "quasar";

export default defineComponent({
  name: 'App',
  setup() {
    return {
      globalStore: useGlobalStore(),
      router: useRouter(),
      constantsStore: useConstantsStore(),
    }
  },
  handleError(message) {
    Notify.create({
      message,
      color: 'negative',
      icon: 'report_problem',
      position: 'top',
    })
  },
  async mounted() {
    if (localStorage.getItem('_darkMode') == 'true') {
      this.globalStore.setDarkMode(this.$q.dark, true)
    } else {
      this.globalStore.setDarkMode(this.$q.dark, false)
    }

    this.$controller.graphql.onConnected = async () => {
      this.globalStore.connectionState = "CONNECTED";
      this.globalStore.connectedBefore = true;
      this.globalStore.subscribeLocalUser((user) => {
        this.constantsStore.subscribePermissions(this.handleError);
        this.constantsStore.subscribeStorageKinds(this.handleError);
        this.constantsStore.subscribeJobStrategies(this.handleError);
        this.constantsStore.subscribeJobScheduleKinds(this.handleError);
        this.constantsStore.subscribeEnvironment(this.handleError);

        const currentPath = this.router.options.history.location;

        if (currentPath.indexOf('/auth') !== -1 && user) {
          this.router.push('/');
        } else if (currentPath.indexOf('/auth') === -1 && !user) {
          this.router.push('/auth');
        }
      }, this.handleError);
    }

    this.$controller.graphql.onReset = () => {
      this.globalStore.connectionState = "CONNECTING";
    };

    this.$controller.graphql.onClose = () => {
      this.globalStore.connectionState = "CLOSED";
      this.globalStore.userLoaded = false;
      this.globalStore.permissionLoaded = false;
      this.globalStore.scheduleKindsLoaded = false;
      this.globalStore.environmentLoaded = false;
      Notify.create({
        message: 'Lost connection to backend server',
        color: 'negative',
        position: 'top-right',
        timeout: 0,
      })
    };

    setInterval(() => {
      if (this.globalStore.connectionState === 'CONNECTING') {
        this.$controller.graphql.reconnect(localStorage.getItem('_auth'));
      }
    }, 5000);

  }
})
</script>
