<template>
  <div class="q-pa-md">
    <q-dialog v-model="showPermissionView" position="right">
      <permission-assign-component
        :current-permissions="user.permissions"
        @add-permission="addPermission"
        @remove-permission="removePermission"
      />
    </q-dialog>
    <q-dialog v-model="showAuthenticationView" position="right">
      <authentication-component @change-password="changePassword"/>
    </q-dialog>

    <q-card>
      <q-card-section>
        <div class="relative-position row items-center">
          <div class="q-table__title">
            User
          </div>
          <div class="col"/>
          <button-group>
            <slixx-button
              color="primary"
              label="Save"
              @s-click="save"
              class="action"
              :disable="!showSaveButton() || (!globalStore.isPermitted('user.update') && !globalStore.isPermitted('user.create'))"
            />
            <slixx-button
              color="primary"
              label="Permissions"
              class="action"
              @s-click="openPermissionView"
              :disable="!showPermissionButton() || (!globalStore.isPermitted('user.permission.update'))"
            />
            <slixx-button
              color="primary"
              label="Authentication"
              class="action"
              @s-click="openAuthenticationView"
              :disable="!showAuthButton() || !globalStore.isPermitted('user.update')"
            />
            <slixx-button
              color="negative"
              label="Delete"
              class="action"
              @s-click="remove"
              :disable="!showDeleteButton() || (!globalStore.isPermitted('user.delete'))"
            />
          </button-group>
        </div>
      </q-card-section>
      <q-separator inset/>
      <q-card-section>
        <div class="row">
          <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                v-model="user.id"
                label="ID"
                readonly
              />
            </div>
          </div>
          <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                v-model="user.name"
                label="Name"
                :readonly="!globalStore.isPermitted('user.update') || (!globalStore.isPermitted('user.create') && isNewUser())"
              />
            </div>
          </div>
          <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                v-model="user.firstName"
                label="First Name"
                :readonly="!globalStore.isPermitted('user.update') || (!globalStore.isPermitted('user.create') && isNewUser())"
              />
            </div>
          </div>
          <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                v-model="user.lastName"
                label="Last Name"
                :readonly="!globalStore.isPermitted('user.update') || (!globalStore.isPermitted('user.create') && isNewUser())"
              />
            </div>
          </div>
          <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                type="email"
                v-model="user.email"
                label="Email"
                :readonly="!globalStore.isPermitted('user.update') || (!globalStore.isPermitted('user.create') && isNewUser())"
              />
            </div>
          </div>
          <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-select
                dense
                filled
                v-model="user.active"
                :options="activeOptions"
                label="Active"
                :readonly="!globalStore.isPermitted('user.update') || (!globalStore.isPermitted('user.create') && isNewUser())"
              />
            </div>
          </div>
          <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                v-model="user.createdAt"
                label="Created at"
                readonly
              />
            </div>
          </div>
          <div class="col-xs-12 col-sm-6 col-md-4 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                v-model="user.updatedAt"
                label="Updated at"
                readonly
              />
            </div>
          </div>
          <div class="col-xs-12 col-sm-12 slixx-pad-5">
            <div class="q-gutter-xl">
              <q-input
                dense
                filled
                v-model="user.description"
                label="Description"
                :readonly="!globalStore.isPermitted('user.update') || (!globalStore.isPermitted('user.create') && isNewUser())"
                type="textarea"
              />
            </div>
          </div>
        </div>
      </q-card-section>
    </q-card>
  </div>
</template>

<script src="./UserDetailsPage.js"/>
