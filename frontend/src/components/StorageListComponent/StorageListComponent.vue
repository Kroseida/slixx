<template>
  <q-table
    style="height: 100%"
    :title="title"
    :rows="rows"
    v-model:pagination="pagination"
    :loading="loading"
    :filter="filter"
    :columns="columns"
    @row-click="rowClick"
    row-key="name"
    @request="subscribe"
  >
    <template v-slot:top-right v-if="enableFilter">
      <q-input dense debounce="300" v-model="filter" placeholder="Search" style="min-width: 350px">
        <template v-slot:append>
          <q-icon name="search"/>
        </template>
      </q-input>
      <slot name="action"/>
    </template>
    <template v-slot:body-cell-kind="props">
      <q-td style="width: 0">
        <q-chip :label="props.row.kind" outline color="primary"/>
      </q-td>
    </template>
  </q-table>
</template>
<script src="./StorageListComponent.js"/>
