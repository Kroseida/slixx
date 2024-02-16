<template>
  <q-card class="right-side-dialog-container">
    <q-item-label header>Schedule</q-item-label>
    <q-separator inset/>
    <div class="slixx-pad-trl-15" v-if="schedule.id">
      <div class="q-gutter-xl">
        <q-input
          v-model="schedule.id"
          readonly
          dense
          filled
          type="text"
          label="ID"
        />
      </div>
    </div>
    <div class="slixx-pad-trl-15">
      <div class="q-gutter-xl">
        <q-input
          v-model="schedule.name"
          dense
          filled
          type="text"
          label="Name"
          :readonly="!globalStore.isPermitted('job.update')"
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
          :readonly="!globalStore.isPermitted('job.update')"
        />
      </div>
    </div>
    <div class="slixx-pad-trl-15">
      <div class="q-gutter-xl">
        <q-input
          dense
          filled
          v-model="schedule.description"
          type="textarea"
          label="Description"
          :readonly="!globalStore.isPermitted('job.update')"
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
            :readonly="!globalStore.isPermitted('job.update')"
          />
        </div>
      </div>
    </div>
    <div class="col-sm-12 col-md-12 slixx-pad-trl-15">
      <button-group style="width: 100%">
        <slixx-button
          color="primary"
          label="Apply"
          class="action"
          @s-click="onApply"
          v-if="globalStore.isPermitted('job.update')"
          :style="showDeleteButton ? 'width: 50%' : 'width: 100%'"
        />
        <slixx-button
          v-if="showDeleteButton && globalStore.isPermitted('job.update')"
          color="negative"
          label="Delete"
          class="action"
          @s-click="onDelete"
          style="width: 50%"
        />
      </button-group>
    </div>
  </q-card>
</template>
<script src="./JobScheduleConfigureComponent.js"/>
