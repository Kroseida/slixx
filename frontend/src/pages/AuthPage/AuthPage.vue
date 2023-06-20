<template>
  <div>
    <q-layout view="lHh Lpr lFf">
      <q-page-container>
        <q-page class="flex flex-center">
          <q-card class="q-pa-md shadow-2" :bordered="true" :dark="globalStore.darkMode">
            <q-card-section class="text-center">
              <q-toolbar-title class="text-h5 text-weight-bold">Sign in</q-toolbar-title>
            </q-card-section>
            <q-card-section>
              <q-input
                dense
                outlined
                v-model="name"
                label="Name"
                input-style="width: 280px"
                :dark="globalStore.darkMode"
              />
              <q-input
                dense
                outlined
                class="q-mt-md"
                v-model="password"
                type="password"
                label="Password"
                input-style="width: 280px"
                :dark="globalStore.darkMode"
              />
            </q-card-section>
            <q-card-section>
              <q-btn
                style="border-radius: 8px;"
                color="primary"
                rounded
                size="md"
                label="Sign in"
                no-caps
                class="full-width"
                @click="authenticate"
              />
            </q-card-section>
          </q-card>
        </q-page>
      </q-page-container>
    </q-layout>
  </div>
</template>
<script src="./AuthPageScript.js"/>
