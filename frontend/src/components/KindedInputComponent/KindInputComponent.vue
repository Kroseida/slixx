<template>
  <div>
    <password-kind-input-component
      :model-value="modelValue"
      :label="label"
      :readonly="readonly"
      v-if="kind === 'PASSWORD'"
      @update:model-value="handleInput"
    />
    <long-kind-input-component
      :model-value="modelValue"
      :label="label"
      :readonly="readonly"
      v-else-if="kind === 'LONG'"
      @update:model-value="handleInput"
    />
    <text-kind-input-component
      :model-value="modelValue"
      :label="label"
      :readonly="readonly"
      v-else
      @update:model-value="handleInput"
    />
  </div>
</template>
<script src="./KindInputComponent.js"/>
