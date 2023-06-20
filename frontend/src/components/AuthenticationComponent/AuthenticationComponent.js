import SlixxButton from "components/SlixxButton/SlixxButton.vue";

export default {
  name: "AuthenticationComponent",
  components: {
    SlixxButton
  },
  setup() {
    return {};
  },
  data() {
    return {
      password: "",
      passwordConfirm: "",
    };
  },
  methods: {
    changePassword(done) {
      this.$emit('change-password', this.password, done);
    }
  }
};
