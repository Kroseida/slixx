import SelectableComponent from "components/SelectableComponent/SelectableComponent.vue";
import StorageListComponent from "components/StorageListComponent/StorageListComponent.vue";

export default {
  name: "StorageSelectableComponent",
  components: {
    SelectableComponent,
    StorageListComponent,
  },
  data() {
    return {
      selectedJob: null,
    };
  },
  props: {
    label: {
      type: String,
      default: "Storage",
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    value: {
      type: String,
      default: "",
    },
  },
  mounted() {
    if (!this.value) {
      return
    }
    this.$controller.storage.subscribeStorage(
      -1,
      this.value,
      (storage, subscribeId) => {
        this.$controller.unsubscribe(subscribeId)
        this.selectedJob = storage;
      },
      (e, subscribeId) => {
        this.$controller.unsubscribe(subscribeId)
        this.$q.notify({
          message: "Error while fetching storage: " + e,
          color: "negative",
          icon: "report_problem",
        });
      }
    )
  },
  methods: {
    selectJob(_, job) {
      this.selectedJob = job;
      this.$emit('input', job);
    }
  }
};
