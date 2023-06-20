<template>
  <q-layout>
    <q-header>
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer"></q-btn>
        <q-toolbar-title>Dashboard</q-toolbar-title>
        <div>slixx v{{ constantsStore.version }}</div>
      </q-toolbar>
    </q-header>
    <sidebar-component />
    <q-page-container>
      <router-view></router-view>
    </q-page-container>
  </q-layout>
</template>
<script src="./MainLayout.js" />
