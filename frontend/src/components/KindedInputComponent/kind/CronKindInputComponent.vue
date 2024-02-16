<template>
  <div>
    <q-input
      dense
      filled
      :label="label"
      type="text"
      :readonly="readonly"
      :model-value="modelValue"
      @update:model-value="handleInput"
      :hint="currentDescription"
    />
  </div>
</template>

<script>
import cronstrue from 'cronstrue';

export default {
  name: "CronKindInputComponent",
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
  methods: {
    handleInput(value) {
      this.$emit("update:modelValue", value);
    },
  },
  computed: {
    currentDescription() {
      try {
        return cronstrue.toString(this.modelValue)
      } catch (e) {
        return e
      }
    }
  }
}
</script>

<style scoped>

</style>
