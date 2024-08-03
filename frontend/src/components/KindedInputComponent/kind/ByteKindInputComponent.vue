<template>
  <div>
    <q-input
      dense
      filled
      :label="label"
      type="number"
      :readonly="readonly"
      v-model="displayValue"
    >
      <template v-slot:append>
        <q-select
          class="no-hover-focus"
          input-class="no-hover-focus"
          bg-color="transparent"
          v-model="unit"
          :options="units"
          outlined
          filled
          borderless
          dense
        />
      </template>
    </q-input>
  </div>
</template>

<script>
const units = ["Bytes", "KB", "MB", "GB", "TB", "PB"];
const unitTransform = {
  "Bytes": 1,
  "KB": 1024,
  "MB": 1024 * 1024,
  "GB": 1024 * 1024 * 1024,
  "TB": 1024 * 1024 * 1024 * 1024,
  "PB": 1024 * 1024 * 1024 * 1024 * 1024,
};

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
  data() {
    return {
      displayValue: 0,
      unit: "Bytes",
      units,
    };
  },
  mounted() {
    this.$emit("update:modelValue", Number(this.modelValue));
    //this.displayValue = this.modelValue;
    this.selectBestUnit();
    this.prepareUnit();
  },
  methods: {
    prepareUnit() {
      this.displayValue = this.modelValue / unitTransform[this.unit];
    },
    selectBestUnit() {
      const value = this.modelValue;
      let bestUnit = "Bytes";
      for (const unit of units) {
        const transformedValue = value / unitTransform[unit];
        if (transformedValue >= 1) {
          bestUnit = unit;
        }
      }
      this.unit = bestUnit;
    }
  },
  watch: {
    displayValue() {
      this.$emit("update:modelValue", Number(this.displayValue * unitTransform[this.unit]));
    },
    unit() {
      this.prepareUnit();
    }
  }
}
</script>

<style scoped>

</style>
