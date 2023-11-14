import PasswordKindInputComponent from "components/KindedInputComponent/kind/PasswordKindInputComponent.vue";
import TextKindInputComponent from "components/KindedInputComponent/kind/TextKindInputComponent.vue";
import LongKindInputComponent from "components/KindedInputComponent/kind/LongKindInputComponent.vue";

export default {
  name: "KindInputComponent",
  components: {
    PasswordKindInputComponent,
    TextKindInputComponent,
    LongKindInputComponent
  },
  emits: ["update:modelValue"],
  props: {
    label: {
      type: String,
      required: true,
    },
    readonly: {
      type: Boolean,
      required: false,
      default: false,
    },
    modelValue: {
      required: false,
    },
    kind: {
      type: String,
      required: true,
    }
  },
  methods: {
    handleInput(value) {
      this.$emit("update:modelValue", value);
    },
  },
}