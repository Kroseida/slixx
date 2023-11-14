import {defineComponent} from 'vue'
import JobList from "components/JobListComponent/JobListComponent.vue";
import {useRouter} from 'vue-router'
import {useGlobalStore} from 'src/stores/global'

export default defineComponent({
  name: 'JobListPage',
  components: {
    JobList
  },
  data() {
    return {
      router: useRouter(),
      globalStore: useGlobalStore(),
    }
  },
  methods: {
    openJob(_, row) {
      this.router.push({path: `/job/${row.id}`})
    },
    createJob() {
      this.router.push({path: `/job/new/`})
    }
  }
})
