import {defineComponent} from 'vue'
import {useRouter} from 'vue-router'
import {useGlobalStore} from 'src/stores/global'
import {useConstantsStore} from "stores/constants";
import SlixxButton from "components/SlixxButton/SlixxButton.vue";
import ButtonGroup from "components/ButtonGroup/ButtonGroup.vue";
import moment from "moment/moment";
import KindInputComponent from "components/KindedInputComponent/KindInputComponent.vue";

export default defineComponent({
  name: 'JobListPage',
  components: {
    SlixxButton,
    ButtonGroup,
    KindInputComponent
  },
  data() {
    const constantsStore = useConstantsStore();

    return {
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
      this.subscriptionId = this.$controller.job.subscribeJob(this.subscriptionId, jobId, (data) => {
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
      }
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
