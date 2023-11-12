<template>
  <q-table
    style="height: 100%"
    title="Storage List"
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
  </q-table>
</template>
<script src="./StorageListComponent.js"/>
