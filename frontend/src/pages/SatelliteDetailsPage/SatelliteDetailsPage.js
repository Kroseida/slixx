import {defineComponent} from 'vue'
import SlixxButton from "components/SlixxButton/SlixxButton.vue";
import ButtonGroup from "components/ButtonGroup/ButtonGroup.vue";
import {useRouter} from 'vue-router'
import moment from "moment";
import {useGlobalStore} from 'src/stores/global'
import LogViewer from "components/LogViewer/LogViewer.vue";

export default defineComponent({
  name: 'SatelliteDetailsPage',
  components: {
    SlixxButton,
    ButtonGroup,
    LogViewer
  },
  data() {
    return {
      paginationOptions:  [15, 30, 50, 100, 150, 200, -1],
      pagination: {
        sortBy: null,
        descending: false,
        page: 1,
        rowsPerPage: 15,
        rowsNumber: 0
      },
      filter: "",
      router: useRouter(),
      tab: 'details',
      satellite: {
        createdAt: "0000-00-00T00:00:00.0000000+00:00",
        id: "00000000-0000-0000-0000-00000000000",
        lastName: "",
        name: "",
        updatedAt: "0000-00-00T00:00:00.0000000+00:00",
        address: "",
        token: "",
        description: "",
      },
      satelliteCopy: JSON.stringify(this.satellite),
      globalStore: useGlobalStore(),
      subscriptionId: -1,
      logSubscriptionId: -1,
      logs: {
        page: {},
        rows: [],
      },
      confirmDeleteActive: false,
      confirmDeletionText: ''
    }
  },
  mounted() {
    this.subscribe(() => {});
    this.subscribeLogs({
      pagination: this.pagination,
      filter: this.filter,
    });
  },
  unmounted() {
    this.$controller.unsubscribe(this.subscriptionId);
    this.$controller.unsubscribe(this.logSubscriptionId);
  },
  methods: {
    remove(done) {
      this.$controller.satellite.deleteSatellite(this.satellite, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Satellite was deleted successfully',
          position: 'top',
        })
        this.router.push('/satellite')
        done();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data,
          position: 'top',
        })
        done();
      })
    },
    create(done) {
      let satellite = {
        name: this.satellite.name,
        address: this.satellite.address,
        token: this.satellite.token,
        description: this.satellite.description
          .replaceAll("\\", "\\\\")
          .replaceAll('"', '\\"'),
      }
      this.$controller.satellite.createSatellite(satellite, (response) => {
        this.$q.notify({
          type: 'positive',
          message: 'Satellite was saved successfully',
          position: 'top',
        })

        // We are routing back to the satellite list and then to user to avoid a "bug"(?) where user details page is not updated
        setTimeout(() => {
          this.router.push('/satellite/' + response.id);
        }, 10);
        this.router.push('/satellite/')
        done();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data,
          position: 'top',
        })
        done();
      });
    },
    resync(done) {
      const satelliteId = this.router.currentRoute.value.params.id;
      if (satelliteId === 'new') {
        return done();
      }
      this.$controller.satellite.resyncSatellite({id: satelliteId}, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Resync was triggered successfully',
          position: 'top',
        })
        done();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data,
          position: 'top',
        })
        done();
      });
    },
    save(done) {
      const satelliteId = this.router.currentRoute.value.params.id;
      if (satelliteId === 'new') {
        this.create(done);
        return;
      }
      let updates = {}
      for (let key in this.satellite) {
        if (this.satellite[key] !== JSON.parse(this.satelliteCopy)[key]) {
          updates[key] = this.satellite[key];

          if (key === 'description') {
            updates[key] = this.satellite[key]
              .replaceAll("\\", "\\\\")
              .replaceAll('"', '\\"');
          }
        }

        updates.id = this.satellite.id;

        this.$controller.satellite.updateSatellite(updates, () => {
          this.$q.notify({
            type: 'positive',
            message: 'Satellite was saved successfully',
            position: 'top',
          })
          this.subscribe(() => {
          });
          done();
        }, (data) => {
          this.$q.notify({
            type: 'negative',
            message: data,
            position: 'top',
          })
          done();
        });
      }
    },
    subscribe(callback) {
      const satelliteId = this.router.currentRoute.value.params.id;
      if (satelliteId === 'new') {
        return callback();
      }
      this.subscriptionId = this.$controller.satellite.subscribeSatellite(this.subscriptionId, satelliteId, (data, subscribeId) => {
        this.subscriptionId = subscribeId;
        if (this.satellite.connected !== data.connected && this.satellite.connected !== undefined) {
          if (data.connected) {
            this.$q.notify({
              type: 'positive',
              message: 'Satellite connection was established',
              position: 'top',
            })
          } else {
            this.$q.notify({
              type: 'negative',
              message: 'Lost connection to satellite',
              position: 'top',
            })
          }
        }

        this.satellite = data;
        this.satellite.createdAt = moment(this.satellite.createdAt).format('YYYY-MM-DD HH:mm:ss');
        this.satellite.updatedAt = moment(this.satellite.updatedAt).format('YYYY-MM-DD HH:mm:ss');

        this.satelliteCopy = JSON.stringify(this.satellite);
        return callback();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data
        })
        return callback();
      });
    },
    subscribeLogs(request) {
      if (!request) {
        request = {
          pagination: this.pagination,
          filter: this.filter,
        }
      }
      const satelliteId = this.router.currentRoute.value.params.id;
      if (satelliteId === 'new') {
        return;
      }
      this.pagination = request.pagination;
      this.filter = request.filter;

      this.loading = true;
      this.logSubscriptionId = this.$controller.satellite.subscribeSatelliteLogs(
        this.logSubscriptionId,
        {
          id: satelliteId,
          limit: this.pagination.rowsPerPage,
          search: this.filter,
          page: this.pagination.page,
        },
        (data, subscribeId) => {
          this.logSubscriptionId = subscribeId;
          this.loading = false;
          this.logs = data;
          this.pagination.rowsNumber = data.page.totalRows;
        },
        (data) => {
          this.loading = false;
          this.$q.notify({
            type: 'negative',
            message: data
          })
        }
      );
    },
    hasChanges() {
      return this.satelliteCopy !== JSON.stringify(this.satellite);
    },
    showSaveButton() {
      if (!this.hasChanges()) {
        return false;
      }
      return true;
    },
    isNewSatellite() {
      return this.router.currentRoute.value.params.id === 'new';
    },
    showDeleteButton() {
      if (this.isNewSatellite()) {
        return false;
      }
      return true;
    },
    changePage(page) {
      this.pagination.page = page;
      this.subscribeLogs();
    },
    changeRowsPerPage(rowsPerPage) {
      this.pagination.rowsPerPage = rowsPerPage;
      this.subscribeLogs();
    },
    confirmDelete(callback) {
      this.confirmDeletionText = '';
      this.confirmDeleteActive = true;
      callback();
    }
  }
})
