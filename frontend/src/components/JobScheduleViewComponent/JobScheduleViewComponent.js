import {defineComponent} from 'vue'
import {Notify} from "quasar";
import JobScheduleCreateComponent from "components/JobScheduleConfigureComponent/JobScheduleConfigureComponent.vue";
import {useGlobalStore} from "stores/global";

export default defineComponent({
  name: 'JobScheduleView',
  components: {
    JobScheduleCreateComponent
  },
  setup() {
    return {
      globalStore: useGlobalStore(),
    }
  },
  data() {
    return {
      createSchedule: false,
      filter: '',
      subscriptionId: -1,
      loading: true,
      pagination: {},
      rows: [],
      newScheduleCopy: {
        name: "",
        description: "",
        kind: "CRON",
        configuration: {}
      },
      updateScheduleOpen: false,
      newSchedule: null,
      updateSchedule: null,
    }
  },
  mounted() {
    this.newSchedule = JSON.parse(JSON.stringify(this.newScheduleCopy))
    this.subscribe({
      pagination: {
        sortBy: null,
        descending: false,
        page: 1,
        rowsPerPage: 9999999,
        rowsNumber: 1
      },
      filter: ""
    });
  },
  props: {
    allowCreation: {
      type: Boolean,
      default: true
    },
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
      data.rows.forEach((row) => {
        row.configuration = JSON.parse(row.configuration)
      })

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
    rowClick(schedule) {
      this.updateScheduleOpen = true
      this.updateSchedule = JSON.parse(JSON.stringify(schedule))
    },
    onCreate(done) {
      const scheduleCreate = JSON.parse(JSON.stringify(this.newSchedule))
      scheduleCreate.jobId = this.jobId
      scheduleCreate.configuration = JSON.stringify(scheduleCreate.configuration)
        .replaceAll("\\", "\\\\")
        .replaceAll('"', '\\"');
      scheduleCreate.description = scheduleCreate.description
        .replaceAll("\\", "\\\\")
        .replaceAll('"', '\\"')

      this.$controller.jobSchedule.createJobSchedule(scheduleCreate, () => {
        done()
        this.subscribe({
          pagination: {
            sortBy: null,
            descending: false,
            page: 1,
            rowsPerPage: 9999999999,
            rowsNumber: 1
          },
          filter: ""
        });
        this.newSchedule = JSON.parse(JSON.stringify(this.newScheduleCopy))
        this.createSchedule = false
        Notify.create({
          message: 'Schedule was created successfully',
          color: 'positiv',
          position: 'top'
        })
      }, (err) => {
        done()
        Notify.create({
          message: err,
          color: 'negative',
          position: 'top',
        })
      })
    },
    onUpdate(done) {
      this.updateSchedule.configuration = JSON.stringify(this.updateSchedule.configuration)
        .replaceAll("\\", "\\\\")
        .replaceAll('"', '\\"');
      this.updateSchedule.description = this.updateSchedule.description
        .replaceAll("\\", "\\\\")
        .replaceAll('"', '\\"')
      delete this.updateSchedule.updatedAt
      delete this.updateSchedule.createdAt

      this.$controller.jobSchedule.updateJobSchedule(this.updateSchedule, () => {
        Notify.create({
          message: 'Schedule was updated successfully',
          color: 'positive',
          position: 'top'
        })
        this.updateScheduleOpen = false
        this.updateSchedule = null
        this.subscribe({
          pagination: {
            sortBy: null,
            descending: false,
            page: 1,
            rowsPerPage: 9999999999,
            rowsNumber: 1
          },
          filter: ""
        });
        done()
      }, (err) => {
        done()
        Notify.create({
          message: err,
          color: 'negative',
          position: 'top',
        })
      })
    },
    onDelete(done) {
      this.$controller.jobSchedule.deleteJobSchedule(this.updateSchedule, () => {
        done()
        this.subscribe({
          pagination: {
            sortBy: null,
            descending: false,
            page: 1,
            rowsPerPage: 9999999999,
            rowsNumber: 1
          },
          filter: ""
        });
        this.newSchedule = JSON.parse(JSON.stringify(this.newScheduleCopy))
        this.createSchedule = false
        this.updateScheduleOpen = false
        this.updateSchedule = null
        Notify.create({
          message: 'Schedule was deleted successfully',
          color: 'positiv',
          position: 'top'
        })
      }, (err) => {
        done()
        Notify.create({
          message: err,
          color: 'negative',
          position: 'top',
        })
      })
    }
  }
})

