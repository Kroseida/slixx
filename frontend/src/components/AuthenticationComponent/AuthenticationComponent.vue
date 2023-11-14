<template>
  <q-card class="right-side-dialog-container">
    <q-item-label header>Password Authentication</q-item-label>
    <q-separator inset/>
    <div class="slixx-pad-trl-15">
      <div class="q-gutter-xl">
        <q-input
          dense
          filled
          type="password"
          v-model="password"
          label="New Password"
        />
      </div>
    </div>
    <div class="slixx-pad-trl-15">
      <div class="q-gutter-xl">
        <q-input
          dense
          filled
          type="password"
          v-model="passwordConfirm"
          label="Repeat Password"
        />
      </div>
    </div>
    <div class="slixx-pad-trl-15">
      <slixx-button
        color="primary"
        size="md"
        label="Update Password"
        no-caps
        class="full-width"
        @s-click="changePassword"
        :disable="password !== passwordConfirm || password.length < 1"
      />
    </div>
  </q-card>
</template>
<style>
@import "./AuthenticationComponent.scss";
</style>
<script src="./AuthenticationComponent.js"/>
