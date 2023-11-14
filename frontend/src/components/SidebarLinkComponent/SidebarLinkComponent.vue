<template>
  <q-item clickable tag="a" :href="link" v-if="globalStore.isPermitted(permission)">
    <q-item-section v-if="icon" avatar style="padding-left: 16px;">
      <q-icon :name="icon" style="font-size: 16px;"/>
    </q-item-section>
    <q-item-section>
      <q-item-label style="font-size: 15px">{{ title }}</q-item-label>
    </q-item-section>
  </q-item>
</template>

<script src="./SidebarLinkComponent.js" />
