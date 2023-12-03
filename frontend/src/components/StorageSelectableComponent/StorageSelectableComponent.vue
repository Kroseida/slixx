<template>
  <selectable-component
    :readonly="readonly"
    :label="label"
    :value="selectedJob">
    <storage-list-component @rowClick="selectJob" :title="''" />
  </selectable-component>
</template>
<script src="./StorageSelectableComponent.js"/>
