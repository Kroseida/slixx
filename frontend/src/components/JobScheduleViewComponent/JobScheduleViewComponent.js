import {defineComponent} from 'vue'
import {Notify} from "quasar";
import JobScheduleCreateComponent from "components/JobScheduleCreateComponent/JobScheduleCreateComponent.vue";

export default defineComponent({
  name: 'JobScheduleView',
  components: {
    JobScheduleCreateComponent
  },
  data() {
    return {
      createSchedule: false,
      filter: '',
      subscriptionId: -1,
      loading: true,
      pagination: {},
      rows: [],
    }
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
      required: true
    }
  },
  unmounted() {
    this.$controller.unsubscribe(this.subscriptionId);
  },
  methods: {
    subscribe(request) {
      this.pagination = request.pagination;
      this.filter = request.filter;

      this.loading = true;
      this.subscriptionId = this.$controller.jobSchedule.subscribeJobSchedules(
        this.subscriptionId,
        {
          limit: this.pagination.rowsPerPage,
          search: this.filter,
          page: this.pagination.page,
          jobId: this.jobId
        },
        this.afterJobSchedulesReceived,
        this.afterJobSchedulesError
      );
    },
    afterJobSchedulesReceived(data, subscriptionId) {
      this.subscriptionId = subscriptionId;
      this.loading = false;
      this.rows = data.rows;
      this.pagination.rowsNumber = data.page.totalRows;
    },
    afterJobSchedulesError(data, subscriptionId) {
      this.subscriptionId = subscriptionId;
      this.loading = false;
      Notify.create({
        message: data,
        color: 'negative',
        position: 'top',
      })
    },
    rowClick(evt, row, index) {
      this.$emit('rowClick', evt, row, index)
    }
  }
})
