<template>
  <satellite-list @row-click="openSatellite">
    <template v-slot:action>
      <div class="action-header">
        <q-btn
          color="primary"
          label="Register Satellite"
          class="action"
          style="width: 177px;"
          @click="createSatellite"
          :disable="!globalStore.isPermitted('satellite.create')"
        />
      </div>
    </template>
  </satellite-list>
</template>
<style lang="scss">
@import "./SatelliteListPage.scss";
</style>
<script src="./SatelliteListPage.js"/>
