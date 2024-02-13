<template>
  <q-dialog v-model="createSchedule" position="right">
    <job-schedule-create-component :schedule="newSchedule" :init-config="true" @s-click="onCreate"/>
  </q-dialog>
  <q-dialog v-model="updateScheduleOpen" position="right">
    <job-schedule-create-component :schedule="updateSchedule" :init-config="false" @s-click="onUpdate" @s-delete="onDelete" show-delete-button/>
  </q-dialog>
  <div class="q-pa-md schedule-list" style="height: 100%">
    <q-card class="schedule" v-for="schedule in rows" :key="schedule.id" @click="rowClick(schedule)">
      <q-card-section style="min-height: 25px">
        <div class="row">
          <div class="col-xl-1 col-6">
            <q-chip :label="schedule.kind" outline color="primary"/>
          </div>
          <div class="col-xl-2 col-6">
            <q-item-label style="margin-top: 7px; font-size: 15px">
              {{ schedule.name }}
            </q-item-label>
          </div>
          <div class="col-xl-9 col-12">
            <label style="color: gray;">{{ schedule.description }}</label>
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-card class="schedule" @click="createSchedule = true"
            v-if="globalStore.isPermitted('job.update') && allowCreation">
      <q-card-section style="min-height: 55px">
        <div class="absolute-center">
          <q-icon name="add" style="font-size: 20px; font-weight: bold"/>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>
<style>
@import "./JobScheduleViewComponent.scss";
</style>
<script src="./JobScheduleViewComponent.js"/>
