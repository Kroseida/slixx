import SelectableComponent from "components/SelectableComponent/SelectableComponent.vue";
import SatelliteListComponent from "components/SatelliteListComponent/SatelliteListComponent.vue";

export default {
  name: "SatelliteSelectableComponent",
  components: {
    SelectableComponent,
    SatelliteListComponent,
  },
  data() {
    return {
      selectedJob: null,
    };
  },
  props: {
    label: {
      type: String,
      default: "Satellite",
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
    this.$controller.satellite.subscribeSatellite(
      -1,
      this.value,
      (satellite, subscribeId) => {
        this.$controller.unsubscribe(subscribeId)
        this.selectedJob = satellite;
      },
      (e, subscribeId) => {
        this.$controller.unsubscribe(subscribeId)
        this.$q.notify({
          message: "Error while fetching satellite: " + e,
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
