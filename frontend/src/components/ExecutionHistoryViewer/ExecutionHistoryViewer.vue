<template>
  <table class="q-table">
    <thead>
    <tr>
      <th class="text-left">Created at</th>
      <th class="text-left">Message</th>
      <th class="text-left">Status</th>
      <th class="text-left">Percentage</th>
    </tr>
    </thead>
    <tbody>
    <tr v-for="entry in history" :key="entry.id"
        :class="{'bg-green-3': entry.statusType === 'FINISHED', 'bg-red-4': entry.statusType === 'FAILED'}">
      <td class="text-left">{{ entry.createdAt }}</td>
      <td class="text-left">{{ entry.message }}</td>
      <td class="text-left">{{ entry.statusType }}</td>
      <td class="text-left">{{ parseInt(entry.percentage * 100.0) / 100.0 }}%</td>
    </tr>
    </tbody>
  </table>
</template>
<script src="./ExecutionHistoryViewer.js"/>
