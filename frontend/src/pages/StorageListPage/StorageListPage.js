import {defineComponent} from 'vue'
import StorageList from "components/StorageListComponent/StorageListComponent.vue";
import {useRouter} from 'vue-router'
import {useGlobalStore} from 'src/stores/global'

export default defineComponent({
  name: 'StorageListPage',
  components: {
    StorageList
  },
  data() {
    return {
      router: useRouter(),
      globalStore: useGlobalStore(),
    }
  },
  methods: {
    openStorage(_, row) {
      this.router.push({path: `/storage/${row.id}`})
    },
    createStorage() {
      this.router.push({path: `/storage/new/`})
    }
  }
})
