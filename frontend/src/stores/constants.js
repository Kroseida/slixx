import {defineStore} from 'pinia';

export const useConstantsStore = defineStore('constants', {
  state: () => ({
    permissionLoaded: false,
    configurationsLoaded: false,
    storageKindsLoaded: false,
    jobStrategiesLoaded: false,
    scheduleKindsLoaded: false,
    environmentLoaded: false,
    version: '0.0.1',
    permissions: [],
    storageKinds: [],
    jobStrategies: [],
    scheduleKinds: []
  }),
  getters: {},
  actions: {
    subscribeEnvironment(error) {
      this.environmentSubscriptionId = this.$controller.environment.environment(
        this.environmentSubscriptionId,
        (data) => {
          this.version = data.version;
          this.environmentLoaded = true;
        },
        (data) => {
          if (!this.permissions) {
            error(data);
          }
        }
      );
    },
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
    },
    subscribeJobScheduleKinds(error) {
        this.subscribeJobScheduleKindsId = this.$controller.jobSchedule.subscribeJobScheduleKinds(
          this.subscribeJobScheduleKindsId,
          (data) => {
            this.scheduleKinds = data
            this.scheduleKindsLoaded = true
          },
          (data) => {
            if (!this.scheduleKinds) {
              error(data);
            }
          }
        )
    }
  },
});
