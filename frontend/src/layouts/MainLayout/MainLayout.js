import {defineComponent, ref} from 'vue'
import SidebarComponent from "components/SidebarComponent/SidebarComponent.vue";
import {useGlobalStore} from "stores/global"
import {useConstantsStore} from "stores/constants";
export default defineComponent({
  name: 'MainLayout',

  components: {
    SidebarComponent
  },

  setup() {
    return {
      globalStore: useGlobalStore(),
      constantsStore: useConstantsStore(),
    }
  },
  methods: {
    toggleLeftDrawer() {
      this.globalStore.setLeftDrawerOpen(!this.globalStore.leftDrawerOpen)
    }
  }
})
