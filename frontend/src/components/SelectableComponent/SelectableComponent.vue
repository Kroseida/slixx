<template>
  <div>
    <div v-if="!readonly">
      <q-input :model-value="!!selected ? displayField(selected) : ''"
               @click="show = true"
               dense
               filled
               class="ignore-readonly-design"
               readonly
               :label="label"/>
      <q-dialog v-model="show" position="right">
        <q-card class="right-side-dialog-container">
          <slot/>
        </q-card>
      </q-dialog>
    </div>
    <div v-else>
      <q-input :model-value="!!selected ? displayField(selected) : ''"
               @click="show = true"
               dense
               filled
               readonly
               :label="label"/>
    </div>
  </div>
</template>
<script src="./SelectableComponent.js"/>
<style>
.q-field--filled.ignore-readonly-design .q-field__control:before {
  background-color: transparent;
  opacity: 1;
  border-bottom-style: none;
}
</style>
