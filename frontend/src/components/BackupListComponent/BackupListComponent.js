import {defineComponent} from 'vue'
import {Notify} from "quasar";
import moment from "moment";

export default defineComponent({
  name: 'BackupList',
  data() {
    return {
      filter: '',
      subscriptionId: -1,
      loading: true,
      pagination: {},
      rows: [],
    }
  },
  emits: ['rowClick'],
  mounted() {
    this.subscribe({
      pagination: {
        sortBy: null,
        descending: false,
        page: 1,
        rowsPerPage: 15,
        rowsNumber: 1,
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
          name: 'name',
          required: true,
          label: 'Name',
          align: 'left',
          field: row => row.name,
          format: val => `${val}`
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
          name: 'createdAt',
          required: true,
          label: 'Created At',
          align: 'left',
          field: row => row.createdAt,
          format: val => `${moment(val).format("YYYY-MM-DD HH:mm:ss")}`
        },
        {
          name: 'updatedAt',
          required: true,
          label: 'Updated At',
          align: 'left',
          field: row => row.updatedAt,
          format: val => `${moment(val).format("YYYY-MM-DD HH:mm:ss")}`
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
      this.subscriptionId = this.$controller.backup.subscribeBackups(
        this.subscriptionId,
        args,
        this.afterBackupsReceived,
        this.afterBackupsError
      );
    },
    afterBackupsReceived(data, subscriptionId) {
      this.subscriptionId = subscriptionId;
      this.loading = false;
      this.rows = data.rows;
      this.pagination.rowsNumber = data.page.totalRows;
    },
    afterBackupsError(data, subscriptionId) {
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
