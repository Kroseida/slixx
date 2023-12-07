<template>
  <q-dialog v-model="showExecutionHistory" persistent>
    <q-card style="width: 100%; max-width: 800px; min-width: 600px; min-height: 400px; max-height: 600px; overflow: auto;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Execution Logs</div>
        <q-space />
        <q-btn icon="close" flat round dense @click="showExecutionHistory = false" />
      </q-card-section>

      <q-card-section class="row items-center q-pt-none scroll" style="overflow-y: scroll; max-height: 500px;">
        <execution-history-viewer :execution-id="selectedExecution"/>
      </q-card-section>
    </q-card>
  </q-dialog>
  <q-table
    style="height: 100%"
    title="Execution List"
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
    <template v-slot:body-cell-status="props">
      <q-td>
        <q-chip :label="props.row.status"
                outline
                :color="colorOfStatus(props.row.status)"/>
      </q-td>
    </template>
  </q-table>
</template>
<script src="./ExecutionListComponent.js"/>
