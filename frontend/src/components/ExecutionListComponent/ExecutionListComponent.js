import {defineComponent} from 'vue'
import {Notify} from "quasar";
import moment from "moment";
import ExecutionHistoryViewer from "components/ExecutionHistoryViewer/ExecutionHistoryViewer.vue";

export default defineComponent({
  name: 'ExecutionList',
  data() {
    return {
      showExecutionHistory: false,
      selectedExecution: null,
      filter: '',
      subscriptionId: -1,
      loading: true,
      pagination: {},
      rows: [],
    }
  },
  components: {
    ExecutionHistoryViewer
  },
  mounted() {
    this.subscribe({
      pagination: {
        sortBy: null,
        descending: false,
        page: 1,
        rowsPerPage: 15,
        rowsNumber: 1
      },
      filter: ""
    });
  },
  props: {
    enableFilter: {
      type: Boolean,
      default: true
    },
    jobId: {
      type: String,
      default: null
    },
    columns: {
      type: Array,
      default: () => [
        {
          name: 'createdAt',
          required: true,
          label: 'Created At',
          align: 'left',
          field: row => row.createdAt,
          format: val => `${moment(val).format("YYYY-MM-DD HH:mm:ss")}`
        },
        {
          name: 'finishedAt',
          required: true,
          label: 'Finished At',
          align: 'left',
          field: row => row.finishedAt,
          format: val => `${moment(val).format("YYYY-MM-DD HH:mm:ss")}`
        },
        {
          name: 'jobId',
          required: true,
          label: 'JobId',
          align: 'left',
          field: row => row.jobId,
          format: val => `${val}`
        },
        {
          name: 'status',
          required: true,
          label: 'Status',
          align: 'left',
          field: row => row.status,
          format: val => `${val}`
        },
      ]
    },
  },
  unmounted() {
    this.$controller.unsubscribe(this.subscriptionId)
  },
  methods: {
    subscribe(request) {
      this.pagination = request.pagination;
      this.filter = request.filter;

      const args = {
        limit: this.pagination.rowsPerPage,
        search: this.filter,
        page: this.pagination.page,
        sort: "created_at DESC"
      }
      if (this.jobId) {
        args.jobId = this.jobId;
      }

      this.loading = true;
      this.subscriptionId = this.$controller.execution.subscribeExecutions(
        this.subscriptionId,
        args,
        this.afterExecutionsReceived,
        this.afterExecutionsError
      );
    },
    afterExecutionsReceived(data, subscriptionId) {
      this.subscriptionId = subscriptionId;
      this.loading = false;
      this.rows = data.rows;
      this.pagination.rowsNumber = data.page.totalRows;
    },
    afterExecutionsError(data, subscriptionId) {
      this.subscriptionId = subscriptionId;
      this.loading = false;
      Notify.create({
        message: data,
        color: 'negative',
        position: 'top',
      })
    },
    rowClick(evt, row, index) {
      this.showExecutionHistory = true;
      this.selectedExecution = row.id;

      this.$emit('rowClick', evt, row, index)
    },
    colorOfStatus(status) {
      switch (status) {
        case "FINISHED":
          return "green";
        case "ERROR":
          return "red";
        default:
          return "amber-6";
      }
    }
  }
})
