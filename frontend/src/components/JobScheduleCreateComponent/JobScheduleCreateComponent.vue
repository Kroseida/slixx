<template>
  <q-card class="right-side-dialog-container">
    <q-item-label header>Schedule</q-item-label>
    <q-separator inset/>
    <div class="slixx-pad-trl-15">
      <div class="q-gutter-xl">
        <q-input
          dense
          filled
          type="text"
          label="Name"
        />
      </div>
    </div>
    <div class="slixx-pad-trl-15">
      <div class="q-gutter-xl">
        <q-select
          dense
          filled
          v-model="schedule.kind"
          :options="kinds"
          label="Kind"
        />
      </div>
    </div>
    <div class="slixx-pad-trl-15">
      <div class="q-gutter-xl">
        <q-input
          dense
          filled
          type="textarea"
          label="Description"
        />
      </div>
    </div>
    <q-item-label header>Configuration</q-item-label>
    <q-separator inset/>
    <div class="row">
      <div class="col-sm-12 col-md-12 slixx-pad-trl-15"
           v-for="configurationEntry in scheduleKindByName(schedule.kind).configuration"
           v-bind:key="configurationEntry.name">
        <div class="q-gutter-xl">
          <kind-input-component
            :kind="configurationEntry.kind"
            v-model="schedule.configuration[configurationEntry.name]"
            :label="configurationEntry.name"
          />
        </div>
      </div>
    </div>
  </q-card>
</template>
<script src="./JobScheduleCreateComponent.js"/>
