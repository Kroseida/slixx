<template>
  <q-dialog v-model="confirmDeleteActive">
    <q-card>
      <q-toolbar style="padding: 25px; padding-bottom: 15px">
        <q-avatar icon="warning" color="negative" text-color="white" />
        <q-toolbar-title><span class="text-weight-bold">Delete a Storage</span></q-toolbar-title>
        <q-btn flat round dense icon="close" v-close-popup />
      </q-toolbar>

      <q-card-section>
        <p>Please note: Deleting this storage will not remove any data from the storage itself.
        However, the storage will no longer be available in the application, and any configurations in the application related to this storage will be deleted.
        To confirm, please type the name of the storage <b>{{storage.name}}</b> and click 'Delete'.</p>

        <p>If this storage is currently in use, deletion is not possible.</p>
        <q-input v-model="confirmDeletionText" dense label="Confirm" style="margin-top: 15px"/>
      </q-card-section>
      <q-card-actions align="right">
        <q-btn flat
               label="Delete"
               color="negative"
               v-close-popup
               :disable="confirmDeletionText !== storage.name"
               @click="remove"/>
      </q-card-actions>
    </q-card>
  </q-dialog>
  <div class="q-pa-md">
    <q-card>
      <q-card-section>
        <div class="relative-position row items-center">
          <div class="q-table__title">
            Storage
          </div>
          <div class="col"/>
          <button-group>
            <slixx-button
              color="primary"
              label="Save"
              class="action"
              @s-click="save"
              :disable="!showSaveButton() || (!globalStore.isPermitted('storage.update') && !globalStore.isPermitted('storage.create'))"
            />
            <slixx-button
              color="negative"
              label="Delete"
              class="action"
              @s-click="confirmDelete"
              :disable="!showDeleteButton() || (!globalStore.isPermitted('storage.delete'))"
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
                  v-model="storage.id"
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
                  v-model="storage.name"
                  :readonly="!globalStore.isPermitted('storage.update') || (!globalStore.isPermitted('storage.create') && isNewStorage())"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-select
                  dense
                  filled
                  v-model="storage.kind"
                  :options="storageKindOptions"
                  label="Kind"
                  :readonly="!globalStore.isPermitted('storage.update') || (!globalStore.isPermitted('storage.create') && isNewStorage())"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-12 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  v-model="storage.description"
                  label="Description"
                  :readonly="!globalStore.isPermitted('storage.update') || (!globalStore.isPermitted('storage.create') && isNewStorage())"
                  type="textarea"
                />
              </div>
            </div>
          </div>
        </q-tab-panel>
        <q-tab-panel name="configuration">
          <div class="row">
            <div class="col-sm-12 col-md-4 slixx-pad-5"
                 v-for="configurationEntry in storageKindByName(storage.kind).configuration"
                 v-bind:key="configurationEntry.name">
              <div class="q-gutter-xl">
                <kind-input-component
                  :kind="configurationEntry.kind"
                  v-model="storage.configuration[configurationEntry.name]"
                  :label="configurationEntry.name"
                  :readonly="!globalStore.isPermitted('storage.update') || (!globalStore.isPermitted('storage.create') && isNewStorage())"
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
@import "./StorageDetailsPage.scss";
</style>
<script src="./StorageDetailsPage.js"/>
