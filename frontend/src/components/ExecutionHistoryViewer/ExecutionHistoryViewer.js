import {Notify} from "quasar";

export default {
  name: "ExecutionHistoryViewer",
  setup() {},
  data() {
    return {
      subscriptionId: -1,
      history: [],
    }
  },
  props: {
    executionId: {
      type: Object,
      required: true,
    },
  },
  methods: {
    load() {
      this.subscriptionId = this.$controller.execution.subscribeExecutionHistory(
        this.subscriptionId,
        {
          executionId: this.executionId,
        },
        this.afterExecutionHistoryReceived,
        this.afterExecutionHistoryError
      );
    },
    afterExecutionHistoryReceived(data, subscriptionId) {
      this.subscriptionId = subscriptionId;
      this.loading = false;
      this.history = data;
    },
    afterExecutionHistoryError(data, subscriptionId) {
      this.subscriptionId = subscriptionId;
      this.loading = false;
      Notify.create({
        message: data,
        color: 'negative',
        position: 'top',
      })
    },
  },
  mounted() {
    this.load();
  }
};
