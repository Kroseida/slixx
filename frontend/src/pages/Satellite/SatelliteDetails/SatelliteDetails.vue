<template>
  <div>
    <h2 class="page-title">
      {{ $route.params.id === 'new' ? 'Create new satellite' : satellite ? 'Satellite - ' + satellite.name : '?' }}
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
                       v-model="satellite.name"
                       id="satelliteDetails__name"
                       :readonly="!isPermitted('satellite.update') && $route.params.id !== 'new'"/>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Address</label>
                <input class="form-control no-border"
                       type="text"
                       name="address"
                       v-model="satellite.address"
                       id="satelliteDetails__address"
                       :readonly="!isPermitted('satellite.update') && $route.params.id !== 'new'"/>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Token</label>
                <input class="form-control no-border"
                       type="password"
                       name="token"
                       v-model="satellite.token"
                       id="satelliteDetails__token"
                       :readonly="!isPermitted('satellite.update') && $route.params.id !== 'new'"/>
              </b-col>
            </b-row>
          </div>
          <div style="float: right;">
            <b-button type="submit"
                      variant="success"
                      :disabled="$route.params.id !== 'new' && !hasChanges()"
                      id="satelliteDetails__save__button"
                      @click="save">
              Save Details
            </b-button>
            <b-button type="submit"
                      style="margin-left: 5px;"
                      variant="danger"
                      :disabled="$route.params.id === 'new' || !isPermitted('satellite.delete') || $route.params.id === localUser.id"
                      id="satelliteDetails__delete__button"
                      @click="deleteSatellite">
              Delete Satellite
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
  name: 'SatelliteDetails',
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
    ...mapState('satellite', ['satellite', 'originalSatellite', 'strategies']),
    ...mapState('layout', ['permissions', 'localUser', 'isPermitted']),
  },
  mounted() {
    this.$store.commit('satellite/reset');
    this.load();
  },
  destroyed() {
    this.unload();
    this.$store.commit('satellite/reset');
  },
  methods: {
    handleConfigurationChange(fieldName, value) {
      this.satellite.configuration[fieldName] = value;
    },
    async load() {
      if (this.$route.params.id === 'new') {
        return;
      }

      this.$store.commit('satellite/subscribeSatellite', {
        satelliteId: this.$route.params.id,
        callback: () => {},
        error: (data) => {
          this.$toasted.error('Error while loading satellite: ' + data.message, {
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
        this.$store.commit('satellite/createSatellite', {
          callback: (satellite) => {
            this.$toasted.success('Satellite was created successfully', {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
            this.$router.push({name: 'SatelliteDetails', params: {id: satellite.id}});
          },
          error: (data) => {
            this.$toasted.error('Error while creating satellite: ' + data.message, {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
          }
        });
        return;
      }
      this.$store.commit('satellite/saveSatellite', {
        callback: () => {
          this.$toasted.success('Satellite was saved successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while saving satellite: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    deleteSatellite() {
      this.$store.commit('satellite/deleteSatellite', {
        callback: () => {
          this.$router.push({name: 'SatelliteList'});
          this.$toasted.success('Satellite was deleted successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while deleting satellite: ' + data.message, {
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
      this.$store.commit('satellite/unsubscribeSatellite');
    },
    hasChanges() {
      let satellite = JSON.parse(JSON.stringify(this.satellite));
      let originalSatellite = JSON.parse(JSON.stringify(this.originalSatellite));

      return JSON.stringify(satellite) !== JSON.stringify(originalSatellite);
    },
    hasConfigurationChanges() {
      return JSON.stringify(this.satellite.configuration) !== JSON.stringify(this.originalSatellite.configuration);
    },
  }
};
</script>

<style src="./SatelliteDetails.scss" lang="scss" scoped/>
