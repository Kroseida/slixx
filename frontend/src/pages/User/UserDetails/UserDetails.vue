<template>
  <div class="tables-basic">
    <h2 class="page-title">{{
        $route.params.id === 'new' ? 'Create new user' : user ? 'User - ' + user.name : '?'
      }}</h2>
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
                       v-model="user.name"
                       id="userDetails__name"
                       :readonly="!isPermitted('user.update') && $route.params.id !== localUser.id && $route.params.id !== 'new'"/>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>First Name</label>
                <input class="form-control no-border"
                       type="text"
                       name="firstName"
                       v-model="user.firstName"
                       id="userDetails__firstName"
                       :readonly="!isPermitted('user.update') && $route.params.id !== localUser.id && $route.params.id !== 'new'"/>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Last Name</label>
                <input class="form-control no-border"
                       type="text"
                       name="lastName"
                       v-model="user.lastName"
                       id="userDetails__lastName"
                       :readonly="!isPermitted('user.update') && $route.params.id !== localUser.id && $route.params.id !== 'new'"/>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Email</label>
                <input class="form-control no-border"
                       type="email"
                       name="email"
                       v-model="user.email"
                       id="userDetails__email"
                       :readonly="!isPermitted('user.update') && $route.params.id !== localUser.id && $route.params.id !== 'new'"/>
              </b-col>
              <b-col lg="4" md="4" xs="12" class="form-col">
                <label>Status</label>
                <select class="form-control no-border"
                        name="active"
                        id="userDetails__active"
                        v-model="user.active">
                  <option value="true" :disabled="!isPermitted('user.update') && $route.params.id !== 'new'">Active
                  </option>
                  <option value="false" :disabled="!isPermitted('user.update') && $route.params.id !== 'new'">Inactive
                  </option>
                </select>
              </b-col>
              <b-col lg="12" xs="12" md="12" class="form-col">
                <label>Description</label>
                <textarea class="form-control no-border"
                          type="text"
                          name="description"
                          placeholder="Description / Notes"
                          style="min-height: 300px; resize: none;"
                          :readonly="!isPermitted('user.update') && $route.params.id !== 'new'"
                          id="userDetails__description"
                          v-model="user.description"/>
              </b-col>
            </b-row>
          </div>
          <div style="float: right;">
            <b-button type="submit"
                      variant="success"
                      :disabled="$route.params.id !== 'new' && !hasChanges()"
                      id="userDetails__save__button"
                      @click="save">
              Save Details
            </b-button>
            <b-button type="submit"
                      style="margin-left: 5px;"
                      variant="danger"
                      :disabled="$route.params.id === 'new' || !isPermitted('user.delete') || $route.params.id === localUser.id"
                      id="userDetails__delete__button"
                      @click="deleteUser">
              Delete Account
            </b-button>
          </div>
        </Widget>
      </b-col>
      <b-col sm="12" md="5" v-if="$route.params.id !== 'new' && isPermitted('user.permission.update')">
        <Widget style="min-height: 546px;">
          <h4>Permissions</h4>
          <div class="permission-list">
            <div v-for="permission in permissions"
                 :key="permission.value"
                 :id="'permission__' + permission.value.replaceAll('.', '_')"
                 @click="togglePermission(permission.value)">
              <Widget class="entry">
                <b style="float: left">{{ permission.name }}</b>
                <div>
                  <b-badge variant="success"
                           v-if="user.permissions && user.permissions.indexOf(permission.value) !== -1"
                           pill>
                    Permitted
                  </b-badge>
                  <b-badge variant="danger"
                           v-else
                           pill>
                    Forbidden
                  </b-badge>
                </div>
              </Widget>
            </div>
          </div>
        </Widget>
      </b-col>
      <b-col sm="12" md="7"
             v-if="$route.params.id !== 'new' && (isPermitted('user.update') || $route.params.id === localUser.id)">
        <Widget>
          <h4>Authentication</h4>
          <div>
            <b-col lg="12" md="12" xs="12" class="form-col">
              <div class="authentication">
                <div class="title">
                  <label>Password Authentication</label>
                </div>
                <b-row>
                  <b-col md="6" xs="12" class="form-col">
                    <label>Password</label>
                    <input class="form-control no-border"
                           type="password"
                           name="password"
                           id="authentication__password"
                           v-model="authentication.password.value"
                           placeholder="************"/>
                  </b-col>
                  <b-col md="6" xs="12" class="form-col">
                    <label>Repeat Password</label>
                    <input class="form-control no-border"
                           type="password"
                           name="repeatPassword"
                           id="authentication__repeat_password"
                           v-model="authentication.password.repeat"
                           placeholder="************"/>
                  </b-col>
                </b-row>
                <b-button type="submit"
                          variant="success"
                          style="float: right; margin-top: 12px;"
                          id="authentication__save__button"
                          :disabled="!authentication.password.value || authentication.password.value !== authentication.password.repeat"
                          @click="updatePassword">
                  Save Password
                </b-button>
              </div>
            </b-col>
          </div>
        </Widget>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import {mapState, mapActions} from 'vuex';

export default {
  name: 'UserDetails',
  components: {},
  data() {
    return {}
  },
  computed: {
    ...mapState('user', ['user', 'originalUser', 'authentication']),
    ...mapState('layout', ['permissions', 'localUser', 'isPermitted']),
  },
  methods: {
    load() {
      if (this.$route.params.id === 'new') {
        return;
      }
      this.$store.commit('user/subscribeUser', {
        userId: this.$route.params.id,
        error: (data) => {
          this.$toasted.error('Error while loading user: ' + data.message, {
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
      this.$store.commit('user/unsubscribeUser');
    },
    deleteUser() {
      this.$store.commit('user/deleteUser', {
        callback: () => {
          this.$router.push({name: 'UserList'});
          this.$toasted.success('User was deleted successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while deleting user: ' + data.message, {
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
        this.$store.commit('user/createUser', {
          callback: (user) => {
            this.$toasted.success('User was created successfully', {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
            this.$router.push({name: 'UserDetails', params: {id: user.id}});
          },
          error: (data) => {
            this.$toasted.error('Error while creating user: ' + data.message, {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
          }
        });
        return;
      }
      this.$store.commit('user/saveUser', {
        callback: () => {
          this.$toasted.success('User was updated successfully', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while updating user: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    updatePassword() {
      this.$store.commit('user/updatePassword', {
        callback: () => {
          this.authentication.password.value = '';
          this.authentication.password.repeat = '';

          this.$toasted.success('Password was updated successfully.', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        },
        error: (data) => {
          this.$toasted.error('Error while updating password: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
    hasChanges() {
      return JSON.stringify(this.user) !== JSON.stringify(this.originalUser);
    },
    togglePermission(permission) {
      if (this.user.permissions.includes(permission)) {
        this.$store.commit('user/removePermission', {
          permission,
          callback: () => {
          },
          error: (data) => {
            this.$toasted.error('Error while removing permission: ' + data.message, {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
          }
        });
      } else {
        this.$store.commit('user/addPermission', {
          permission,
          callback: () => {
          },
          error: (data) => {
            this.$toasted.error('Error while adding permission: ' + data.message, {
              duration: 5000,
              position: 'top-right',
              fullWidth: true,
              fitToScreen: true,
            });
          }
        });
      }
    }
  },
  mounted() {
    this.$store.commit('user/reset');
    this.load();
  },
  destroyed() {
    this.unload();
    this.$store.commit('user/reset');
  }
};
</script>

<style src="./UserDetails.scss" lang="scss" scoped/>
