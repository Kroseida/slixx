<template>
  <q-dialog v-model="confirmDeleteActive">
    <q-card>
      <q-toolbar style="padding: 25px; padding-bottom: 15px">
        <q-avatar icon="warning" color="negative" text-color="white" />
        <q-toolbar-title><span class="text-weight-bold">Delete a Backup</span></q-toolbar-title>
        <q-btn flat round dense icon="close" v-close-popup />
      </q-toolbar>

      <q-card-section>
        <p>Please take note: Deleting this backup will permanently remove all associated data, configurations, dates, and logs.
          The backup will be irreversibly deleted and cannot be restored.
          To confirm, kindly type the name of the backup <b>{{backup.name}}</b> and click 'Delete'.</p>

        <q-input v-model="confirmDeletionText" dense label="Confirm" style="margin-top: 15px"/>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat
               label="Delete"
               color="negative"
               v-close-popup
               :disable="confirmDeletionText !== backup.name"
               @click="remove"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
  <q-dialog v-model="showExecutionHistory" persistent>
    <q-card
      style="width: 100%; max-width: 800px; min-width: 600px; min-height: 400px; max-height: 600px; overflow: auto;">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">Execution Logs</div>
        <q-space/>
        <q-btn icon="close" flat round dense @click="showExecutionHistory = false"/>
      </q-card-section>

      <q-card-section class="row items-center q-pt-none scroll" style="overflow-y: scroll; max-height: 500px;">
        <execution-history-viewer :execution-id="backup.executionId"/>
      </q-card-section>
    </q-card>
  </q-dialog>
  <div class="q-pa-md">
    <q-card>
      <q-card-section>
        <div class="relative-position row items-center">
          <div class="q-table__title">
            Backup
          </div>
          <div class="col"/>
          <button-group>
            <slixx-button
              color="primary"
              label="Restore"
              class="action"
              @s-click="restore"
              :disable="!globalStore.isPermitted('backup.restore')"
            />
            <slixx-button
              color="primary"
              label="Logs"
              class="action"
              @s-click="openLogs"
              :disable="!globalStore.isPermitted('backup.view')"
            />
            <slixx-button
              color="negative"
              label="Delete"
              class="action"
              @s-click="confirmDelete"
              :disable="!globalStore.isPermitted('backup.delete')"
            />
          </button-group>
        </div>
      </q-card-section>
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
                  v-model="backup.id"
                  readonly
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Execution ID"
                  v-model="backup.executionId"
                  readonly
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Created at"
                  v-model="backup.createdAt"
                  readonly
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Updated at"
                  v-model="backup.updatedAt"
                  readonly
                />
              </div>
            </div>
          </div>
        </q-tab-panel>
      </q-tab-panels>

    </q-card>
  </div>
</template>
<style lang="scss">
@import "./BackupDetailsPage.scss";
</style>
<script src="./BackupDetailsPage.js"/>
