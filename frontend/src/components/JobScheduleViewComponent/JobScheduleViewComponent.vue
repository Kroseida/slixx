<template>
  <q-dialog v-model="createSchedule" position="right">
    <job-schedule-create-component :job-id="jobId"   />
  </q-dialog>
  <div class="q-pa-md" style="height: 100%">
    <q-card v-for="schedule in rows" :key="schedule.id">
      {{schedule}}
    </q-card>

    <q-card class="schedule" @click="createSchedule = true">
      <q-card-section>
        Create
      </q-card-section>
    </q-card>
  </div>
</template>
<style>
@import "./JobScheduleViewComponent.scss";
</style>
<script src="./JobScheduleViewComponent.js"/>
