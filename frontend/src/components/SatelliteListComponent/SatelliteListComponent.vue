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
    <template v-slot:top-right v-if="enableFilter && title !== ''">
      <q-input dense debounce="300" v-model="filter" placeholder="Search" style="min-width: 350px">
        <template v-slot:append>
          <q-icon name="search"/>
        </template>
      </q-input>
      <slot name="action"/>
    </template>
    <template v-slot:top v-if="enableFilter && title === ''">
      <q-input dense debounce="300" v-model="filter" placeholder="Search" style="min-width: 100%">
        <template v-slot:append>
          <q-icon name="search"/>
        </template>
      </q-input>
      <slot name="action"/>
    </template>
    <template v-slot:body-cell-connected="props">
      <q-td :props="props" class="text-center" style="text-align: center; width: 80px">
        <q-icon name="favorite" v-if="props.value === 'true'" color="green" style="font-size: 20px;"/>
        <q-icon name="favorite_border" v-else color="red" style="font-size: 22px;"/>
      </q-td>
    </template>
  </q-table>
</template>
<script src="./SatelliteListComponent.js"/>
