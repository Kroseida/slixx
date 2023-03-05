<template>
  <div class="table-responsive">
    <table class="table table-hover">
      <thead>
      <tr>
        <th v-for="header in headers" :key="header.field">{{ header.name }}</th>
      </tr>
      </thead>
      <tbody v-if="entries.length !== 0">
      <tr v-for="entry in entries"
          :key="entry.id"
          @click="onRowClick ? onRowClick(entry) : () => true"
          :style="onRowClick ? 'cursor: pointer;' : ''">
        <td v-for="header in headers" :key="header.field">
          <slot :name="header.field" :entry="entry">
            {{ entry[header.field] }}
          </slot>
        </td>
      </tr>
      </tbody>
      <tbody v-else>
      <tr>
        <td colspan="100%" class="text-center missing-entries">
          <label>
            No entries found
          </label>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>
<script>
export default {
  name: 'SlixxTable',
  props: {
    headers: {
      type: Array,
      default: () => []
    },
    entries: {
      type: Array,
      default: () => {
      }
    },
    onRowClick: {
      type: Function,
      default: undefined
    },
  }
};
</script>
<style src="./SlixxTable.scss" lang="scss"/>