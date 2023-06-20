import {defineComponent} from 'vue'
import SlixxButton from "components/SlixxButton/SlixxButton.vue";
import ButtonGroup from "components/ButtonGroup/ButtonGroup.vue";
import {useRouter} from 'vue-router'
import moment from "moment";

export default defineComponent({
  name: 'SatelliteDetailsPage',
  components: {
    SlixxButton,
    ButtonGroup
  },
  data() {
    return {
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
    }
  },
  mounted() {
    this.subscribe(() => {});
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
      this.subscriptionId = this.$controller.satellite.subscribeSatellite(this.subscriptionId, satelliteId, (data) => {
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
  }
})
