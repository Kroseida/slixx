<template>
  <div class="q-pa-md" style="height: 100%">
    <user-list @row-click="openUser">
      <template v-slot:action>
        <div class="action-header">
          <q-btn
            color="primary"
            label="Create User"
            class="action"
            @click="createUser"
            :disable="!globalStore.isPermitted('user.create')"
          />
        </div>
      </template>
    </user-list>
  </div>
</template>
<style lang="scss">
@import "./UserListPage.scss";
</style>
<script src="./UserListPage.js"/>
