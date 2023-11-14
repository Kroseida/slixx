import {defineStore} from 'pinia';

export const useConstantsStore = defineStore('constants', {
  state: () => ({
    permissionLoaded: false,
    configurationsLoaded: false,
    storageKindsLoaded: false,
    jobStrategiesLoaded: false,
    version: '0.0.1',
    permissions: [],
    storageKinds: [],
    jobStrategies: [],
  }),
  getters: {},
  actions: {
    subscribePermissions(error) {
      this.permissionsSubscriptionId = this.$controller.user.subscribePermissions(
        this.permissionsSubscriptionId,
        (data) => {
          this.permissions = data;
          this.permissionLoaded = true;
        },
        (data) => {
          if (!this.permissions) {
            error(data);
          }
        }
      );
    },
    subscribeStorageKinds(error) {
      this.storageKindsSubscriptionId = this.$controller.storage.subscribeStorageKinds(
        this.storageKindsSubscriptionId,
        (data) => {
          this.storageKinds = data;
          this.storageKindsLoaded = true;
        },
        (data) => {
          if (!this.storageKinds) {
            error(data);
          }
        }
      );
    },
    subscribeJobStrategies(error) {
      this.jobStrategiesSubscriptionId = this.$controller.job.subscribeJobStrategies(
        this.jobStrategiesSubscriptionId,
        (data) => {
          this.jobStrategies = data;
          this.jobStrategiesLoaded = true;
        },
        (data) => {
          if (!this.jobStrategies) {
            error(data);
          }
        }
      );
    }
  },
});
