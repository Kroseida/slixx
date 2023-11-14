import {defineComponent} from 'vue'
import {Notify} from "quasar";
import moment from "moment";

export default defineComponent({
  name: 'StorageList',
  data() {
    return {
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
    columns: {
      type: Array,
      default: () => [
        {
          name: 'id',
          required: true,
          label: 'Id',
          align: 'left',
          field: row => row.id,
          format: val => `${val}`
        },
        {
          name: 'name',
          required: true,
          label: 'Name',
          align: 'left',
          field: row => row.name,
          format: val => `${val}`
        },
        {
          name: 'kind',
          required: true,
          label: 'Kind',
          align: 'left',
          field: row => row.kind,
          format: val => `${val}`
        },
        {
          name: 'createdAt',
          required: true,
          label: 'Created At',
          align: 'left',
          field: row => row.createdAt,
          format: val => `${moment(val).format("YYYY-MM-DD")}`
        },
        {
          name: 'updatedAt',
          required: true,
          label: 'Updated At',
          align: 'left',
          field: row => row.updatedAt,
          format: val => `${moment(val).format("YYYY-MM-DD")}`
        },
      ]
    },
  },
  methods: {
    subscribe(request) {
      this.pagination = request.pagination;
      this.filter = request.filter;

      this.loading = true;
      this.$controller.storage.subscribeStorages(
        this.subscriptionId,
        {
          limit: this.pagination.rowsPerPage,
          search: this.filter,
          page: this.pagination.page,
        },
        this.afterStoragesReceived,
        this.afterStoragesError
      );
    },
    afterStoragesReceived(data) {
      this.loading = false;
      this.rows = data.rows;
      this.pagination.rowsNumber = data.page.totalRows;
    },
    afterStoragesError(data) {
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