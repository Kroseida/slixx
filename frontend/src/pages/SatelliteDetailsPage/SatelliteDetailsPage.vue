<template>
  <div class="q-pa-md">
    <q-card>
      <q-card-section>
        <div class="relative-position row items-center">
          <div class="q-table__title">
            Satellite
          </div>
          <div class="col"/>
          <button-group>
            <slixx-button
              color="primary"
              label="Resync"
              class="action"
              @s-click="resync"
              v-if="!isNewSatellite()"
              :disable="isNewSatellite() || !satellite.connected || !globalStore.isPermitted('satellite.resync')"
            />
            <slixx-button
              color="primary"
              label="Save"
              class="action"
              @s-click="save"
              :disable="!showSaveButton() || (!globalStore.isPermitted('satellite.update') && !globalStore.isPermitted('satellite.create'))"
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
        <q-tab name="logs" label="Logs"/>
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
                  v-model="satellite.connected"
                  label="Connected"
                  readonly
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  v-model="satellite.id"
                  label="ID"
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
                  v-model="satellite.name"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Address"
                  v-model="satellite.address"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Token"
                  v-model="satellite.token"
                  type="password"
                />
              </div>
            </div>
            <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
              <div class="q-gutter-xl">
                <q-input
                  dense
                  filled
                  label="Created at"
                  v-model="satellite.createdAt"
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
                  v-model="satellite.updatedAt"
                  readonly
                />
              </div>
            </div>
          </div>
          <div class="col-xs-12 col-sm-12 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                v-model="satellite.description"
                label="Description"
                type="textarea"
              />
            </div>
          </div>
        </q-tab-panel>
        <q-tab-panel name="logs">
          <log-viewer :logs="logs.rows"/>
          <div class="q-table__bottom row items-center justify-end">
            <div class="q-table__control">
              <span class="q-table__bottom-item">Records per page:</span>
              <q-select
                dense
                bg-color="white"
                class="q-table__bottom-item"
                v-model="pagination.rowsPerPage"
                :options="paginationOptions"
                @update:model-value="changeRowsPerPage"
              />
            </div>
            <div class="q-table__control">
              <q-pagination
                color="grey"
                @update:model-value="changePage"
                v-model="pagination.page"
                :max="logs.page.totalPages"
                input
              />
            </div>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card>
  </div>
</template>
<style lang="scss">
@import "./SatelliteDetailsPage.scss";
</style>
<script src="./SatelliteDetailsPage.js"/>
