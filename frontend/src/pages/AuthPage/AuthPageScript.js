import {defineComponent} from 'vue'
import {useGlobalStore} from "stores/global";
import {useRouter} from 'vue-router'
import {Notify} from "quasar";

export default defineComponent({
  name: 'IndexPage',
  data() {
    return {
      globalStore: useGlobalStore(),
      name: '',
      password: '',
      dismissNotification: null,
      authenticating: false,
      authTimeout: null,
      router: useRouter(),
    }
  },
  mounted() {
    if (this.globalStore.localUser !== null) {
      this.router.push('/');
    }
  },
  methods: {
    async authenticate() {
      if (this.authenticating) {
        return;
      }
      this.authenticating = true;
      this.dismissNotification = Notify.create({
        message: 'Authenticating...',
        timeout: 0,
        spinner: true,
        position: 'top',
      })
      this.$controller.user.authenticate(this.name, this.password, async (data) => {
        await this.afterAuthentication(data);
      }, this.afterAuthenticationError);
    },
    async afterAuthentication(data) {
      this.authenticating = false;
      this.globalStore.session = data.token;
      this.dismissNotification();
      Notify.create({
        message: 'Authentication successful',
        color: 'positive',
        position: 'top',
      })
      localStorage.setItem('_auth', data.token);
      await this.$controller.graphql.reconnect(data.token);
      if (data !== null) {
        await this.router.push('/');
      }
    },
    afterAuthenticationError(data) {
      this.authenticating = false;
      this.dismissNotification();
      Notify.create({
        message: 'Authentication failed: ' + data,
        color: 'negative',
        position: 'top',
      })
    }
  }
})
