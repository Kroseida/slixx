<template>
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
</template>
<style lang="scss">
@import "./UserListPage.scss";
</style>
<script src="./UserListPage.js"/>
