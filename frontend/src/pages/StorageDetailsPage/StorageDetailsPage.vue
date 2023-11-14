<template>
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
              @s-click="remove"
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
