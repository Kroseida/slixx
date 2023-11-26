<template>
  <div class="q-pa-md" style="height: 100%">
    <q-table
      style="height: 100%"
      title="Job List"
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
      <template v-slot:body-cell-strategy="props">
        <q-td style="width: 0">
          <q-chip :label="props.row.strategy" outline color="primary"/>
        </q-td>
      </template>
    </q-table>
  </div>
</template>
<script src="./JobListComponent.js"/>
