import {defineComponent} from 'vue'
import {useConstantsStore} from "stores/constants";
import {useGlobalStore} from "stores/global";
import {useRouter} from "vue-router";
import KindInputComponent from "components/KindedInputComponent/KindInputComponent.vue";

export default defineComponent({
  name: 'JobScheduleCreate',
  components: {KindInputComponent},
  data() {
    return {
      schedule: {
        kind: "CRON",
        configuration: {}
      }
    }
  },
  setup() {
    return {
      globalStore: useGlobalStore(),
      router: useRouter(),
      constantsStore: useConstantsStore(),
    }
  },
  props: {
    jobId: {
      type: String,
      required: true
    }
  },
  methods: {
    scheduleKindByName(name) {
      return this.constantsStore.scheduleKinds.find((strategy) => {
        return strategy.name === name;
      });
    }
  },
  computed: {
    kinds() {
      return this.constantsStore.scheduleKinds.map((kind) => {
        return kind.name;
      });
    }
  }
})
