import {defineComponent} from 'vue'
import UserList from "components/UserListComponent/UserListComponent.vue";
import {useRouter} from 'vue-router'
import {useGlobalStore} from 'src/stores/global'

export default defineComponent({
  name: 'UserListPage',
  components: {
    UserList
  },
  data() {
    return {
      router: useRouter(),
      globalStore: useGlobalStore(),
    }
  },
  methods: {
    openUser(_, row) {
      this.router.push({path: `/user/${row.id}`})
    },
    createUser() {
      this.router.push({path: `/user/new/`})
    }
  }
})
