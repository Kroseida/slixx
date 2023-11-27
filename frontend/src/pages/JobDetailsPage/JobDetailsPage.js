import {defineComponent} from 'vue'
import {useRouter} from 'vue-router'
import {useGlobalStore} from 'src/stores/global'
import {useConstantsStore} from "stores/constants";
import SlixxButton from "components/SlixxButton/SlixxButton.vue";
import ButtonGroup from "components/ButtonGroup/ButtonGroup.vue";
import moment from "moment/moment";
import KindInputComponent from "components/KindedInputComponent/KindInputComponent.vue";
import BackupListComponent from "components/BackupListComponent/BackupListComponent.vue";
import ExecutionListComponent from "components/ExecutionListComponent/ExecutionListComponent.vue";

export default defineComponent({
  name: 'JobListPage',
  components: {
    SlixxButton,
    ButtonGroup,
    KindInputComponent,
    BackupListComponent,
    ExecutionListComponent,
  },
  data() {
    const constantsStore = useConstantsStore();

    return {
      backupColumns: [
        {
          name: 'name',
          required: true,
          label: 'Name',
          align: 'left',
          field: row => row.name,
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
      ],
      executionColumns: [
        {
          name: 'createdAt',
          required: true,
          label: 'Started At',
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
          format: val => val === null ? '' : `${moment(val).format("YYYY-MM-DD HH:mm:ss")}`
        },
        {
          name: 'duration',
          required: true,
          label: 'Duration',
          align: 'left',
          field: row => row.finishedAt,
          format: (val, row) => {
            if (val === null) {
              return '';
            }
            const duration = moment.duration(moment(row.finishedAt).diff(moment(row.createdAt)));
            let time = '';
            if (duration.hours() > 0) {
              time += `${duration.hours()}h `;
            }
            if (duration.minutes() > 0) {
              time += `${duration.minutes()}m `;
            }
            if (duration.seconds() > 0) {
              time += `${duration.seconds()}s `;
            }
            if (duration.seconds() === 0 && duration.minutes() === 0 && duration.hours() === 0) {
              time += `${duration.milliseconds()}ms`;
            }
            return time;
          }
        },
        {
          name: 'status',
          required: true,
          label: 'Status',
          align: 'left',
          field: row => row.status,
          format: val => `${val}`
        },
      ],
      tab: 'details',
      job: {
        createdAt: "0000-00-00T00:00:00.0000000+00:00",
        id: "00000000-0000-0000-0000-00000000000",
        strategy: constantsStore.jobStrategies[0].name,
        name: "",
        updatedAt: "0000-00-00T00:00:00.0000000+00:00",
        destinationStorageId: "",
        originStorageId: "",
        description: "",
        configuration: {},
        executorSatelliteId: "",
      },
      jobCopy: JSON.stringify(this.storage),
      router: useRouter(),
      globalStore: useGlobalStore(),
      constantsStore,
    }
  },
  mounted() {
    this.subscribe(() => {
      this.tab = 'configuration';
      setTimeout(() => {
        this.tab = 'details';
      });
    })
  },
  unmounted() {
    this.$controller.unsubscribe(this.subscriptionId);
  },
  computed: {
    jobStrategyOptions() {
      return this.constantsStore.jobStrategies.map((strategy) => {
        return strategy.name;
      });
    }
  },
  methods: {
    subscribe(callback) {
      const jobId = this.router.currentRoute.value.params.id;
      if (jobId === 'new') {
        // Initialize the configuration
        this.jobStrategyByName(this.job.strategy).configuration.forEach((config) => {
          this.job.configuration[config.name] = config['default'];
        }, 10);
        return callback();
      }
      this.subscriptionId = this.$controller.job.subscribeJob(this.subscriptionId, jobId, (data, subscriptionId) => {
        this.subscriptionId = subscriptionId;
        this.job = data;
        this.job.createdAt = moment(this.job.createdAt).format('YYYY-MM-DD HH:mm:ss');
        this.job.updatedAt = moment(this.job.updatedAt).format('YYYY-MM-DD HH:mm:ss');
        this.job.configuration = JSON.parse(this.job.configuration);

        // Add any missing configuration values
        this.jobStrategyByName(this.job.strategy).configuration.forEach((config) => {
          if (!this.job.configuration[config.name]) {
            this.job.configuration[config.name] = config['default'];
          }
        });

        this.jobCopy = JSON.stringify(this.job);
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
    save(done) {
      const jobId = this.router.currentRoute.value.params.id;
      if (jobId === 'new') {
        this.create(done);
        return;
      }
      let updates = {}
      for (let key in this.job) {
        if (this.job[key] !== JSON.parse(this.jobCopy)[key]) {
          updates[key] = this.job[key];

          if (key === 'description') {
            updates[key] = this.job[key]
              .replaceAll("\\", "\\\\")
              .replaceAll('"', '\\"');
          }
          if (key === 'configuration') {
            updates[key] = JSON.stringify(this.job[key])
              .replaceAll("\\", "\\\\")
              .replaceAll('"', '\\"');
          }
        }
      }
      updates.id = this.job.id;

      this.$controller.job.updateJob(updates, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Job was saved successfully',
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
    },
    executeBackup(done) {
      const job = {
        jobId: this.job.id,
      }
      this.$controller.job.executeBackup(job, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Backup was triggered successfully',
          position: 'top',
        })
        done();
      }, (err) => {
        this.$q.notify({
          type: 'negative',
          message: err,
          position: 'top',
        })
        done();
      });
    },
    create(done) {
      let job = {
        name: this.job.name,
        strategy: this.job.strategy,
        configuration: JSON.stringify(this.job.configuration)
          .replaceAll("\\", "\\\\")
          .replaceAll('"', '\\"'),
        description: this.job.description
          .replaceAll("\\", "\\\\")
          .replaceAll('"', '\\"'),
        originStorageId: this.job.originStorageId,
        destinationStorageId: this.job.destinationStorageId,
        executorSatelliteId: this.job.executorSatelliteId,
      }
      this.$controller.job.createJob(job, (response) => {
        this.$q.notify({
          type: 'positive',
          message: 'Job was saved successfully',
          position: 'top',
        })

        // We are routing back to the storage list and then to storage to avoid a "bug"(?) where storage details page is not updated
        setTimeout(() => {
          this.router.push('/job/' + response.id);
        }, 10);
        this.router.push('/job/')
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
    remove(done) {
      this.$controller.job.deleteJob(this.job, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Job was deleted successfully',
          position: 'top',
        })
        this.router.push('/job')
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
    hasChanges() {
      return this.jobCopy !== JSON.stringify(this.job);
    },
    jobStrategyByName(name) {
      return this.constantsStore.jobStrategies.find((strategy) => {
        return strategy.name === name;
      });
    },
    isNewJob() {
      return this.router.currentRoute.value.params.id === 'new';
    },
    showSaveButton() {
      if (!this.hasChanges()) {
        return false;
      }
      return true;
    },
    showDeleteButton() {
      if (this.isNewJob()) {
        return false;
      }
      return true;
    }
  }
})
