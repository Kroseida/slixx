<template>
  <div class="auth-page">
    <b-container>
      <h5 class="auth-logo">
        <i class="fa fa-circle text-primary"></i>
        Slixx Dashboard
        <i class="fa fa-circle text-danger"></i>
      </h5>
      <Widget class="widget-auth mx-auto" title="<h3 class='mt-0'>Management Login</h3>" customHeader>
        <form class="mt" @submit.prevent="login" style="padding-top: 15px">
          <div class="form-group">
            <input class="form-control no-border"
                   type="text"
                   name="name"
                   placeholder="Name"
                   v-model="authentication.name"
                   id="authentication__name"
                   required/>
          </div>
          <div class="form-group">
            <input class="form-control no-border"
                   type="password"
                   name="password"
                   placeholder="Password"
                   v-model="authentication.password"
                   id="authentication__password"
                   required/>
          </div>
          <div style="padding: 10px 22px 0px;">
            <b-button
                type="submit"
                size="sm"
                class="auth-btn mb-3"
                variant="inverse"
                v-if="!loggingIn"
                id="authentication__button_execute"
                @click="login">
              Login
            </b-button>
            <b-button
                type="submit"
                size="sm"
                class="auth-btn mb-3"
                variant="inverse"
                id="authentication__button_execute"
                v-else
                disabled>
              <i class="fa fa-circle-o-notch fa-spin"></i>
            </b-button>
          </div>
        </form>
      </Widget>
    </b-container>
    <footer class="auth-footer">
      Slixx Backup Solution © 2023 (v0.0.0) (based on Sing App Vue Dashboard <a
        href="https://flatlogic.com/">Flatlogic</a>)
    </footer>
  </div>
</template>

<script>
import Widget from '@/components/Widget/Widget';
import {mapState} from "vuex";

export default {
  name: 'LoginPage',
  components: {Widget},
  computed: {
    ...mapState('login', ['authentication', 'loggingIn']),
    ...mapState('layout', ['localUser'])
  },
  methods: {
    login() {
      this.$store.commit('login/login', {
        callback: (data) => {
          this.$toasted.success('Logged in successfully!', {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
          this.$store.commit('layout/afterAuthentication', {token: data.message[0].authentication.token});
        },
        error: (data) => {
          this.$toasted.error('Error while logging in: ' + data.message, {
            duration: 5000,
            position: 'top-right',
            fullWidth: true,
            fitToScreen: true,
          });
        }
      });
    },
  },
  created() {
    if (this.localUser) {
      this.$router.push('/app/dashboard');
    }
  },
};
</script>
