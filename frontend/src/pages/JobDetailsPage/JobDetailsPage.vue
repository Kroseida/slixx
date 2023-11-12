<template>
  <div class="q-pa-md">
    <q-card>
      <q-card-section>
        <div class="relative-position row items-center">
          <div class="q-table__title">
            Job
          </div>
          <div class="col"/>
          <button-group>
            <slixx-button
              color="primary"
              label="Save"
              class="action"
              @s-click="save"
              :disable="!showSaveButton() || (!globalStore.isPermitted('job.update') && !globalStore.isPermitted('job.create'))"
            />
            <slixx-button
              color="negative"
              label="Delete"
              class="action"
              @s-click="remove"
              :disable="!showDeleteButton() || (!globalStore.isPermitted('job.delete'))"
            />
          </button-group>
        </div>
      </q-card-section>
      <q-separator/>
      <q-tabs
        v-model="tab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="justify"
        narrow-indicator
      >
        <q-tab name="details" label="Details"/>
        <q-tab name="configuration" label="Configuration"/>
        <q-tab name="backups" label="Backups" v-if="!isNewJob()"/>
      </q-tabs>
      <q-separator/>
      <q-tab-panels v-model="tab" animated>
        <q-tab-panel name="details">
          <div class="row">
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="ID"
                  v-model="job.id"
                  readonly
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Name"
                  v-model="job.name"
                  :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-select
                  dense
                  filled
                  v-model="job.strategy"
                  :options="jobStrategyOptions"
                  label="Strategy"
                  :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Origin Storage"
                  v-model="job.originStorageId"
                  :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Destination Storage"
                  v-model="job.destinationStorageId"
                  :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Executor Satellite"
                  v-model="job.executorSatelliteId"
                  :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-12 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  v-model="job.description"
                  label="Description"
                  :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"
                  type="textarea"
                />
              </div>
            </div>
          </div>
        </q-tab-panel>
        <q-tab-panel name="configuration">
          <div class="row">
            <div class="col-sm-12 col-md-4 slixx-pad-5"
                 v-for="configurationEntry in jobStrategyByName(job.strategy).configuration"
                 v-bind:key="configurationEntry.name">
              <div class="q-gutter-xl">
                <kind-input-component
                  :kind="configurationEntry.kind"
                  v-model="job.configuration[configurationEntry.name]"
                  :label="configurationEntry.name"
                  :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"
                />
              </div>
            </div>
          </div>
        </q-tab-panel>
        <q-tab-panel name="backups" class="no-padding">
          <backup-list-component :job-id="job.id"/>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </div>
</template>
<style lang="scss">
@import "./JobDetailsPage.scss";

.no-padding {
  padding: 0;
}

</style>
<script src="./JobDetailsPage.js"/>
