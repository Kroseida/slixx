import {defineComponent} from 'vue'
import BackupList from "components/BackupListComponent/BackupListComponent.vue";
import {useRouter} from 'vue-router'
import {useGlobalStore} from 'src/stores/global'

export default defineComponent({
  name: 'BackupListPage',
  components: {
    BackupList
  },
  data() {
    return {
      router: useRouter(),
      globalStore: useGlobalStore(),
    }
  },
  methods: {
    openBackup(_, row) {
      this.router.push({path: `/backup/${row.id}`})
    },
    createBackup() {
      this.router.push({path: `/backup/new/`})
    }
  }
})
