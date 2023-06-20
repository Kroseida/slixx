import {defineComponent} from 'vue'
import {useRouter} from 'vue-router'
import {useGlobalStore} from "stores/global";
import moment from "moment/moment";
import ButtonGroup from "components/ButtonGroup/ButtonGroup.vue";
import PermissionAssignComponent from "components/PermissionAssignComponent/PermissionAssignComponent.vue";
import AuthenticationComponent from "components/AuthenticationComponent/AuthenticationComponent.vue";
import SlixxButton from "components/SlixxButton/SlixxButton.vue";

export default defineComponent({
  name: 'UserDetailsPage',
  components: {
    PermissionAssignComponent,
    ButtonGroup,
    AuthenticationComponent,
    SlixxButton,
  },
  data() {
    return {
      activeOptions: [
        {
          label: 'true',
          value: true
        },
        {
          label: 'false',
          value: false
        }
      ],
      user: {
        active: false,
        createdAt: "0000-00-00T00:00:00.0000000+00:00",
        email: "",
        firstName: "",
        id: "00000000-0000-0000-0000-00000000000",
        lastName: "",
        name: "",
        updatedAt: "0000-00-00T00:00:00.0000000+00:00",
        description: "",
      },
      userCopy: JSON.stringify(this.user),
      router: useRouter(),
      globalStore: useGlobalStore(),
      showPermissionView: false,
      showAuthenticationView: false,
    }
  },
  mounted() {
    this.subscribe()
  },
  methods: {
    remove(done) {
      this.$controller.user.deleteUser(this.user, () => {
        this.$q.notify({
          type: 'positive',
          message: 'User was deleted successfully',
          position: 'top',
        })
        this.router.push('/user')
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
      let active = this.user.active;
      if (active != true && active != false) {
        active = active.value;
      }
      let user = {
        name: this.user.name,
        firstName: this.user.firstName,
        lastName: this.user.lastName,
        email: this.user.email,
        active: active,
        description: this.user.description
          .replaceAll("\\", "\\\\")
          .replaceAll('"', '\\"'),
      }
      this.$controller.user.createUser(user, (response) => {
        this.$q.notify({
          type: 'positive',
          message: 'User was saved successfully',
          position: 'top',
        })

        // We are routing back to the user list and then to user to avoid a "bug"(?) where user details page is not updated
        setTimeout(() => {
          this.router.push('/user/' + response.id);
        }, 10);
        this.router.push('/user/')
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
      const userId = this.router.currentRoute.value.params.id;
      if (userId === 'new') {
        this.create(done);
        return;
      }
      let updates = {}
      for (let key in this.user) {
        if (this.user[key] !== JSON.parse(this.userCopy)[key]) {
          if (key === 'permissions') {
            continue;
          }
          updates[key] = this.user[key];

          if (key === 'active') {
            updates[key] = this.user[key].value;
          }
          if (key === 'description') {
            updates[key] = this.user[key].replaceAll("\\", "\\\\").replaceAll('"', '\\"');
          }
        }
      }
      updates.id = this.user.id;

      this.$controller.user.updateUser(updates, () => {
        this.$q.notify({
          type: 'positive',
          message: 'User was saved successfully',
          position: 'top',
        })
        this.subscribe();
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
    subscribe() {
      const userId = this.router.currentRoute.value.params.id;
      if (userId === 'new') {
        return;
      }
      this.subscriptionId = this.$controller.user.subscribeUser(this.subscriptionId, userId, (data) => {
        this.user = data;
        this.user.createdAt = moment(this.user.createdAt).format('YYYY-MM-DD HH:mm:ss');
        this.user.updatedAt = moment(this.user.updatedAt).format('YYYY-MM-DD HH:mm:ss');
        this.user.permissions = JSON.parse(this.user.permissions);

        this.userCopy = JSON.stringify(this.user);
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data
        })
      });
    },
    hasChanges() {
      return this.userCopy !== JSON.stringify(this.user);
    },
    isNewUser() {
      return this.router.currentRoute.value.params.id === 'new';
    },
    isSameUser() {
      return this.user.id === this.globalStore.localUser.id;
    },
    openPermissionView(done) {
      this.showPermissionView = true;
      done();
    },
    openAuthenticationView(done) {
      this.showAuthenticationView = true;
      done();
    },
    addPermission(permission) {
      const args = {
        id: this.user.id,
        permissions: [permission]
      }

      this.$controller.user.addPermission(args, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Permission was added successfully',
          position: 'top',
        });
        this.subscribe();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data,
          position: 'top',
        })
      });
    },
    removePermission(permission) {
      const args = {
        id: this.user.id,
        permissions: [permission]
      }

      this.$controller.user.removePermission(args, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Permission was removed successfully',
          position: 'top',
        });
        this.subscribe();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data,
          position: 'top',
        });
      });
    },
    changePassword(password, done) {
      const args = {
        id: this.user.id,
        password: password
      }
      this.$controller.user.changePassword(args, () => {
        this.$q.notify({
          type: 'positive',
          message: 'Password was changed successfully',
          position: 'top',
        });
        this.subscribe();
        this.showAuthenticationView = false;
        done();
      }, (data) => {
        this.$q.notify({
          type: 'negative',
          message: data,
          position: 'top',
        });
      });
    },
    showSaveButton() {
      if (!this.hasChanges()) {
        return false;
      }
      return true;
    },
    showPermissionButton() {
      if (this.isNewUser()) {
        return false;
      }
      return true
    },
    showAuthButton() {
      if (this.isNewUser()) {
        return false;
      }
      return true;
    },
    showDeleteButton() {
      if (this.isNewUser()) {
        return false;
      }
      if (this.isSameUser()) {
        return false;
      }
      return true;
    }
  }
})
