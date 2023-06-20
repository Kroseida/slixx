<template>
  <job-list @row-click="openJob">
    <template v-slot:action>
      <div class="action-header">
        <q-btn
          color="primary"
          label="Create Job"
          class="action"
          @click="createJob"
          :disable="!globalStore.isPermitted('job.create')"
        />
      </div>
    </template>
  </job-list>
</template>
<style lang="scss">
@import "JobListPage";
</style>
<script src="./JobListPage.js"/>
