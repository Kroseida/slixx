export default {
  name: "SlixxButton",
  data() {
    return {
      loading: false,
    };
  },
  emits: ["s-click"],
  props: {
    color: {
      type: String,
      default: "primary",
    },
    label: {
      type: String,
      default: "",
    },
    icon: {
      type: String,
      default: undefined,
    },
    disable: {
      type: Boolean,
      default: false,
    },
    noCaps: {
      type: Boolean,
      default: false,
    },
    size: {
      type: String,
      default: "md",
    }
  },
  methods: {
    async onClick() {
      this.loading = true;
      this.$emit('s-click', () => {
        this.loading = false;
      });
    }
  }
}
