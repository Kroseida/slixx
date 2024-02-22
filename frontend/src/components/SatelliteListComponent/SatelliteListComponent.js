import {defineComponent} from 'vue'
import {Notify} from "quasar";
import moment from "moment";

export default defineComponent({
  name: 'SatelliteList',
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
        rowsNumber: 1
      },
      filter: ""
    });
  },
  props: {
    title: {
      type: String,
      default: 'Satellite List'
    },
    enableFilter: {
      type: Boolean,
      default: true
    },
    columns: {
      type: Array,
      default: () => [
        {
          name: 'connected',
          required: true,
          label: '',
          align: 'left',
          field: row => row.connected,
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
          name: 'address',
          required: true,
          label: 'Address',
          align: 'left',
          field: row => row.address,
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
  unmounted() {
    this.$controller.unsubscribe(this.subscriptionId);
  },
  methods: {
    subscribe(request) {
      this.pagination = request.pagination;
      this.filter = request.filter;

      this.loading = true;
      this.subscriptionId = this.$controller.satellite.subscribeSatellites(
        this.subscriptionId,
        {
          limit: this.pagination.rowsPerPage,
          search: this.filter,
          page: this.pagination.page,
        },
        this.afterSatellitesReceived,
        this.afterSatellitesError
      );
    },
    afterSatellitesReceived(data, subscriptionId) {
      this.subscriptionId = subscriptionId;
      this.loading = false;
      this.rows = data.rows;
      this.pagination.rowsNumber = data.page.totalRows;
    },
    afterSatellitesError(data, subscriptionId) {
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
