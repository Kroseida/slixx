<template>
  <div class="q-pa-md" style="height: 100%">
    <backup-list @row-click="openBackup">
      <template v-slot:action>
        <div class="action-header">
          <q-btn
            color="primary"
            label="Create Backup"
            class="action"
            style="width: 158px;"
            @click="createBackup"
            :disable="!globalStore.isPermitted('backup.create')"
          />
        </div>
      </template>
    </backup-list>
  </div>
</template>
<style lang="scss">
@import "BackupListPage";
</style>
<script src="./BackupListPage.js"/>
