<template>
  <q-drawer v-model="globalStore.leftDrawerOpen" show-if-above :width="245">
    <q-scroll-area :style="globalStore.darkMode ? 'height: 100%;' : 'height: 100%;border-right: 1px solid #ddd'">
      <q-list padding>
        <sidebar-link-component
          :key="link.title"
          v-bind="link"
          :icon="link.icon"
          v-for="link in links"
          :permission="link.permission"
        />
      </q-list>
    </q-scroll-area>
  </q-drawer>
</template>

<script src="./SidebarComponent.js" />
