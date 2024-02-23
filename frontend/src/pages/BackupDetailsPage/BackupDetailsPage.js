import {defineComponent} from 'vue'
import ExecutionHistoryViewer from "components/ExecutionHistoryViewer/ExecutionHistoryViewer.vue";
import moment from "moment";
import {useRouter} from "vue-router";
import ButtonGroup from "components/ButtonGroup/ButtonGroup.vue";
import SlixxButton from "components/SlixxButton/SlixxButton.vue";
import {useGlobalStore} from "stores/global";
import SatelliteSelectableComponent from "components/SatelliteSelectableComponent/SatelliteSelectableComponent";
import JobSelectableComponent from "components/StorageSelectableComponent/StorageSelectableComponent";

export default defineComponent({
  name: 'SatelliteDetailsPage',
  components: {
    JobSelectableComponent,
    SatelliteSelectableComponent,
    SlixxButton,
    ButtonGroup,
    ExecutionHistoryViewer,
  },
  data(){
    return {
      showExecutionHistory: false,
      tab: 'details',
      backupCopy: '',
      subscriptionId: -1,
      backup: {
        name: '',
        id: '',
        executionId: '',
        createdAt: '',
        updatedAt: '',
        configuration: {},
      },
      router: useRouter(),
      globalStore: useGlobalStore(),
      confirmDeleteActive: false,
      confirmDeletionText: ''
    }
  },
  methods: {
    subscribe(callback) {
      const backupId = this.router.currentRoute.value.params.id;
      this.subscriptionId = this.$controller.backup.subscribeBackup(this.subscriptionId, backupId, (data, subscriptionId) => {
        this.subscriptionId = subscriptionId;
        this.backup = data;
        this.backup.createdAt = moment(this.backup.createdAt).format('YYYY-MM-DD HH:mm:ss');
        this.backup.updatedAt = moment(this.backup.updatedAt).format('YYYY-MM-DD HH:mm:ss');

        this.backupCopy = JSON.stringify(this.backup);
        return callback();
      }, (error) => {
        this.$q.notify({
          type: 'negative',
          message: error,
          position: 'top',
        })
        return callback();
      });
    },
    openLogs(callback) {
      this.showExecutionHistory = true;
      return callback();
    },
    restore(callback) {
      const args = {
        backupId: this.backup.id,
        jobId: this.backup.jobId
      }

      this.$controller.backup.restoreBackup(args, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Deletion started',
          position: 'top',
        })
        callback();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data,
          position: 'top',
        })
        callback();
      });
    },
    remove(callback) {
      const args = {
        backupId: this.backup.id,
        jobId: this.backup.jobId
      }

      this.$controller.backup.deleteBackup(args, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Deletion started',
          position: 'top',
        })
        this.router.push({path: `/job/${this.backup.jobId}`})
        callback();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data,
          position: 'top',
        })
        callback();
      });
    },
    confirmDelete(callback) {
      this.confirmDeletionText = '';
      this.confirmDeleteActive = true;
      callback();
    }
  },
  mounted() {
    this.subscribe(() => {});
  }
})
