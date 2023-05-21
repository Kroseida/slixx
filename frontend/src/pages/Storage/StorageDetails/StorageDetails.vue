<template>
  <div>
    <h2 class="page-title">
      {{ $route.params.id === 'new' ? 'Create new storage' : storage ? 'Storage - ' + storage.name : '?' }}
    </h2>
    <b-row>
      <b-col sm="12" md="12">
        <Widget>
          <h4>Details</h4>
          <div class="form-group">
            <b-row>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Name</label>
                <input class="form-control no-border"
                       type="text"
                       name="name"
                       v-model="storage.name"
                       id="storageDetails__name"
                       :readonly="!isPermitted('storage.update') && $route.params.id !== 'new'"/>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Kind</label>
                <select class="form-control no-border"
                        name="kind"
                        id="storageDetails__kind"
                        v-model="storage.kind">
                  <option v-for="kind in Object.keys(kinds)"
                          :key="kind"
                          :value="kind"
                          :disabled="!isPermitted('storage.update') && $route.params.id !== 'new'">
                    {{ kind }}
                  </option>
                </select>
              </b-col>
            </b-row>
          </div>
          <div style="float: right;">
            <b-button type="submit"
                      variant="success"
                      :disabled="$route.params.id !== 'new' && !hasChanges()"
                      id="storageDetails__save__button"
                      @click="save">
              Save Details
            </b-button>
            <b-button type="submit"
                      style="margin-left: 5px;"
                      variant="danger"
                      :disabled="$route.params.id === 'new' || !isPermitted('storage.delete') || $route.params.id === localUser.id"
                      id="storageDetails__delete__button"
                      @click="deleteStorage">
              Delete Storage
            </b-button>
          </div>
        </Widget>
      </b-col>
      <b-col sm="12" md="12"
             v-if="(isPermitted('storage.update') || isPermitted('storage.view'))
              && !!storage.configuration && $route.params.id !== 'new'">
        <Widget>
          <h4>Configuration</h4>
          <div class="form-group">
            <b-row>
              <b-col md="6" class="form-col"
                     :key="configurationDescription[0]"
                     v-for="configurationDescription in Object.entries(kinds[storage.kind] || {})">
                <label>{{ configurationDescription[0] }}</label>
                <ConfigurationKindLong
                    v-if="configurationDescription[1].kind === 'LONG'"
                    :isPermitted="isPermitted('storage.update')"
                    :field="storage.configuration[configurationDescription[0]]"
                    :handle="handleConfigurationChange"
                    :fieldName="configurationDescription[0]"/>
                <ConfigurationKindPassword
                    v-else-if="configurationDescription[1].kind === 'PASSWORD'"
                    :isPermitted="isPermitted('storage.update')"
                    :field="storage.configuration[configurationDescription[0]]"
                    :handle="handleConfigurationChange"
                    :fieldName="configurationDescription[0]"/>
                <ConfigurationKindString
                    v-else
                    :isPermitted="isPermitted('storage.update')"
                    :field="storage.configuration[configurationDescription[0]]"
                    :handle="handleConfigurationChange"
                    :fieldName="configurationDescription[0]"/>
              </b-col>
            </b-row>
          </div>
          <div style="float: right;">
            <b-button type="submit"
                      variant="success"
                      id="storageDetails__save__button"
                      :disabled="!hasConfigurationChanges()"
                      @click="saveConfiguration">
              Save Configuration
            </b-button>
          </div>
        </Widget>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import {mapState} from 'vuex';
import ConfigurationKindLong from "@/components/Configuration/ConfigurationKindLong.vue";
import ConfigurationKindPassword from "@/components/Configuration/ConfigurationKindPassword.vue";
import ConfigurationKindString from "@/components/Configuration/ConfigurationKindString.vue";

export default {
  name: 'StorageDetails',
  components: {
    ConfigurationKindLong,
    ConfigurationKindPassword,
    ConfigurationKindString
  },
  data() {
    return {}
  },
  computed: {
    ...mapState('storage', ['storage', 'originalStorage', 'kinds']),
    ...mapState('layout', ['permissions', 'localUser', 'isPermitted']),
  },
  mounted() {
    this.$store.commit('storage/reset');
    this.load();
  },
  destroyed() {
    this.unload();
    this.$store.commit('storage/reset');
  },
  methods: {
    handleConfigurationChange(fieldName, value) {
      this.storage.configuration[fieldName] = value;
    },
    loadKinds() {
      return new Promise((resolve, reject) => {
        this.$store.commit('storage/subscribeKinds', {
          callback: () => {
            resolve();
          },
          error: (data) => {
            reject(data);
          }
        });
      });
    },
    async load() {
      try {
        await this.loadKinds();
        this.storage.kind = this.storage.kind || Object.keys(this.kinds)[0];
        this.applyConfiguration();
      } catch (e) {
        this.$toasted.error('Error while loading storage kinds: ' + e.message, {
          duration: 5000,
          position: 'top-right',
          fullWidth: true,
          fitToScreen: true,
        });
      }

      if (this.$route.params.id === 'new') {
        return;
      }

      this.$store.commit('storage/subscribeStorage', {
        storageId: this.$route.params.id,
        callback: () => {
          this.applyConfiguration();
        },
        error: (data) => {
          this.$toasted.error('Error while loading storage: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    applyConfiguration() {
      for (let entry of Object.entries(this.kinds[this.storage.kind])) {
        if (this.storage.configuration[entry[0]]) {
          continue;
        }
        if (entry[1].kind === 'LONG') {
          this.storage.configuration[entry[0]] = Number(entry[1].default || 0);
        } else if (entry[1].kind === 'DOUBLE') {
          this.storage.configuration[entry[0]] = Number(entry[1].default || 0);
        } else if (entry[1].kind === 'BOOLEAN') {
          this.storage.configuration[entry[0]] = Boolean(entry[1].default || false);
        } else {
          this.storage.configuration[entry[0]] = entry[1].default || '';
        }
      }
    },
    saveConfiguration() {
      this.$store.commit('storage/saveConfiguration', {
        callback: () => {
          this.$toasted.success('Storage configuration was saved successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while saving storage configuration: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    save() {
      if (this.$route.params.id === 'new') {
        this.$store.commit('storage/createStorage', {
          callback: (storage) => {
            this.$toasted.success('Storage was created successfully', {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
            this.$router.push({name: 'StorageDetails', params: {id: storage.id}});
          },
          error: (data) => {
            this.$toasted.error('Error while creating storage: ' + data.message, {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
          }
        });
        return;
      }
      this.$store.commit('storage/saveStorage', {
        callback: () => {
          this.$toasted.success('Storage was saved successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while saving storage: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    deleteStorage() {
      this.$store.commit('storage/deleteStorage', {
        callback: () => {
          this.$router.push({name: 'StorageList'});
          this.$toasted.success('Storage was deleted successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while deleting storage: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    unload() {
      if (this.$route.params.id === 'new') {
        return;
      }
      this.$store.commit('storage/unsubscribeStorage');
    },
    hasChanges() {
      let storage = JSON.parse(JSON.stringify(this.storage));
      delete storage.configuration;
      let originalStorage = JSON.parse(JSON.stringify(this.originalStorage));
      delete originalStorage.configuration;

      return JSON.stringify(storage) !== JSON.stringify(originalStorage);
    },
    hasConfigurationChanges() {
      return JSON.stringify(this.storage.configuration) !== JSON.stringify(this.originalStorage.configuration);
    },
  }
};
</script>

<style src="./StorageDetails.scss" lang="scss" scoped/>
