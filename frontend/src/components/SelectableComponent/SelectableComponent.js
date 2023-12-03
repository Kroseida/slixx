export default {
  name: "SelectableComponent",
  setup() {

  },
  data() {
    return {
      selected: null,
      show: false,
    };
  },
  props: {
    readonly: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      default: "Select",
    },
    value: {
      type: Object,
      default: () => {
        return {};
      },
    },
    displayField: {
      type: Function,
      default: function (field) {
        return field.name;
      }
    },
  },
  watch: {
    value: {
      handler: function (val) {
        this.show = false;
        this.selected = val;
      }
    }
  }
};
