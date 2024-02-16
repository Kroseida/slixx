<template>
  <q-card class="right-side-dialog-container">
    <q-item-label header>User Permissions</q-item-label>
    <q-separator inset/>
    <q-item
      tag="label" v-ripple
      v-for="permission in sortedPermissions()"
      v-bind:key="permission.name">
      <q-item-section side top>
        <q-checkbox
          v-model="permissionMap[permission.value]"
          @update:model-value="(value) => togglePermission(permission.value, value)"
        />
      </q-item-section>

      <q-item-section>
        <q-item-label>{{ permission.name }}</q-item-label>
        <q-item-label caption>
          {{ permission.value }}
        </q-item-label>
      </q-item-section>
    </q-item>
  </q-card>
</template>
<style>
@import "./PermissionAssignComponent.scss";
</style>
<script src="./PermissionAssignComponent.js"/>
