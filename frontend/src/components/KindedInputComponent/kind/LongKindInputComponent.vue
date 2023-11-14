<template>
  <div>
    <q-input
      dense
      filled
      :label="label"
      type="number"
      :readonly="readonly"
      :model-value="modelValue"
      @update:model-value="handleInput"
    />
  </div>
</template>

<script>
export default {
  name: "TextKindInputComponent",
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
      required: true,
    },
  },
  mounted() {
    this.$emit("update:modelValue", Number(this.modelValue));
  },
  methods: {
    handleInput(value) {
      this.$emit("update:modelValue", Number(value));
    },
  },
}
</script>

<style scoped>

</style>
