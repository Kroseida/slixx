<template>
  <div>
    <input class="form-control no-border"
           type="text"
           v-model="value"
           :id="'storageConfiguration__' + fieldName"
           :readonly="!isPermitted"/>
  </div>
</template>

<script>
export default {
  name: 'StorageConfigurationKindLong',
  components: {},
  watch: {
    value() {
      this.handle(this.fieldName, this.value);
    },
    storage() {
      this.$set(this, 'value', this.field);
    }
  },
  props: {
    isPermitted: {
      type: Boolean,
      required: true
    },
    storage: {
      type: Object,
      required: true
    },
    fieldName: {
      type: String,
      required: true
    },
    field: {
      required: true
    },
    handle: {
      type: Function,
      required: true
    }
  },
  mounted() {
    this.$set(this, 'value', this.field);
  },
  data() {
    return {
      value: ''
    }
  },
};
</script>