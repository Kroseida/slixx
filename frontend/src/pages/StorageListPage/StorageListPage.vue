<template>
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
</template>
<style lang="scss">
@import "./StorageListPage.scss";
</style>
<script src="./StorageListPage.js"/>
