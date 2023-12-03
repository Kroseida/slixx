<template>
  <selectable-component
    :readonly="readonly"
    :label="label"
    :value="selectedJob">
    <satellite-list-component @rowClick="selectJob" :title="''" />
  </selectable-component>
</template>
<script src="./SatelliteSelectableComponent.js"/>
