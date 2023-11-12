<template>
  <div class="q-pa-md" style="height: 100%">
    <storage-list @row-click="openStorage">
      <template v-slot:action>
        <div class="action-header">
          <q-btn
            color="primary"
            label="Create Storage"
            class="action"
            style="width: 158px;"
            @click="createStorage"
            :disable="!globalStore.isPermitted('storage.create')"
          />
        </div>
      </template>
    </storage-list>
  </div>
</template>
<style lang="scss">
@import "./StorageListPage.scss";
</style>
<script src="./StorageListPage.js"/>
