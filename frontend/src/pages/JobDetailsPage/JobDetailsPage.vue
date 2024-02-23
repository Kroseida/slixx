<template>
  <q-dialog v-model="confirmDeleteActive">
    <q-card>
      <q-toolbar style="padding: 25px; padding-bottom: 15px">
        <q-avatar icon="warning" color="negative" text-color="white" />
        <q-toolbar-title><span class="text-weight-bold">Delete a Backup</span></q-toolbar-title>
        <q-btn flat round dense icon="close" v-close-popup />
      </q-toolbar>

      <q-card-section>
        <p>Please take note: Deleting this job will not remove any data from the storage itself.
          However, the job will no longer be available in the application, and any configurations in the application related to this job will be deleted.
          To confirm, please type the name of the job <b>{{job.name}}</b> and click 'Delete'.</p>
        <q-input v-model="confirmDeletionText" dense label="Confirm" style="margin-top: 15px"/>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat
               label="Delete"
               color="negative"
               v-close-popup
               :disable="confirmDeletionText !== job.name"
               @click="remove"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
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
              label="Execute"
              class="action"
              @s-click="executeBackup"
              v-if="!isNewJob()"
              :disable="hasChanges() || (!globalStore.isPermitted('job.execute'))"
            />
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
              @s-click="confirmDelete"
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
        <q-tab name="schedule" label="Schedule Plan"/>
        <q-tab name="configuration" label="Configuration"/>
        <q-tab name="backups" label="Backups" v-if="!isNewJob() && globalStore.isPermitted('backup.view')"/>
        <q-tab name="executions" label="Executions" v-if="!isNewJob() && globalStore.isPermitted('execution.view')"/>
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
                <job-selectable-component @input="onOriginStorageSelected"
                                          :value="job.originStorageId"
                                          label="Origin Storage"
                                          :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"/>
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <job-selectable-component @input="onDestinationStorageSelected"
                                          :value="job.destinationStorageId"
                                          label="Destination Storage"
                                          :readonly="!globalStore.isPermitted('job.update') || (!globalStore.isPermitted('job.create') && isNewJob())"/>
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <satellite-selectable-component @input="onExecutorSatelliteSelected"
                                                :value="job.executorSatelliteId"
                                                label="Executor Satellite"
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
        <q-tab-panel name="schedule" class="no-padding">
          <job-schedule-view-component :job-id="job.id" :allow-creation="true"/>
        </q-tab-panel>
        <q-tab-panel name="backups" class="no-padding">
          <backup-list-component :job-id="job.id" :columns="backupColumns" @rowClick="openBackup"/>
        </q-tab-panel>
        <q-tab-panel name="executions" class="no-padding">
          <execution-list-component :job-id="job.id" :columns="executionColumns"/>
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
