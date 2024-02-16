import {defineComponent} from 'vue'
import {useConstantsStore} from "stores/constants";
import {useGlobalStore} from "stores/global";
import {useRouter} from "vue-router";
import KindInputComponent from "components/KindedInputComponent/KindInputComponent.vue";
import SlixxButton from "components/SlixxButton/SlixxButton.vue";
import ButtonGroup from "components/ButtonGroup/ButtonGroup.vue";

export default defineComponent({
  name: 'JobScheduleCreate',
  components: {ButtonGroup, SlixxButton, KindInputComponent},
  data() {
    return {}
  },
  mounted() {
    if (this.initConfig) {
      this.scheduleKindByName(this.schedule.kind).configuration.forEach((config) => {
        // TODO: Dont do it like that. Its actually bad and can cause some bugs? For now its fine.
        // eslint-disable-next-line
        this.schedule.configuration[config.name] = config['default'];
      })
    }
  },
  setup() {
    return {
      globalStore: useGlobalStore(),
      router: useRouter(),
      constantsStore: useConstantsStore(),
    }
  },
  emits: ["s-click", "s-delete"],
  props: {
    schedule: {
      type: Object,
      required: true,
      default: () => ({
        name: "",
        description: "",
        kind: "CRON",
        configuration: {}
      })
    },
    initConfig: {
      type: Boolean,
      default: false
    },
    showDeleteButton: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    scheduleKindByName(name) {
      return this.constantsStore.scheduleKinds.find((strategy) => {
        return strategy.name === name;
      });
    },
    async onApply(done) {
      this.$emit('s-click', () => {
        done()
      });
    },
    async onDelete(done) {
      this.$emit('s-delete', () => {
        done()
      });
    }
  },
  computed: {
    kinds() {
      return this.constantsStore.scheduleKinds.map((kind) => {
        return kind.name;
      });
    },
  }
})
