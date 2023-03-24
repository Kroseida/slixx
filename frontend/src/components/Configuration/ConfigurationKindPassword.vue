<template>
  <div>
    <input class="form-control no-border"
           type="password"
           v-model="value"
           :id="'configuration__' + fieldName"
           :readonly="!isPermitted"/>
  </div>
</template>

<script>
export default {
  name: 'ConfigurationKindLong',
  components: {},
  watch: {
    field() {
      this.$set(this, 'value', this.field);
    },
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