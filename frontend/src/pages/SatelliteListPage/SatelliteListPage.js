import {defineComponent} from 'vue'
import SatelliteList from "components/SatelliteListComponent/SatelliteListComponent.vue";
import {useRouter} from 'vue-router'
import {useGlobalStore} from 'src/stores/global'

export default defineComponent({
  name: 'SatelliteListPage',
  components: {
    SatelliteList
  },
  data() {
    return {
      router: useRouter(),
      globalStore: useGlobalStore(),
    }
  },
  methods: {
    openSatellite(_, row) {
      this.router.push({path: `/satellite/${row.id}`})
    },
    createSatellite() {
      this.router.push({path: `/satellite/new/`})
    }
  }
})
