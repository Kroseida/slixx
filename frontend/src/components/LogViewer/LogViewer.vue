<template>
  <div>
    <q-toolbar>
      <q-toolbar-title>
      </q-toolbar-title>
      <q-btn flat round dense>
        <q-icon name="view_stream" @click="displayType = true"/>
      </q-btn>
      <q-btn flat round dense>
        <q-icon name="article" @click="displayType = false"/>
      </q-btn>
      <q-btn flat round dense>
        <q-icon name="download" @click="downloadLogs"/>
      </q-btn>
    </q-toolbar>
    <div class="content">
      <div class="q-table__container q-table--horizontal-separator" v-if="displayType">
        <div class="q-table__middle scroll">
          <table class="q-table">
            <thead>
            <tr>
              <th class="text-left">Time</th>
              <th class="text-left">Sender</th>
              <th class="text-left">Level</th>
              <th class="text-left">Message</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="log in logs" :key="log.id">
              <td class="text-left" style="width: 300px">{{ log.loggedAt }}</td>
              <td class="text-left" style="width: 140px" :class="log.sender === 'supervisor' ? 'bg-blue-grey-2' : 'bg-teal-2'">{{ log.sender }}</td>
              <td class="text-left" :class="classFromLevel(log.level)" style="width: 100px">{{ log.level }}</td>
              <td class="text-left">{{ log.message }}</td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div class="console" v-if="!displayType">
        <textarea class="console-textarea" v-model="logsConsole" readonly></textarea>
      </div>
    </div>
  </div>
</template>
<script src="./LogViewer.js"/>
<style>
@import "./LogViewer.scss";
</style>
