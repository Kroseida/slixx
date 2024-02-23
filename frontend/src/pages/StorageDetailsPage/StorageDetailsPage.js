import {defineComponent} from 'vue'
import {useRouter} from 'vue-router'
import {useGlobalStore} from 'src/stores/global'
import {useConstantsStore} from "stores/constants";
import SlixxButton from "components/SlixxButton/SlixxButton.vue";
import ButtonGroup from "components/ButtonGroup/ButtonGroup.vue";
import KindInputComponent from "components/KindedInputComponent/KindInputComponent.vue";
import moment from "moment";

export default defineComponent({
  name: 'StorageDetailsPage',
  components: {
    SlixxButton,
    ButtonGroup,
    KindInputComponent
  },
  data() {
    let constantsStore = useConstantsStore();

    return {
      router: useRouter(),
      globalStore: useGlobalStore(),
      constantsStore,
      tab: 'details', // Pre Initialize the tab
      subscriptionId: -1,
      storage: {
        createdAt: "0000-00-00T00:00:00.0000000+00:00",
        id: "00000000-0000-0000-0000-00000000000",
        kind: constantsStore.storageKinds[0].name,
        name: "",
        updatedAt: "0000-00-00T00:00:00.0000000+00:00",
        description: "",
        configuration: {},
      },
      storageCopy: JSON.stringify(this.storage),
      confirmDeleteActive: false,
      confirmDeletionText: '',
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
    storageKindOptions() {
      return this.constantsStore.storageKinds.map((storageKind) => {
        return storageKind.name;
      });
    }
  },
  methods: {
    subscribe(callback) {
      const storageId = this.router.currentRoute.value.params.id;
      if (storageId === 'new') {
        // Initialize the configuration
        this.storageKindByName(this.storage.kind).configuration.forEach((config) => {
          this.storage.configuration[config.name] = config['default'];
        }, 10);
        return callback();
      }
      this.subscriptionId = this.$controller.storage.subscribeStorage(this.subscriptionId, storageId, (data, subscriptionId) => {
        this.subscriptionId = subscriptionId;
        this.storage = data;
        this.storage.createdAt = moment(this.storage.createdAt).format('YYYY-MM-DD HH:mm:ss');
        this.storage.updatedAt = moment(this.storage.updatedAt).format('YYYY-MM-DD HH:mm:ss');
        this.storage.configuration = JSON.parse(this.storage.configuration);

        // Add any missing configuration values
        this.storageKindByName(this.storage.kind).configuration.forEach((config) => {
          if (!this.storage.configuration[config.name]) {
            this.storage.configuration[config.name] = config['default'];
          }
        });

        this.storageCopy = JSON.stringify(this.storage);
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
    remove(done) {
      this.$controller.storage.deleteStorage(this.storage, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Storage was deleted successfully',
          position: 'top',
        })
        this.router.push('/storage')
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
      let storage = {
        name: this.storage.name,
        kind: this.storage.kind,
        configuration: JSON.stringify(this.storage.configuration)
          .replaceAll("\\", "\\\\")
          .replaceAll('"', '\\"'),
        description: this.storage.description
          .replaceAll("\\", "\\\\")
          .replaceAll('"', '\\"'),
      }
      this.$controller.storage.createStorage(storage, (response) => {
        this.$q.notify({
          type: 'positive',
          message: 'Storage was saved successfully',
          position: 'top',
        })

        // We are routing back to the storage list and then to storage to avoid a "bug"(?) where storage details page is not updated
        setTimeout(() => {
          this.router.push('/storage/' + response.id);
        }, 10);
        this.router.push('/storage/')
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
      const storageId = this.router.currentRoute.value.params.id;
      if (storageId === 'new') {
        this.create(done);
        return;
      }
      let updates = {}
      for (let key in this.storage) {
        if (this.storage[key] !== JSON.parse(this.storageCopy)[key]) {
          updates[key] = this.storage[key];

          if (key === 'description') {
            updates[key] = this.storage[key]
              .replaceAll("\\", "\\\\")
              .replaceAll('"', '\\"');
          }
          if (key === 'configuration') {
            updates[key] = JSON.stringify(this.storage[key])
              .replaceAll("\\", "\\\\")
              .replaceAll('"', '\\"');
          }
        }
      }
      updates.id = this.storage.id;

      this.$controller.storage.updateStorage(updates, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Storage was saved successfully',
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
    hasChanges() {
      return this.storageCopy !== JSON.stringify(this.storage);
    },
    isNewStorage() {
      return this.router.currentRoute.value.params.id === 'new';
    },
    storageKindByName(name) {
      return this.constantsStore.storageKinds.find((storageKind) => {
        return storageKind.name === name;
      });
    },
    showSaveButton() {
      if (!this.hasChanges()) {
        return false;
      }
      return true;
    },
    showDeleteButton() {
      if (this.isNewStorage()) {
        return false;
      }
      return true;
    },
    confirmDelete(callback) {
      this.confirmDeletionText = '';
      this.confirmDeleteActive = true;
      callback();
    }
  }
})
