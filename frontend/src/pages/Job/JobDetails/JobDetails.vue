<template>
  <div>
    <h2 class="page-title">
      {{ $route.params.id === 'new' ? 'Create new job' : job ? 'Job - ' + job.name : '?' }}
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
                       v-model="job.name"
                       id="jobDetails__name"
                       :readonly="!isPermitted('job.update') && $route.params.id !== 'new'"/>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Strategy</label>
                <select class="form-control no-border"
                        name="strategy"
                        id="jobDetails__strategy"
                        v-model="job.strategy">
                  <option v-for="strategy in Object.keys(strategies)"
                          :key="strategy"
                          :value="strategy"
                          :disabled="!isPermitted('job.update') && $route.params.id !== 'new'">
                    {{ strategy }}
                  </option>
                </select>
              </b-col>
              <b-col lg="4" md="4" xs="0"/>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Origin Storage</label>
                <div v-if="isPermitted('job.update') || $route.params.id === 'new'">
                  <vue-simple-suggest
                      display-attribute="id"
                      valueAttribute="id"
                      v-model="job.originStorageId"
                      :debounce="500"
                      :list="getStorages">
                    <input class="form-control no-border"
                           :value="job.originStorageId || ''"
                           type="text"/>
                    <div slot="suggestion-item" slot-scope="{ suggestion }">
                      <span>{{ suggestion.name }}</span>
                    </div>
                  </vue-simple-suggest>
                </div>
                <div v-else>
                  <input class="form-control no-border"
                         :value="job.originStorageId || ''"
                         type="text"
                         readonly/>
                </div>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Destination Storage</label>
                <div v-if="isPermitted('job.update') || $route.params.id === 'new'">
                  <vue-simple-suggest
                      display-attribute="id"
                      value-attribute="id"
                      v-model="job.destinationStorageId"
                      :debounce="500"
                      :list="getStorages">
                    <input class="form-control no-border"
                           :value="job.destinationStorageId || ''"
                           type="text"/>
                    <div slot="suggestion-item" slot-scope="{ suggestion }">
                      <span>{{ suggestion.name }}</span>
                    </div>
                  </vue-simple-suggest>
                </div>
                <div v-else>
                  <input class="form-control no-border"
                         :value="job.destinationStorageId || ''"
                         type="text"
                         readonly/>
                </div>
              </b-col>
            </b-row>
          </div>
          <div style="float: right;">
            <b-button type="submit"
                      variant="success"
                      :disabled="$route.params.id !== 'new' && !hasChanges()"
                      id="jobDetails__save__button"
                      @click="save">
              Save Details
            </b-button>
            <b-button type="submit"
                      style="margin-left: 5px;"
                      variant="danger"
                      :disabled="$route.params.id === 'new' || !isPermitted('job.delete') || $route.params.id === localUser.id"
                      id="jobDetails__delete__button"
                      @click="deleteJob">
              Delete Job
            </b-button>
          </div>
        </Widget>
      </b-col>
      <b-col sm="12" md="12"
             v-if="(isPermitted('job.update') || isPermitted('job.view'))
              && !!job.configuration && $route.params.id !== 'new'">
        <Widget v-if="job">
          <h4>Configuration</h4>
          <div class="form-group">
            <b-row>
              <b-col md="6" class="form-col"
                     :key="configurationDescription[0]"
                     v-for="configurationDescription in Object.entries(strategies[job.strategy] || {})">
                <label>{{ configurationDescription[0] }}</label>
                <ConfigurationKindLong
                    v-if="configurationDescription[1].kind === 'LONG'"
                    :isPermitted="isPermitted('job.update')"
                    :job="job"
                    :field="job.configuration[configurationDescription[0]]"
                    :handle="handleConfigurationChange"
                    :fieldName="configurationDescription[0]"/>
                <ConfigurationKindPassword
                    v-else-if="configurationDescription[1].kind === 'PASSWORD'"
                    :isPermitted="isPermitted('job.update')"
                    :job="job"
                    :field="job.configuration[configurationDescription[0]]"
                    :handle="handleConfigurationChange"
                    :fieldName="configurationDescription[0]"/>
                <ConfigurationKindString
                    v-else
                    :isPermitted="isPermitted('job.update')"
                    :job="job"
                    :field="job.configuration[configurationDescription[0]]"
                    :handle="handleConfigurationChange"
                    :fieldName="configurationDescription[0]"/>
              </b-col>
            </b-row>
          </div>
          <div style="float: right;">
            <b-button type="submit"
                      variant="success"
                      id="jobDetails__save__button"
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
import VueSimpleSuggest from 'vue-simple-suggest';
import ConfigurationKindLong from "@/components/Configuration/ConfigurationKindLong.vue";
import ConfigurationKindPassword from "@/components/Configuration/ConfigurationKindPassword.vue";
import ConfigurationKindString from "@/components/Configuration/ConfigurationKindString.vue";

export default {
  name: 'JobDetails',
  components: {
    VueSimpleSuggest,
    ConfigurationKindLong,
    ConfigurationKindPassword,
    ConfigurationKindString
  },
  data() {
    return {}
  },
  computed: {
    ...mapState('job', ['job', 'originalJob', 'strategies']),
    ...mapState('layout', ['permissions', 'localUser', 'isPermitted']),
  },
  mounted() {
    this.$store.commit('job/reset');
    this.load();
  },
  destroyed() {
    this.unload();
    this.$store.commit('job/reset');
  },
  methods: {
    async getStorages(filter) {
      return new Promise((resolve) => {
        this.$store.commit('job/subscribeStorages', {
          filter,
          callback: (data) => {
            resolve(data.rows);
          },
          error: (e) => {
            this.$toasted.error('Error while loading storages: ' + e.message, {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
          }
        });
      });
    },
    handleConfigurationChange(fieldName, value) {
      this.job.configuration[fieldName] = value;
    },
    loadStrategies() {
      return new Promise((resolve, reject) => {
        this.$store.commit('job/subscribeStrategies', {
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
        await this.loadStrategies();
        this.job.strategy = this.job.strategy || Object.keys(this.strategies)[0];
        this.applyConfiguration();
      } catch (e) {
        this.$toasted.error('Error while loading job strategies: ' + e.message, {
          duration: 5000,
          position: 'top-right',
          fullWidth: true,
          fitToScreen: true,
        });
      }

      if (this.$route.params.id === 'new') {
        return;
      }

      this.$store.commit('job/subscribeJob', {
        jobId: this.$route.params.id,
        callback: () => {
          this.applyConfiguration();
        },
        error: (data) => {
          this.$toasted.error('Error while loading job: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    applyConfiguration() {
      for (let entry of Object.entries(this.strategies[this.job.strategy])) {
        if (this.job.configuration[entry[0]]) {
          continue;
        }
        if (entry[1].kind === 'LONG') {
          this.job.configuration[entry[0]] = Number(entry[1].default || 0);
        } else if (entry[1].kind === 'DOUBLE') {
          this.job.configuration[entry[0]] = Number(entry[1].default || 0);
        } else if (entry[1].kind === 'BOOLEAN') {
          this.job.configuration[entry[0]] = Boolean(entry[1].default || false);
        } else {
          this.job.configuration[entry[0]] = entry[1].default || '';
        }
      }
    },
    saveConfiguration() {
      this.$store.commit('job/saveConfiguration', {
        callback: () => {
          this.$toasted.success('Job configuration was saved successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while saving job configuration: ' + data.message, {
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
        this.$store.commit('job/createJob', {
          callback: (job) => {
            this.$toasted.success('Job was created successfully', {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
            this.$router.push({name: 'JobDetails', params: {id: job.id}});
          },
          error: (data) => {
            this.$toasted.error('Error while creating job: ' + data.message, {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
          }
        });
        return;
      }
      this.$store.commit('job/saveJob', {
        callback: () => {
          this.$toasted.success('Job was saved successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while saving job: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    deleteJob() {
      this.$store.commit('job/deleteJob', {
        callback: () => {
          this.$router.push({name: 'JobList'});
          this.$toasted.success('Job was deleted successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while deleting job: ' + data.message, {
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
      this.$store.commit('job/unsubscribeJob');
    },
    hasChanges() {
      let job = JSON.parse(JSON.stringify(this.job));
      delete job.configuration;
      let originalJob = JSON.parse(JSON.stringify(this.originalJob));
      delete originalJob.configuration;

      return JSON.stringify(job) !== JSON.stringify(originalJob);
    },
    hasConfigurationChanges() {
      return JSON.stringify(this.job.configuration) !== JSON.stringify(this.originalJob.configuration);
    },
  }
};
</script>

<style src="./JobDetails.scss" lang="scss" scoped/>
